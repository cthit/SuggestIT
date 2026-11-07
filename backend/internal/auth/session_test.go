package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func testKeys(t *testing.T, seed string) Keys {
	t.Helper()
	k, err := DeriveKeys([]byte(strings.Repeat(seed, 32)))
	if err != nil {
		t.Fatal(err)
	}
	return k
}

var alice = User{CID: "alice", Nickname: "Ali", Name: "Alice A", Picture: "https://example.com/a.png"}

// issue runs Set and returns the resulting cookie.
func issue(t *testing.T, s *Sessions, u User) *http.Cookie {
	t.Helper()
	rec := httptest.NewRecorder()
	if err := s.Set(rec, u); err != nil {
		t.Fatal(err)
	}
	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("got %d cookies", len(cookies))
	}
	return cookies[0]
}

func reqWith(c *http.Cookie) *http.Request {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	if c != nil {
		r.AddCookie(c)
	}
	return r
}

func TestDeriveKeys(t *testing.T) {
	if _, err := DeriveKeys([]byte("short")); err == nil {
		t.Error("expected error for short secret")
	}
	a := testKeys(t, "a")
	if len(a.Hash) != 32 || len(a.Encrypt) != 32 {
		t.Fatalf("key lengths %d/%d", len(a.Hash), len(a.Encrypt))
	}
	if bytes.Equal(a.Hash, a.Encrypt) {
		t.Error("hash and encryption keys must differ")
	}
	if again := testKeys(t, "a"); !bytes.Equal(a.Hash, again.Hash) || !bytes.Equal(a.Encrypt, again.Encrypt) {
		t.Error("derivation must be deterministic")
	}
	if b := testKeys(t, "b"); bytes.Equal(a.Hash, b.Hash) {
		t.Error("different secrets must give different keys")
	}
}

func TestSessionRoundTrip(t *testing.T) {
	s := NewSessions(testKeys(t, "k"), CookieOptions{Secure: true})
	c := issue(t, s, alice)

	if c.Name != "__Host-suggestit_session" {
		t.Errorf("name = %q", c.Name)
	}
	if !c.HttpOnly || !c.Secure || c.SameSite != http.SameSiteLaxMode || c.Path != "/" || c.Domain != "" {
		t.Errorf("bad cookie attributes: %+v", c)
	}
	if c.MaxAge != int(SessionMaxAge/time.Second) {
		t.Errorf("MaxAge = %d", c.MaxAge)
	}
	if strings.Contains(c.Value, "alice") {
		t.Error("cookie value must be encrypted")
	}

	got, err := s.Get(reqWith(c))
	if err != nil {
		t.Fatal(err)
	}
	if got != alice {
		t.Errorf("got %+v, want %+v", got, alice)
	}
}

func TestSessionCookieNaming(t *testing.T) {
	keys := testKeys(t, "k")
	cases := []struct {
		opts       CookieOptions
		name       string
		domain     string
		wantSecure bool
	}{
		{CookieOptions{Secure: true}, "__Host-suggestit_session", "", true},
		{CookieOptions{Secure: true, Domain: "chalmers.it"}, "suggestit_session", "chalmers.it", true},
		{CookieOptions{Secure: false}, "suggestit_session", "", false},
	}
	for _, tc := range cases {
		s := NewSessions(keys, tc.opts)
		c := issue(t, s, alice)
		if c.Name != tc.name || c.Domain != tc.domain || c.Secure != tc.wantSecure || s.CookieName() != tc.name {
			t.Errorf("%+v: got name=%q domain=%q secure=%v", tc.opts, c.Name, c.Domain, c.Secure)
		}
	}
}

