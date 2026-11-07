package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"chalmers.it/suggestit/internal/auth"
	"chalmers.it/suggestit/internal/store"
)

// fakeStore is an in-memory Store.
type fakeStore struct {
	mu      sync.Mutex
	items   []store.Suggestion // newest first
	clock   time.Time
	pingErr error
	failAll error
	deleted [][]string
}

func (f *fakeStore) Ping(context.Context) error { return f.pingErr }

func (f *fakeStore) List(context.Context) ([]store.Suggestion, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failAll != nil {
		return nil, f.failAll
	}
	if len(f.items) == 0 {
		return nil, nil // the handler must still return []
	}
	return slices.Clone(f.items), nil
}

func (f *fakeStore) Get(_ context.Context, id string) (store.Suggestion, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failAll != nil {
		return store.Suggestion{}, f.failAll
	}
	for _, s := range f.items {
		if s.ID == id {
			return s, nil
		}
	}
	return store.Suggestion{}, store.ErrNotFound
}

func (f *fakeStore) Create(_ context.Context, n store.NewSuggestion) (store.Suggestion, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failAll != nil {
		return store.Suggestion{}, f.failAll
	}
	f.clock = f.clock.Add(time.Minute)
	s := store.Suggestion{ID: store.NewUUID(), Timestamp: f.clock, Title: n.Title, Text: n.Text, Author: n.Author}
	f.items = append([]store.Suggestion{s}, f.items...)
	return s, nil
}

func (f *fakeStore) Delete(_ context.Context, id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failAll != nil {
		return f.failAll
	}
	before := len(f.items)
	f.items = slices.DeleteFunc(f.items, func(s store.Suggestion) bool { return s.ID == id })
	if len(f.items) == before {
		return store.ErrNotFound
	}
	return nil
}

func (f *fakeStore) DeleteMany(_ context.Context, ids []string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failAll != nil {
		return f.failAll
	}
	f.deleted = append(f.deleted, ids)
	f.items = slices.DeleteFunc(f.items, func(s store.Suggestion) bool { return slices.Contains(ids, s.ID) })
	return nil
}

type testEnv struct {
	store    *fakeStore
	sessions *auth.Sessions
	handler  http.Handler
	session  *http.Cookie
}

func newEnv(t *testing.T) *testEnv {
	t.Helper()
	keys, err := auth.DeriveKeys([]byte(strings.Repeat("k", 32)))
	if err != nil {
		t.Fatal(err)
	}
	sessions := auth.NewSessions(keys, auth.CookieOptions{Secure: true})
	fs := &fakeStore{clock: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)}
	h, err := New(Options{
		Store:    fs,
		Sessions: sessions,
		Login: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://gamma.example/authorize", http.StatusFound)
		}),
		Callback: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/", http.StatusFound)
		}),
		Logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		TrustedOrigins: []string{"https://suggestit.example"},
	})
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	if err := sessions.Set(rec, auth.User{CID: "alice", Nickname: "Ali"}); err != nil {
		t.Fatal(err)
	}
	return &testEnv{store: fs, sessions: sessions, handler: h, session: rec.Result().Cookies()[0]}
}

type reqOpt func(*http.Request)

func authed(e *testEnv) reqOpt { return func(r *http.Request) { r.AddCookie(e.session) } }

func header(k, v string) reqOpt { return func(r *http.Request) { r.Header.Set(k, v) } }

func (e *testEnv) do(method, path, body string, opts ...reqOpt) *httptest.ResponseRecorder {
	var rd io.Reader
	if body != "" {
		rd = strings.NewReader(body)
	}
	r := httptest.NewRequest(method, path, rd)
	if body != "" {
		r.Header.Set("Content-Type", "application/json")
	}
	for _, o := range opts {
		o(r)
	}
	rec := httptest.NewRecorder()
	e.handler.ServeHTTP(rec, r)
	return rec
}

