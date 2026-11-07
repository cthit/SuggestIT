// Package auth implements Gamma (OIDC) login and stateless, encrypted
// cookie sessions.
package auth

import (
	"context"
	"crypto/hkdf"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	httphelper "github.com/zitadel/oidc/v3/pkg/http"
)

const (
	// SessionMaxAge is how long a login lasts.
	SessionMaxAge = 24 * time.Hour
	// flowMaxAge bounds how long the state/PKCE cookies of a login attempt
	// stay valid.
	flowMaxAge = 15 * time.Minute

	minSecretLen = 32
	keyLen       = 32
	hashKeyInfo  = "suggestit cookie hmac-sha256 key v1"
	encKeyInfo   = "suggestit cookie aes-256 key v1"
)

// ErrNoSession is returned when a request carries no valid session.
var ErrNoSession = errors.New("no valid session")

// Keys are the cookie keys derived from SESSION_SECRET.
type Keys struct {
	Hash    []byte // HMAC-SHA256 authentication key
	Encrypt []byte // AES-256 encryption key
}

// DeriveKeys derives independent HMAC and AES keys from the session secret
// using HKDF-SHA256 with distinct info labels.
func DeriveKeys(secret []byte) (Keys, error) {
	if len(secret) < minSecretLen {
		return Keys{}, fmt.Errorf("session secret must be at least %d bytes", minSecretLen)
	}
	hashKey, err := hkdf.Key(sha256.New, secret, nil, hashKeyInfo, keyLen)
	if err != nil {
		return Keys{}, err
	}
	encKey, err := hkdf.Key(sha256.New, secret, nil, encKeyInfo, keyLen)
	if err != nil {
		return Keys{}, err
	}
	return Keys{Hash: hashKey, Encrypt: encKey}, nil
}

// CookieOptions controls cookie attributes shared by all auth cookies.
type CookieOptions struct {
	Secure bool
	Domain string // empty: host-only cookie
}

func handlerOpts(secure bool, domain string, maxAge time.Duration) []httphelper.CookieHandlerOpt {
	opts := []httphelper.CookieHandlerOpt{
		httphelper.WithPath("/"),
		httphelper.WithSameSite(http.SameSiteLaxMode),
		httphelper.WithMaxAge(int(maxAge / time.Second)),
	}
	if domain != "" {
		opts = append(opts, httphelper.WithDomain(domain))
	}
	if !secure {
		opts = append(opts, httphelper.WithUnsecure())
	}
	return opts
}

// NewFlowCookieHandler returns the cookie handler used by the relying party
// for its short-lived state and PKCE verifier cookies. They are always
// host-only.
func NewFlowCookieHandler(keys Keys, secure bool) *httphelper.CookieHandler {
	return httphelper.NewCookieHandler(keys.Hash, keys.Encrypt, handlerOpts(secure, "", flowMaxAge)...)
}

// SessionCookieName returns the session cookie name. The __Host- prefix is
// used whenever the browser allows it (Secure, Path=/, no Domain).
func SessionCookieName(o CookieOptions) string {
	if o.Secure && o.Domain == "" {
		return "__Host-suggestit_session"
	}
	return "suggestit_session"
}

// User is the logged-in Gamma user, as exposed by GET /api/me.
type User struct {
	CID      string `json:"cid"`
	Nickname string `json:"nickname"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
}

// payload is what is stored (encrypted) in the session cookie. It never
// contains tokens.
type payload struct {
	User
	IssuedAt int64 `json:"iat"`
}

// LoginCookieName returns the name of the cookie that binds a login attempt
// to the browser that started it (see LoginHandler). It is host-only and
// uses the __Host- prefix when Secure, so unlike rp's fixed "state" and
// "pkce" cookies it cannot be planted from a sibling subdomain.
func LoginCookieName(secure bool) string {
	if secure {
		return "__Host-suggestit_login"
	}
	return "suggestit_login"
}

// Sessions reads and writes the encrypted, authenticated session cookie.
type Sessions struct {
	cookies   *httphelper.CookieHandler
	name      string
	loginName string
	maxAge    time.Duration
	now       func() time.Time
}

// NewSessions creates a session manager. The cookie is HttpOnly, SameSite=Lax
// and expires after SessionMaxAge both in the browser and server side.
func NewSessions(keys Keys, o CookieOptions) *Sessions {
	return &Sessions{
		cookies:   httphelper.NewCookieHandler(keys.Hash, keys.Encrypt, handlerOpts(o.Secure, o.Domain, SessionMaxAge)...),
		name:      SessionCookieName(o),
		loginName: LoginCookieName(o.Secure),
		maxAge:    SessionMaxAge,
		now:       time.Now,
	}
}

// CookieName returns the name of the session cookie.
func (s *Sessions) CookieName() string { return s.name }

// Set starts a session for u.
func (s *Sessions) Set(w http.ResponseWriter, u User) error {
	if u.CID == "" {
		return errors.New("user has no cid")
	}
	b, err := json.Marshal(payload{User: u, IssuedAt: s.now().Unix()})
	if err != nil {
		return err
	}
	return s.cookies.SetCookie(w, s.name, string(b))
}

// Get returns the user of the request's session. It returns ErrNoSession if
// there is no cookie, or if it is invalid, tampered with or expired.
func (s *Sessions) Get(r *http.Request) (User, error) {
	raw, err := s.cookies.CheckCookie(r, s.name)
	if err != nil {
		return User{}, ErrNoSession
	}
	var p payload
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		return User{}, ErrNoSession
	}
	iat := time.Unix(p.IssuedAt, 0)
	now := s.now()
	if p.CID == "" || now.Sub(iat) > s.maxAge || iat.After(now.Add(time.Minute)) {
		return User{}, ErrNoSession
	}
	return p.User, nil
}

// Clear deletes the session cookie.
func (s *Sessions) Clear(w http.ResponseWriter) {
	s.cookies.DeleteCookie(w, s.name)
}

// HasCookie reports whether the request carries a session cookie at all
// (valid or not).
func (s *Sessions) HasCookie(r *http.Request) bool {
	_, err := r.Cookie(s.name)
	return err == nil
}

// Current returns the session user and deletes the cookie if it is present
// but invalid, so stale cookies do not linger.
func (s *Sessions) Current(w http.ResponseWriter, r *http.Request) (User, bool) {
	u, err := s.Get(r)
	if err != nil {
		if s.HasCookie(r) {
			s.Clear(w)
		}
		return User{}, false
	}
	return u, true
}

type userKey struct{}

// UserFrom returns the user stored in ctx by Require.
func UserFrom(ctx context.Context) (User, bool) {
	u, ok := ctx.Value(userKey{}).(User)
	return u, ok
}

// Require only lets requests with a valid session through and responds
// 401 {"error": ...} otherwise. The user is available via UserFrom.
func (s *Sessions) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, ok := s.Current(w, r)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"not logged in"}` + "\n"))
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userKey{}, u)))
	})
}
