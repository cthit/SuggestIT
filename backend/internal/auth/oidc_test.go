package auth

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
)

const (
	testClientID     = "suggestit"
	testClientSecret = "s3cret"
	testRedirect     = "http://localhost:3000/callback"
	testSubject      = "user-1"
)

var b64 = base64.RawURLEncoding

// fakeGamma is a minimal OIDC provider: discovery, JWKS, token (with PKCE
// check) and userinfo. ID tokens are RS256 JWTs signed with stdlib crypto.
type fakeGamma struct {
	t        *testing.T
	srv      *httptest.Server
	key      *rsa.PrivateKey
	issuer   string // advertised issuer (defaults to srv.URL)
	userinfo map[string]any

	mu        sync.Mutex
	challenge string // code_challenge from the authorize request
	code      string
	tokenErr  bool
}

func newFakeGamma(t *testing.T) *fakeGamma {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	g := &fakeGamma{t: t, key: key, userinfo: map[string]any{
		"sub": testSubject, "cid": "alice", "nickname": "Ali", "name": "Alice A", "picture": "https://example.com/a.png",
	}}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /.well-known/openid-configuration", g.discovery)
	mux.HandleFunc("GET /jwks", g.jwks)
	mux.HandleFunc("POST /token", g.token)
	mux.HandleFunc("GET /userinfo", g.userinfoHandler)
	g.srv = httptest.NewServer(mux)
	t.Cleanup(g.srv.Close)
	g.issuer = g.srv.URL
	return g
}

func (g *fakeGamma) discovery(w http.ResponseWriter, _ *http.Request) {
	u := g.srv.URL
	writeTestJSON(w, map[string]any{
		"issuer":                                g.issuer,
		"authorization_endpoint":                u + "/authorize",
		"token_endpoint":                        u + "/token",
		"userinfo_endpoint":                     u + "/userinfo",
		"jwks_uri":                              u + "/jwks",
		"code_challenge_methods_supported":      []string{"S256"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
		"token_endpoint_auth_methods_supported": []string{"client_secret_basic"},
		"response_types_supported":              []string{"code"},
	})
}

func (g *fakeGamma) jwks(w http.ResponseWriter, _ *http.Request) {
	pub := g.key.PublicKey
	writeTestJSON(w, map[string]any{"keys": []map[string]any{{
		"kty": "RSA", "kid": "k1", "use": "sig", "alg": "RS256",
		"n": b64.EncodeToString(pub.N.Bytes()),
		"e": b64.EncodeToString(big.NewInt(int64(pub.E)).Bytes()),
	}}})
}

// authorize simulates the user logging in at Gamma: it records the PKCE
// challenge and returns the redirect back to the client.
func (g *fakeGamma) authorize(authURL string) (code, state string) {
	g.t.Helper()
	u, err := url.Parse(authURL)
	if err != nil {
		g.t.Fatal(err)
	}
	q := u.Query()
	if u.Path != "/authorize" || q.Get("client_id") != testClientID || q.Get("redirect_uri") != testRedirect ||
		q.Get("response_type") != "code" || q.Get("scope") != "openid profile" || q.Get("code_challenge_method") != "S256" {
		g.t.Fatalf("unexpected authorize URL %s", authURL)
	}
	if q.Get("state") == "" || q.Get("code_challenge") == "" {
		g.t.Fatalf("missing state or PKCE challenge in %s", authURL)
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.challenge = q.Get("code_challenge")
	g.code = "code-" + q.Get("state")
	return g.code, q.Get("state")
}

func (g *fakeGamma) token(w http.ResponseWriter, r *http.Request) {
	g.mu.Lock()
	defer g.mu.Unlock()
	id, secret, ok := r.BasicAuth()
	if !ok {
		id, secret = r.PostFormValue("client_id"), r.PostFormValue("client_secret")
	}
	verifier := r.PostFormValue("code_verifier")
	sum := sha256.Sum256([]byte(verifier))
	switch {
	case g.tokenErr:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, `{"error":"invalid_grant"}`)
		return
	case id != testClientID || secret != testClientSecret:
		http.Error(w, `{"error":"invalid_client"}`, http.StatusUnauthorized)
		return
	case r.PostFormValue("code") != g.code || g.code == "",
		r.PostFormValue("redirect_uri") != testRedirect,
		b64.EncodeToString(sum[:]) != g.challenge:
		http.Error(w, `{"error":"invalid_grant"}`, http.StatusBadRequest)
		return
	}
	g.code = "" // single use
	now := time.Now()
	writeTestJSON(w, map[string]any{
		"access_token": "access-token",
		"token_type":   "Bearer",
		"expires_in":   3600,
		"id_token": g.sign(map[string]any{
			"iss": g.issuer, "sub": testSubject, "aud": testClientID,
			"iat": now.Unix(), "exp": now.Add(time.Hour).Unix(),
		}),
	})
}

func (g *fakeGamma) userinfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer access-token" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	writeTestJSON(w, g.userinfo)
}