func expect(t *testing.T, rec *httptest.ResponseRecorder, status int) {
	t.Helper()
	if rec.Code != status {
		t.Fatalf("status = %d, want %d (body %q)", rec.Code, status, rec.Body.String())
	}
	if status >= 400 && rec.Header().Get("Content-Type") == "application/json" {
		var body map[string]string
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body["error"] == "" {
			t.Errorf("error body = %q", rec.Body.String())
		}
	}
}

func decode[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q", ct)
	}
	var v T
	if err := json.Unmarshal(rec.Body.Bytes(), &v); err != nil {
		t.Fatalf("decode %q: %v", rec.Body.String(), err)
	}
	return v
}

func (e *testEnv) create(t *testing.T, title string) store.Suggestion {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"title": title, "text": "text of " + title, "author": "Bob"})
	rec := e.do(http.MethodPost, "/api/suggestions", string(body))
	expect(t, rec, http.StatusCreated)
	return decode[store.Suggestion](t, rec)
}

func TestAuthRequired(t *testing.T) {
	e := newEnv(t)
	id := store.NewUUID()
	for _, tc := range []struct{ method, path, body string }{
		{http.MethodGet, "/api/suggestions", ""},
		{http.MethodGet, "/api/suggestions/" + id, ""},
		{http.MethodDelete, "/api/suggestions/" + id, ""},
		{http.MethodPost, "/api/suggestions/delete", `{"ids":[]}`},
		{http.MethodGet, "/api/me", ""},
	} {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			expect(t, e.do(tc.method, tc.path, tc.body), http.StatusUnauthorized)
			bad := &http.Cookie{Name: e.sessions.CookieName(), Value: "forged"}
			expect(t, e.do(tc.method, tc.path, tc.body, func(r *http.Request) { r.AddCookie(bad) }), http.StatusUnauthorized)
		})
	}
}

func TestMe(t *testing.T) {
	e := newEnv(t)
	rec := e.do(http.MethodGet, "/api/me", "", authed(e))
	expect(t, rec, http.StatusOK)
	got := decode[map[string]string](t, rec)
	want := map[string]string{"cid": "alice", "nickname": "Ali", "name": "", "picture": ""}
	if len(got) != len(want) {
		t.Fatalf("got %v", got)
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s = %q, want %q", k, got[k], v)
		}
	}
}

func TestMeDeletesInvalidCookie(t *testing.T) {
	e := newEnv(t)
	bad := &http.Cookie{Name: e.sessions.CookieName(), Value: "forged"}
	rec := e.do(http.MethodGet, "/api/me", "", func(r *http.Request) { r.AddCookie(bad) })
	expect(t, rec, http.StatusUnauthorized)
	c := rec.Result().Cookies()
	if len(c) != 1 || c[0].Name != e.sessions.CookieName() || c[0].MaxAge >= 0 {
		t.Errorf("invalid cookie not cleared: %+v", c)
	}
}

func TestLogout(t *testing.T) {
	e := newEnv(t)
	rec := e.do(http.MethodPost, "/api/auth/logout", "", authed(e))
	expect(t, rec, http.StatusNoContent)
	c := rec.Result().Cookies()
	if len(c) != 1 || c[0].Name != e.sessions.CookieName() || c[0].MaxAge >= 0 {
		t.Errorf("logout did not clear the session cookie: %+v", c)
	}
}

func TestAuthRoutesWired(t *testing.T) {
	e := newEnv(t)
	rec := e.do(http.MethodGet, "/api/auth/login", "")
	expect(t, rec, http.StatusFound)
	if loc := rec.Header().Get("Location"); loc != "https://gamma.example/authorize" {
		t.Errorf("login location = %q", loc)
	}
	rec = e.do(http.MethodGet, "/api/auth/callback?code=x&state=y", "")
	expect(t, rec, http.StatusFound)
}

