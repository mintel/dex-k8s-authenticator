package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type authState struct {
	state string
	pkce  string
}

type State struct {
	Store StateStore
}

func (s State) Create() (cookie *http.Cookie, state, pkceHash string) {
	if s.Store == nil {
		return nil, defaultState, ""
	}
	return s.Store.Create()
}

func (s State) Validate(r *http.Request) (codeVerifier string, validationError error) {
	if s.Store == nil {
		return "", nil
	}
	return s.Store.Validate(r)
}

// defaultState is used when no state store is provided
// not ideal that it is constant for the life of the server, but better than nothing
var defaultState = uuid.NewString()

var _ StateStore = State{}

type MapStateStore map[string]authState

type StateStore interface {
	Create() (cookie *http.Cookie, state, pkceHash string)
	Validate(r *http.Request) (codeVerifier string, validationError error)
}

func AuthUrl(state StateStore, cfg *oauth2.Config) (cookie *http.Cookie, url string) {
	cookie, nonce, pkce := state.Create()
	opts := []oauth2.AuthCodeOption{oauth2.AccessTypeOffline}
	if cookie != nil {
		opts = append(opts, oauth2.SetAuthURLParam("code_challenge_method", "S256"),
			oauth2.SetAuthURLParam("code_challenge", pkce))
	}
	return cookie, cfg.AuthCodeURL(nonce, opts...)
}

func CompleteExchange(ctx context.Context, r *http.Request, state StateStore, cfg *oauth2.Config, code string) (*oauth2.Token, error) {
	codeVerifier, err := state.Validate(r)
	if err != nil {
		return nil, fmt.Errorf("handleCallback: invalid state: %w", err)
	}
	var opts []oauth2.AuthCodeOption
	if codeVerifier != "" {
		opts = append(opts, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	}
	return cfg.Exchange(ctx, code, opts...)
}

const cookieName = "dex-k8s-auth-state"

// Create generates unique state and pkce parameters and sets a cookie to reference them
func (states MapStateStore) Create() (cookie *http.Cookie, state, pkceHash string) {
	key, state := uuid.NewString(), uuid.NewString()

	// compute the PKCE challenge
	{
		var pk [32]byte
		_, _ = rand.Read(pk[:])
		// store the original string
		original := base64.RawURLEncoding.EncodeToString(pk[:])
		states[key] = authState{
			state: state,
			pkce:  original,
		}
		// return the hash
		sum256 := sha256.Sum256([]byte(original))
		pkceHash = base64.RawURLEncoding.EncodeToString(sum256[:])
	}

	// TODO this should be configurable
	pathPrefix := "/"

	return &http.Cookie{
		Name:     cookieName,
		Value:    key,
		HttpOnly: true,
		// Lax is required since a redirect will return us to the authenticator
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		MaxAge:   300,
		Path:     pathPrefix,
	}, state, pkceHash
}

// Validate checks the state parameter against the cookie
func (states MapStateStore) Validate(r *http.Request) (codeVerifier string, _ error) {

	stateKey, err := r.Cookie(cookieName)
	if err != nil {
		return "", fmt.Errorf("no state cookie found")
	}
	expectedState, set := states[stateKey.Value]
	if !set {
		return "", fmt.Errorf("handleCallback: invalid state cookie: %q", stateKey.Value)
	}
	// single use
	delete(states, stateKey.Value)

	if state := r.FormValue("state"); state != expectedState.state {
		return "", fmt.Errorf("handleCallback: expected state %q got %q", expectedState, state)
	}
	return expectedState.pkce, nil
}
