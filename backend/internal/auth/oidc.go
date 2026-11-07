package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"log/slog"
	"net/http"
	"time"

	"github.com/zitadel/oidc/v3/pkg/client/rp"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

// LoginFailedPath is where the browser is sent when a login fails for any
// reason. The frontend shows a toast for it.
const LoginFailedPath = "/?login=failed"

// Scopes requested from Gamma.
var Scopes = []string{oidc.ScopeOpenID, oidc.ScopeProfile}

// OIDCConfig configures the Gamma relying party.
type OIDCConfig struct {
	Issuer       string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	// HTTPClient is used for all calls to the provider. Nil means a client
	// with a 10 second timeout.
	HTTPClient *http.Client
}

// NewRelyingParty runs OIDC discovery against the issuer and returns a
// relying party using PKCE and the given cookie handler for its state and
// verifier cookies. All flow errors redirect to LoginFailedPath and are
// logged; raw error text is never rendered to the user.
func NewRelyingParty(ctx context.Context, cfg OIDCConfig, cookies *httphelper.CookieHandler, logger *slog.Logger) (rp.RelyingParty, error) {
	client := cfg.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	fail := func(w http.ResponseWriter, r *http.Request, msg string, attrs ...any) {
		logger.WarnContext(r.Context(), "login failed: "+msg, attrs...)
		http.Redirect(w, r, LoginFailedPath, http.StatusFound)
	}
	return rp.NewRelyingPartyOIDC(ctx, cfg.Issuer, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL, Scopes,
		rp.WithPKCE(cookies),
		rp.WithCookieHandler(cookies),
		rp.WithHTTPClient(client),
		rp.WithLogger(logger),
		rp.WithVerifierOpts(rp.WithIssuedAtOffset(30*time.Second)),
		rp.WithErrorHandler(func(w http.ResponseWriter, r *http.Request, errorType, errorDesc, _ string) {
			fail(w, r, "provider returned an error", "error", errorType, "description", errorDesc)
		}),
		rp.WithUnauthorizedHandler(func(w http.ResponseWriter, r *http.Request, desc, _ string) {
			fail(w, r, "unauthorized", "reason", desc)
		}),
	)
}

// LoginHandler redirects to Gamma's authorization endpoint with a random
// state and a PKCE challenge (both stored in encrypted cookies by rp).
//
// rp's flow cookies have the fixed names "state" and "pkce" and cannot use
// the __Host- prefix, so a sibling *.chalmers.it host could plant its own
// (validly encrypted) values for them and log the victim into the
// attacker's account. To prevent that, the state is also stored in a
// __Host- prefixed login cookie that CallbackHandler checks first.
func LoginHandler(provider rp.RelyingParty, sessions *Sessions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state := rand.Text()
		if err := provider.CookieHandler().SetCookie(w, sessions.loginName, state); err != nil {
			slog.ErrorContext(r.Context(), "login failed: could not set login cookie", "error", err)
			http.Redirect(w, r, LoginFailedPath, http.StatusFound)
			return
		}
		rp.AuthURLHandler(func() string { return state }, provider).ServeHTTP(w, r)
	})
}

// CallbackHandler exchanges the authorization code, fetches userinfo,
// starts a session and redirects to "/". Users without a cid are rejected,
// as are callbacks whose state does not match the browser's login cookie.
func CallbackHandler(provider rp.RelyingParty, sessions *Sessions, logger *slog.Logger) http.Handler {
	onUserinfo := func(w http.ResponseWriter, r *http.Request, _ *oidc.Tokens[*oidc.IDTokenClaims], _ string, _ rp.RelyingParty, info *oidc.UserInfo) {
		u := UserFromInfo(info)
		if u.CID == "" {
			logger.WarnContext(r.Context(), "login failed: userinfo has no cid", "sub", info.Subject)
			http.Redirect(w, r, LoginFailedPath, http.StatusFound)
			return
		}
		if err := sessions.Set(w, u); err != nil {
			logger.ErrorContext(r.Context(), "login failed: could not set session cookie", "error", err)
			http.Redirect(w, r, LoginFailedPath, http.StatusFound)
			return
		}
		logger.InfoContext(r.Context(), "user logged in", "cid", u.CID)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	exchange := rp.CodeExchangeHandler(rp.UserinfoCallback(onUserinfo), provider)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies := provider.CookieHandler()
		want, err := cookies.CheckCookie(r, sessions.loginName)
		cookies.DeleteCookie(w, sessions.loginName)
		if err != nil || subtle.ConstantTimeCompare([]byte(want), []byte(r.FormValue("state"))) != 1 {
			logger.WarnContext(r.Context(), "login failed: state does not match login cookie")
			http.Redirect(w, r, LoginFailedPath, http.StatusFound)
			return
		}
		exchange.ServeHTTP(w, r)
	})
}

// UserFromInfo maps Gamma userinfo to a User. cid is a custom claim.
func UserFromInfo(info *oidc.UserInfo) User {
	if info == nil {
		return User{}
	}
	cid, _ := info.Claims["cid"].(string)
	return User{
		CID:      cid,
		Nickname: info.Nickname,
		Name:     info.Name,
		Picture:  info.Picture,
	}
}