func TestCreateAndList(t *testing.T) {
	e := newEnv(t)

	rec := e.do(http.MethodGet, "/api/suggestions", "", authed(e))
	expect(t, rec, http.StatusOK)
	if strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Errorf("empty list = %q, want []", rec.Body.String())
	}

	rec = e.do(http.MethodPost, "/api/suggestions", `{"title":"  Coffee  ","text":"\n More coffee \n","author":""}`)
	expect(t, rec, http.StatusCreated)
	first := decode[store.Suggestion](t, rec)
	if first.Title != "Coffee" || first.Text != "More coffee" || first.Author != "Anonymous" || !store.ValidUUID(first.ID) {
		t.Errorf("created = %+v", first)
	}
	second := e.create(t, "Tea")
	if second.Author != "Bob" {
		t.Errorf("author = %q", second.Author)
	}

	rec = e.do(http.MethodGet, "/api/suggestions", "", authed(e))
	expect(t, rec, http.StatusOK)
	list := decode[[]store.Suggestion](t, rec)
	if len(list) != 2 || list[0].ID != second.ID || list[1].ID != first.ID {
		t.Errorf("list not newest first: %+v", list)
	}

	var raw []map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &raw)
	ts, _ := raw[0]["timestamp"].(string)
	if _, err := time.Parse(time.RFC3339, ts); err != nil || !strings.HasSuffix(ts, "Z") {
		t.Errorf("timestamp %q is not RFC3339 UTC", ts)
	}
	for _, k := range []string{"id", "timestamp", "title", "text", "author"} {
		if _, ok := raw[0][k]; !ok {
			t.Errorf("suggestion JSON lacks %q: %v", k, raw[0])
		}
	}
}

func TestCreateAnonymousToggle(t *testing.T) {
	e := newEnv(t)
	for _, author := range []string{"", "   ", "anonymous", "ANONYMOUS"} {
		body, _ := json.Marshal(map[string]string{"title": "t", "text": "x", "author": author})
		rec := e.do(http.MethodPost, "/api/suggestions", string(body))
		expect(t, rec, http.StatusCreated)
		if got := decode[store.Suggestion](t, rec).Author; got != "Anonymous" {
			t.Errorf("author %q -> %q", author, got)
		}
	}
}

func TestCreateValidation(t *testing.T) {
	e := newEnv(t)
	long := func(n int) string { return strings.Repeat("å", n) } // multi-byte: limits are in runes
	ok := func(title, text, author string) string {
		b, _ := json.Marshal(map[string]string{"title": title, "text": text, "author": author})
		return string(b)
	}

	expect(t, e.do(http.MethodPost, "/api/suggestions", ok(long(MaxTitleLen), long(MaxTextLen), long(MaxAuthorLen))), http.StatusCreated)

	cases := map[string]struct {
		body   string
		status int
		opts   []reqOpt
	}{
		"missing title":     {ok("", "x", ""), http.StatusBadRequest, nil},
		"blank title":       {ok("   ", "x", ""), http.StatusBadRequest, nil},
		"missing text":      {ok("t", "", ""), http.StatusBadRequest, nil},
		"blank text":        {ok("t", " \n\t ", ""), http.StatusBadRequest, nil},
		"title too long":    {ok(long(MaxTitleLen+1), "x", ""), http.StatusBadRequest, nil},
		"text too long":     {ok("t", long(MaxTextLen+1), ""), http.StatusBadRequest, nil},
		"author too long":   {ok("t", "x", long(MaxAuthorLen+1)), http.StatusBadRequest, nil},
		"NUL byte":          {ok("t\x00", "x", ""), http.StatusBadRequest, nil},
		"invalid json":      {`{"title":`, http.StatusBadRequest, nil},
		"wrong type":        {`{"title":1,"text":"x"}`, http.StatusBadRequest, nil},
		"trailing data":     {ok("t", "x", "") + `{}`, http.StatusBadRequest, nil},
		"array body":        {`[]`, http.StatusBadRequest, nil},
		"too large":         {ok("t", strings.Repeat("a", MaxBodyBytes+1), ""), http.StatusRequestEntityTooLarge, nil},
		"form content type": {ok("t", "x", ""), http.StatusUnsupportedMediaType, []reqOpt{header("Content-Type", "application/x-www-form-urlencoded")}},
		"no content type":   {ok("t", "x", ""), http.StatusUnsupportedMediaType, []reqOpt{header("Content-Type", "")}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			expect(t, e.do(http.MethodPost, "/api/suggestions", tc.body, tc.opts...), tc.status)
		})
	}

	expect(t, e.do(http.MethodPost, "/api/suggestions", ok("t", "x", ""), header("Content-Type", "application/json; charset=utf-8")), http.StatusCreated)
	if n := len(e.store.items); n != 2 {
		t.Errorf("store has %d items, want 2 (invalid requests must not be stored)", n)
	}
}