func (g *fakeGamma) sign(claims map[string]any) string {
	header, _ := json.Marshal(map[string]string{"alg": "RS256", "kid": "k1", "typ": "JWT"})
	body, _ := json.Marshal(claims)
	input := b64.EncodeToString(header) + "." + b64.EncodeToString(body)
	digest := sha256.Sum256([]byte(input))
	sig, err := rsa.SignPKCS1v15(nil, g.key, crypto.SHA256, digest[:])
	if err != nil {
		g.t.Fatal(err)
	}
	return input + "." + b64.EncodeToString(sig)
}

func writeTestJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

type loginEnv struct {
	gamma    *fakeGamma
	sessions *Sessions
	login    http.Handler
	callback http.Handler
}

func newLoginEnv(t *testing.T) *loginEnv {
	t.Helper()
	g := newFakeGamma(t)
	keys := testKeys(t, "k")
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	provider, err := NewRelyingParty(context.Background(), OIDCConfig{
		Issuer: g.srv.URL, ClientID: testClientID, ClientSecret: testClientSecret, RedirectURL: testRedirect,
	}, NewFlowCookieHandler(keys, false), logger)
	if err != nil {
		t.Fatal(err)
	}
	if !provider.IsPKCE() {
		t.Fatal("PKCE must be enabled")
	}
	sessions := NewSessions(keys, CookieOptions{Secure: false})
	return &loginEnv{gamma: g, sessions: sessions, login: LoginHandler(provider, sessions), callback: CallbackHandler(provider, sessions, logger)}
}

// startLogin hits the login endpoint and returns the flow cookies and the
// code/state Gamma would send back.
func (e *loginEnv) startLogin(t *testing.T) (cookies []*http.Cookie, code, state string) {
	t.Helper()
	rec := httptest.NewRecorder()
	e.login.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/auth/login", nil))
	if rec.Code != http.StatusFound {
		t.Fatalf("login: code %d", rec.Code)
	}
	cookies = rec.Result().Cookies()
	names := map[string]bool{}
	for _, c := range cookies {
		names[c.Name] = true
		if !c.HttpOnly || c.SameSite != http.SameSiteLaxMode || c.MaxAge <= 0 || c.MaxAge > int(flowMaxAge/time.Second) {
			t.Errorf("bad flow cookie attributes: %+v", c)
		}
	}
	if !names["state"] || !names["pkce"] || !names[LoginCookieName(false)] {
		t.Fatalf("expected state, pkce and login cookies, got %v", names)
	}
	code, state = e.gamma.authorize(rec.Header().Get("Location"))
	return cookies, code, state
}

func (e *loginEnv) finish(cookies []*http.Cookie, query string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodGet, "/api/auth/callback?"+query, nil)
	for _, c := range cookies {
		r.AddCookie(c)
	}
	rec := httptest.NewRecorder()
	e.callback.ServeHTTP(rec, r)
	return rec
}

func (e *loginEnv) sessionFrom(rec *httptest.ResponseRecorder) *http.Cookie {
	for _, c := range rec.Result().Cookies() {
		if c.Name == e.sessions.CookieName() && c.MaxAge > 0 {
			return c
		}
	}
	return nil
}

func TestLoginFlow(t *testing.T) {
	e := newLoginEnv(t)
	cookies, code, state := e.startLogin(t)

	rec := e.finish(cookies, url.Values{"code": {code}, "state": {state}}.Encode())
	if rec.Code != http.StatusFound || rec.Header().Get("Location") != "/" {
		t.Fatalf("callback: code=%d location=%q body=%q", rec.Code, rec.Header().Get("Location"), rec.Body.String())
	}
	sc := e.sessionFrom(rec)
	if sc == nil {
		t.Fatal("no session cookie set")
	}
	for _, c := range rec.Result().Cookies() {
		if (c.Name == "state" || c.Name == "pkce" || c.Name == LoginCookieName(false)) && c.MaxAge >= 0 {
			t.Errorf("flow cookie %q not deleted", c.Name)
		}
	}
	if strings.Contains(sc.Value, "access-token") {
		t.Error("session must not contain the access token")
	}
	u, err := e.sessions.Get(reqWith(sc))
	if err != nil {
		t.Fatal(err)
	}
	if want := (User{CID: "alice", Nickname: "Ali", Name: "Alice A", Picture: "https://example.com/a.png"}); u != want {
		t.Errorf("user = %+v, want %+v", u, want)
	}

	// The code is single use: replaying the callback fails gracefully.
	rec = e.finish(cookies, url.Values{"code": {code}, "state": {state}}.Encode())
	assertLoginFailed(t, e, rec)
}