func TestSessionRejectsTampering(t *testing.T) {
	s := NewSessions(testKeys(t, "k"), CookieOptions{Secure: true})
	c := issue(t, s, alice)

	b := []byte(c.Value)
	i := len(b) / 2
	if b[i] == 'A' {
		b[i] = 'B'
	} else {
		b[i] = 'A'
	}
	tampered := *c
	tampered.Value = string(b)
	if _, err := s.Get(reqWith(&tampered)); err != ErrNoSession {
		t.Errorf("tampered cookie: err = %v", err)
	}

	garbage := *c
	garbage.Value = "not-a-cookie"
	if _, err := s.Get(reqWith(&garbage)); err != ErrNoSession {
		t.Errorf("garbage cookie: err = %v", err)
	}

	other := NewSessions(testKeys(t, "x"), CookieOptions{Secure: true})
	if _, err := other.Get(reqWith(c)); err != ErrNoSession {
		t.Errorf("cookie from another secret: err = %v", err)
	}

	if _, err := s.Get(reqWith(nil)); err != ErrNoSession {
		t.Errorf("no cookie: err = %v", err)
	}
}

func TestSessionExpiry(t *testing.T) {
	s := NewSessions(testKeys(t, "k"), CookieOptions{Secure: true})

	s.now = func() time.Time { return time.Now().Add(-SessionMaxAge - time.Minute) }
	old := issue(t, s, alice)
	s.now = time.Now
	if _, err := s.Get(reqWith(old)); err != ErrNoSession {
		t.Errorf("expired session accepted: %v", err)
	}

	s.now = func() time.Time { return time.Now().Add(-SessionMaxAge + time.Minute) }
	almost := issue(t, s, alice)
	s.now = time.Now
	if _, err := s.Get(reqWith(almost)); err != nil {
		t.Errorf("session within max age rejected: %v", err)
	}

	s.now = func() time.Time { return time.Now().Add(time.Hour) }
	future := issue(t, s, alice)
	s.now = time.Now
	if _, err := s.Get(reqWith(future)); err != ErrNoSession {
		t.Errorf("session issued in the future accepted: %v", err)
	}
}

func TestSessionRequiresCID(t *testing.T) {
	s := NewSessions(testKeys(t, "k"), CookieOptions{})
	if err := s.Set(httptest.NewRecorder(), User{Nickname: "nocid"}); err == nil {
		t.Error("Set without cid must fail")
	}
}

func TestClear(t *testing.T) {
	s := NewSessions(testKeys(t, "k"), CookieOptions{Secure: true})
	rec := httptest.NewRecorder()
	s.Clear(rec)
	c := rec.Result().Cookies()
	if len(c) != 1 || c[0].Name != s.CookieName() || c[0].MaxAge >= 0 || !c[0].Secure || c[0].Path != "/" {
		t.Errorf("bad delete cookie: %+v", c)
	}
}

func TestRequire(t *testing.T) {
	s := NewSessions(testKeys(t, "k"), CookieOptions{Secure: true})
	var seen User
	h := s.Require(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, ok := UserFrom(r.Context())
		if !ok {
			t.Error("no user in context")
		}
		seen = u
		w.WriteHeader(http.StatusTeapot)
	}))

	t.Run("no cookie", func(t *testing.T) {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, reqWith(nil))
		assertUnauthorized(t, rec)
		if len(rec.Result().Cookies()) != 0 {
			t.Error("must not set cookies when none was sent")
		}
	})

	t.Run("invalid cookie is deleted", func(t *testing.T) {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, reqWith(&http.Cookie{Name: s.CookieName(), Value: "bogus"}))
		assertUnauthorized(t, rec)
		c := rec.Result().Cookies()
		if len(c) != 1 || c[0].MaxAge >= 0 {
			t.Errorf("invalid cookie not deleted: %+v", c)
		}
	})

	t.Run("valid", func(t *testing.T) {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, reqWith(issue(t, s, alice)))
		if rec.Code != http.StatusTeapot || seen != alice {
			t.Errorf("code=%d user=%+v", rec.Code, seen)
		}
	})
}

func assertUnauthorized(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("code = %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q", ct)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body["error"] == "" {
		t.Errorf("body = %q", rec.Body.String())
	}
}
