// Package api implements the SuggestIT HTTP API.
package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"mime"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"chalmers.it/suggestit/internal/auth"
	"chalmers.it/suggestit/internal/store"
)

// Input limits (in runes). The UI is stricter.
const (
	MaxTitleLen  = 100
	MaxAuthorLen = 100
	MaxTextLen   = 10000

	// MaxBodyBytes limits request bodies.
	MaxBodyBytes = 64 << 10
)

// Store is the data access the API needs.
type Store interface {
	Ping(ctx context.Context) error
	List(ctx context.Context) ([]store.Suggestion, error)
	Get(ctx context.Context, id string) (store.Suggestion, error)
	Create(ctx context.Context, s store.NewSuggestion) (store.Suggestion, error)
	Delete(ctx context.Context, id string) error
	DeleteMany(ctx context.Context, ids []string) error
}

// Options wires the API together.
type Options struct {
	Store    Store
	Sessions *auth.Sessions
	// Login and Callback are the OIDC handlers (see auth.LoginHandler and
	// auth.CallbackHandler).
	Login    http.Handler
	Callback http.Handler
	Logger   *slog.Logger
	// TrustedOrigins are origins (scheme://host[:port]) whose unsafe
	// requests are always allowed. The site's public origin belongs here:
	// browsers without Sec-Fetch-Site are otherwise checked by comparing
	// Origin with Host, which fails if a reverse proxy rewrites Host.
	TrustedOrigins []string
}

type server struct {
	store    Store
	sessions *auth.Sessions
	logger   *slog.Logger
}

// New returns the complete HTTP handler: routes, CSRF protection for unsafe
// methods, security headers and access logging.
func New(o Options) (http.Handler, error) {
	s := &server{store: o.Store, sessions: o.Sessions, logger: o.Logger}
	authed := o.Sessions.Require

	mux := http.NewServeMux()
	mux.Handle("GET /api/auth/login", o.Login)
	mux.Handle("GET /api/auth/callback", o.Callback)
	mux.HandleFunc("POST /api/auth/logout", s.logout)
	mux.HandleFunc("GET /api/me", s.me)

	mux.Handle("GET /api/suggestions", authed(http.HandlerFunc(s.listSuggestions)))
	mux.HandleFunc("POST /api/suggestions", s.createSuggestion)
	mux.Handle("GET /api/suggestions/{id}", authed(http.HandlerFunc(s.getSuggestion)))
	mux.Handle("DELETE /api/suggestions/{id}", authed(http.HandlerFunc(s.deleteSuggestion)))
	mux.Handle("POST /api/suggestions/delete", authed(http.HandlerFunc(s.deleteSuggestions)))

	mux.HandleFunc("GET /healthz", s.healthz)
	// JSON 404/405 for everything else under /api/. This catch-all matches
	// every method, so it must answer 405 itself for known paths.
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		if allow := allowedMethods(mux, r); len(allow) > 0 {
			w.Header().Set("Allow", strings.Join(allow, ", "))
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		writeError(w, http.StatusNotFound, "not found")
	})

	cop := http.NewCrossOriginProtection()
	for _, origin := range o.TrustedOrigins {
		if err := cop.AddTrustedOrigin(origin); err != nil {
			return nil, err
		}
	}
	cop.SetDenyHandler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusForbidden, "cross-origin request rejected")
	}))

	return s.logRequests(securityHeaders(cop.Handler(mux))), nil
}

// probeMethods are the methods checked when building an Allow header.
var probeMethods = []string{
	http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut,
	http.MethodPatch, http.MethodDelete, http.MethodOptions,
}

// allowedMethods returns the methods for which r's path is routed to a real
// endpoint, i.e. to a pattern other than the "/api/" catch-all.
func allowedMethods(mux *http.ServeMux, r *http.Request) []string {
	var allow []string
	for _, m := range probeMethods {
		if m == r.Method {
			continue
		}
		probe := &http.Request{Method: m, Host: r.Host, URL: r.URL, Header: http.Header{}}
		if _, pattern := mux.Handler(probe); pattern != "" && pattern != "/api/" {
			allow = append(allow, m)
		}
	}
	return allow
}

func (s *server) me(w http.ResponseWriter, r *http.Request) {
	u, ok := s.sessions.Current(w, r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "not logged in")
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func (s *server) logout(w http.ResponseWriter, _ *http.Request) {
	s.sessions.Clear(w)
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) listSuggestions(w http.ResponseWriter, r *http.Request) {
	list, err := s.store.List(r.Context())
	if err != nil {
		s.internalError(w, r, "list suggestions", err)
		return
	}
	if list == nil {
		list = []store.Suggestion{}
	}
	writeJSON(w, http.StatusOK, list)
}

type suggestionInput struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	Author string `json:"author"`
}