func TestLoginFailures(t *testing.T) {
	t.Run("state mismatch", func(t *testing.T) {
		e := newLoginEnv(t)
		cookies, code, _ := e.startLogin(t)
		assertLoginFailed(t, e, e.finish(cookies, url.Values{"code": {code}, "state": {"forged"}}.Encode()))
	})
	t.Run("no flow cookies", func(t *testing.T) {
		e := newLoginEnv(t)
		_, code, state := e.startLogin(t)
		assertLoginFailed(t, e, e.finish(nil, url.Values{"code": {code}, "state": {state}}.Encode()))
	})
	t.Run("missing login cookie", func(t *testing.T) {
		e := newLoginEnv(t)
		cookies, code, state := e.startLogin(t)
		cookies = slices.DeleteFunc(cookies, func(c *http.Cookie) bool { return c.Name == LoginCookieName(false) })
		assertLoginFailed(t, e, e.finish(cookies, url.Values{"code": {code}, "state": {state}}.Encode()))
	})
	t.Run("planted state and pkce cookies", func(t *testing.T) {
		// An attacker on a sibling subdomain starts a login, then plants
		// its state/pkce cookies in the victim's browser and sends the
		// victim to the callback with the attacker's code and state. The
		// victim's own login cookie (or its absence) must stop this.
		e := newLoginEnv(t)
		victim, _, _ := e.startLogin(t)
		attacker, code, state := e.startLogin(t) // fake Gamma now expects the attacker's code and PKCE verifier
		var planted []*http.Cookie
		for _, c := range attacker {
			if c.Name == "state" || c.Name == "pkce" {
				planted = append(planted, c)
			}
		}
		for _, c := range victim {
			if c.Name == LoginCookieName(false) {
				planted = append(planted, c)
			}
		}
		assertLoginFailed(t, e, e.finish(planted, url.Values{"code": {code}, "state": {state}}.Encode()))
		// Without any login cookie of the victim's it fails as well.
		assertLoginFailed(t, e, e.finish(planted[:2], url.Values{"code": {code}, "state": {state}}.Encode()))
	})
	t.Run("provider error", func(t *testing.T) {
		e := newLoginEnv(t)
		cookies, _, state := e.startLogin(t)
		assertLoginFailed(t, e, e.finish(cookies, url.Values{"error": {"access_denied"}, "state": {state}}.Encode()))
	})
	t.Run("token endpoint error", func(t *testing.T) {
		e := newLoginEnv(t)
		cookies, code, state := e.startLogin(t)
		e.gamma.tokenErr = true
		assertLoginFailed(t, e, e.finish(cookies, url.Values{"code": {code}, "state": {state}}.Encode()))
	})
	t.Run("userinfo without cid", func(t *testing.T) {
		e := newLoginEnv(t)
		delete(e.gamma.userinfo, "cid")
		cookies, code, state := e.startLogin(t)
		assertLoginFailed(t, e, e.finish(cookies, url.Values{"code": {code}, "state": {state}}.Encode()))
	})
	t.Run("userinfo subject mismatch", func(t *testing.T) {
		e := newLoginEnv(t)
		e.gamma.userinfo["sub"] = "someone-else"
		cookies, code, state := e.startLogin(t)
		assertLoginFailed(t, e, e.finish(cookies, url.Values{"code": {code}, "state": {state}}.Encode()))
	})
}

func assertLoginFailed(t *testing.T, e *loginEnv, rec *httptest.ResponseRecorder) {
	t.Helper()
	if rec.Code != http.StatusFound || rec.Header().Get("Location") != LoginFailedPath {
		t.Errorf("code=%d location=%q body=%q", rec.Code, rec.Header().Get("Location"), rec.Body.String())
	}
	if e.sessionFrom(rec) != nil {
		t.Error("session cookie set on failed login")
	}
}

func TestDiscoveryIssuerMismatch(t *testing.T) {
	g := newFakeGamma(t)
	g.issuer = "https://evil.example.com"
	_, err := NewRelyingParty(context.Background(), OIDCConfig{
		Issuer: g.srv.URL, ClientID: testClientID, ClientSecret: testClientSecret, RedirectURL: testRedirect,
	}, NewFlowCookieHandler(testKeys(t, "k"), false), slog.New(slog.NewTextHandler(io.Discard, nil)))
	if !errors.Is(err, oidc.ErrIssuerInvalid) {
		t.Fatalf("err = %v, want ErrIssuerInvalid", err)
	}
}

func TestUserFromInfo(t *testing.T) {
	if (UserFromInfo(nil) != User{}) {
		t.Error("nil info must give zero user")
	}
	info := &oidc.UserInfo{Subject: "s"}
	info.Nickname = "nick"
	info.AppendClaims("cid", 123) // wrong type is ignored
	if u := UserFromInfo(info); u.CID != "" || u.Nickname != "nick" {
		t.Errorf("got %+v", u)
	}
}