func TestGetSuggestion(t *testing.T) {
	e := newEnv(t)
	s := e.create(t, "Coffee")

	rec := e.do(http.MethodGet, "/api/suggestions/"+s.ID, "", authed(e))
	expect(t, rec, http.StatusOK)
	if got := decode[store.Suggestion](t, rec); got != s {
		t.Errorf("got %+v, want %+v", got, s)
	}

	expect(t, e.do(http.MethodGet, "/api/suggestions/"+store.NewUUID(), "", authed(e)), http.StatusNotFound)
	expect(t, e.do(http.MethodGet, "/api/suggestions/not-a-uuid", "", authed(e)), http.StatusNotFound)
	expect(t, e.do(http.MethodGet, "/api/suggestions/delete", "", authed(e)), http.StatusNotFound)
}

func TestDeleteSuggestion(t *testing.T) {
	e := newEnv(t)
	s := e.create(t, "Coffee")
	keep := e.create(t, "Tea")

	expect(t, e.do(http.MethodDelete, "/api/suggestions/"+s.ID, "", authed(e)), http.StatusNoContent)
	expect(t, e.do(http.MethodDelete, "/api/suggestions/"+s.ID, "", authed(e)), http.StatusNotFound)
	expect(t, e.do(http.MethodDelete, "/api/suggestions/nope", "", authed(e)), http.StatusNotFound)
	if len(e.store.items) != 1 || e.store.items[0].ID != keep.ID {
		t.Errorf("items = %+v", e.store.items)
	}
}

func TestDeleteMany(t *testing.T) {
	e := newEnv(t)
	a, b, c := e.create(t, "a"), e.create(t, "b"), e.create(t, "c")

	body, _ := json.Marshal(map[string][]string{"ids": {a.ID, c.ID}})
	expect(t, e.do(http.MethodPost, "/api/suggestions/delete", string(body), authed(e)), http.StatusNoContent)
	if len(e.store.items) != 1 || e.store.items[0].ID != b.ID {
		t.Errorf("items = %+v", e.store.items)
	}
	if len(e.store.deleted) != 1 || len(e.store.deleted[0]) != 2 {
		t.Errorf("DeleteMany calls = %v, want one call with both ids", e.store.deleted)
	}

	expect(t, e.do(http.MethodPost, "/api/suggestions/delete", `{"ids":["nope"]}`, authed(e)), http.StatusBadRequest)
	expect(t, e.do(http.MethodPost, "/api/suggestions/delete", `{"ids":"x"}`, authed(e)), http.StatusBadRequest)
	expect(t, e.do(http.MethodPost, "/api/suggestions/delete", `{"ids":[]}`, authed(e)), http.StatusNoContent)
	if len(e.store.items) != 1 {
		t.Error("invalid requests must not delete anything")
	}
}

func TestStoreErrorsAre500(t *testing.T) {
	e := newEnv(t)
	s := e.create(t, "x")
	e.store.failAll = errors.New("db down: secret detail")
	for _, rec := range []*httptest.ResponseRecorder{
		e.do(http.MethodGet, "/api/suggestions", "", authed(e)),
		e.do(http.MethodGet, "/api/suggestions/"+s.ID, "", authed(e)),
		e.do(http.MethodDelete, "/api/suggestions/"+s.ID, "", authed(e)),
		e.do(http.MethodPost, "/api/suggestions/delete", `{"ids":["`+s.ID+`"]}`, authed(e)),
		e.do(http.MethodPost, "/api/suggestions", `{"title":"t","text":"x"}`),
	} {
		expect(t, rec, http.StatusInternalServerError)
		if strings.Contains(rec.Body.String(), "secret detail") {
			t.Error("internal error details leaked")
		}
	}
}