// validate trims and checks the input. It returns a user-facing message on
// failure.
func (in suggestionInput) validate() (store.NewSuggestion, string) {
	n := store.NewSuggestion{
		Title:  strings.TrimSpace(in.Title),
		Text:   strings.TrimSpace(in.Text),
		Author: strings.TrimSpace(in.Author),
	}
	if n.Author == "" || strings.EqualFold(n.Author, store.AnonymousAuthor) {
		n.Author = store.AnonymousAuthor
	}
	for _, f := range []struct {
		name, val string
		max       int
	}{
		{"title", n.Title, MaxTitleLen},
		{"text", n.Text, MaxTextLen},
		{"author", n.Author, MaxAuthorLen},
	} {
		switch {
		case f.val == "":
			return n, f.name + " is required"
		case utf8.RuneCountInString(f.val) > f.max:
			return n, f.name + " is too long"
		case strings.ContainsRune(f.val, 0):
			return n, f.name + " contains invalid characters"
		}
	}
	return n, ""
}

func (s *server) createSuggestion(w http.ResponseWriter, r *http.Request) {
	var in suggestionInput
	if !decodeJSON(w, r, &in) {
		return
	}
	n, msg := in.validate()
	if msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}
	created, err := s.store.Create(r.Context(), n)
	if err != nil {
		s.internalError(w, r, "create suggestion", err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (s *server) getSuggestion(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !store.ValidUUID(id) {
		writeError(w, http.StatusNotFound, "suggestion not found")
		return
	}
	sug, err := s.store.Get(r.Context(), id)
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, "suggestion not found")
	case err != nil:
		s.internalError(w, r, "get suggestion", err)
	default:
		writeJSON(w, http.StatusOK, sug)
	}
}

func (s *server) deleteSuggestion(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !store.ValidUUID(id) {
		writeError(w, http.StatusNotFound, "suggestion not found")
		return
	}
	err := s.store.Delete(r.Context(), id)
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, "suggestion not found")
	case err != nil:
		s.internalError(w, r, "delete suggestion", err)
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *server) deleteSuggestions(w http.ResponseWriter, r *http.Request) {
	var in struct {
		IDs []string `json:"ids"`
	}
	if !decodeJSON(w, r, &in) {
		return
	}
	for _, id := range in.IDs {
		if !store.ValidUUID(id) {
			writeError(w, http.StatusBadRequest, "ids must be UUIDs")
			return
		}
	}
	if err := s.store.DeleteMany(r.Context(), in.IDs); err != nil {
		s.internalError(w, r, "delete suggestions", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) healthz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if err := s.store.Ping(ctx); err != nil {
		s.logger.WarnContext(r.Context(), "health check failed", "error", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("database unavailable\n"))
		return
	}
	_, _ = w.Write([]byte("ok\n"))
}

func (s *server) internalError(w http.ResponseWriter, r *http.Request, op string, err error) {
	s.logger.ErrorContext(r.Context(), op+" failed", "error", err)
	writeError(w, http.StatusInternalServerError, "internal server error")
}

// decodeJSON reads a single JSON object from a size-limited body. It writes
// an error response and returns false on failure.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type")); err != nil || mt != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return false
	}
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, MaxBodyBytes))
	if err := dec.Decode(dst); err != nil {
		if mbe := (*http.MaxBytesError)(nil); errors.As(err, &mbe) {
			writeError(w, http.StatusRequestEntityTooLarge, "request body too large")
		} else {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
		}
		return false
	}
	if dec.More() {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	if r.status == 0 {
		r.status = code
	}
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(b)
}

func (r *statusRecorder) Unwrap() http.ResponseWriter { return r.ResponseWriter }

// logRequests logs method, path (never the query, which may hold OAuth
// codes), status and duration.
func (s *server) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(rec, r)
		if rec.status == 0 {
			rec.status = http.StatusOK
		}
		level := slog.LevelInfo
		if r.URL.Path == "/healthz" && rec.status == http.StatusOK {
			level = slog.LevelDebug
		}
		s.logger.Log(r.Context(), level, "request",
			"method", r.Method, "path", r.URL.Path, "status", rec.status,
			"duration", time.Since(start).Round(time.Microsecond))
	})
}