func TestHealthz(t *testing.T) {
	e := newEnv(t)
	rec := e.do(http.MethodGet, "/healthz", "")
	expect(t, rec, http.StatusOK)
	if strings.TrimSpace(rec.Body.String()) != "ok" {
		t.Errorf("body = %q", rec.Body.String())
	}
	e.store.pingErr = errors.New("down")
	expect(t, e.do(http.MethodGet, "/healthz", ""), http.StatusServiceUnavailable)
}

func TestCrossOriginProtection(t *testing.T) {
	e := newEnv(t)
	body := `{"title":"t","text":"x"}`

	expect(t, e.do(http.MethodPost, "/api/suggestions", body, header("Sec-Fetch-Site", "cross-site")), http.StatusForbidden)
	expect(t, e.do(http.MethodPost, "/api/suggestions", body, header("Sec-Fetch-Site", "same-site")), http.StatusForbidden)
	expect(t, e.do(http.MethodPost, "/api/auth/logout", "", authed(e), header("Sec-Fetch-Site", "cross-site")), http.StatusForbidden)
	expect(t, e.do(http.MethodDelete, "/api/suggestions/"+store.NewUUID(), "", authed(e), header("Origin", "https://evil.example")), http.StatusForbidden)
	if len(e.store.items) != 0 {
		t.Fatal("cross-origin request reached the store")
	}

	expect(t, e.do(http.MethodPost, "/api/suggestions", body, header("Sec-Fetch-Site", "same-origin")), http.StatusCreated)
	expect(t, e.do(http.MethodPost, "/api/suggestions", body, header("Origin", "http://example.com")), http.StatusCreated) // httptest Host is example.com
	// Old browser (no Sec-Fetch-Site) behind a proxy that rewrites Host: the
	// public origin is trusted explicitly.
	behindProxy := func(r *http.Request) { r.Host = "backend:3001" }
	expect(t, e.do(http.MethodPost, "/api/suggestions", body, behindProxy, header("Origin", "https://suggestit.example")), http.StatusCreated)
	expect(t, e.do(http.MethodPost, "/api/suggestions", body, behindProxy, header("Origin", "https://evil.example")), http.StatusForbidden)
	// Safe methods are never blocked.
	expect(t, e.do(http.MethodGet, "/api/suggestions", "", authed(e), header("Sec-Fetch-Site", "cross-site")), http.StatusOK)
}

func TestHeadersAndNotFound(t *testing.T) {
	e := newEnv(t)
	rec := e.do(http.MethodGet, "/api/me", "", authed(e))
	if rec.Header().Get("X-Content-Type-Options") != "nosniff" || rec.Header().Get("Cache-Control") != "no-store" {
		t.Errorf("missing security headers: %v", rec.Header())
	}
	rec = e.do(http.MethodGet, "/api/nope", "")
	expect(t, rec, http.StatusNotFound)
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Error("unknown API route must return JSON")
	}
	if rec.Header().Get("Allow") != "" {
		t.Errorf("404 must not have an Allow header, got %q", rec.Header().Get("Allow"))
	}

	// Known paths with the wrong method give a JSON 405 with Allow.
	for _, tc := range []struct{ method, path, allow string }{
		{http.MethodPut, "/api/suggestions", "GET, HEAD, POST"},
		{http.MethodPatch, "/api/me", "GET, HEAD"},
		{http.MethodGet, "/api/auth/logout", "POST"},
		{http.MethodPost, "/api/suggestions/" + store.NewUUID(), "GET, HEAD, DELETE"},
	} {
		rec := e.do(tc.method, tc.path, "", authed(e))
		expect(t, rec, http.StatusMethodNotAllowed)
		if got := rec.Header().Get("Allow"); got != tc.allow {
			t.Errorf("%s %s: Allow = %q, want %q", tc.method, tc.path, got, tc.allow)
		}
		if rec.Header().Get("Content-Type") != "application/json" || !strings.Contains(rec.Body.String(), "method not allowed") {
			t.Errorf("%s %s: want JSON 405 body, got %q", tc.method, tc.path, rec.Body.String())
		}
	}
}
