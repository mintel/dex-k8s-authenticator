package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func Test_DefaultState(t *testing.T) {
	_, state, _ := State{}.Create()
	assert.NotEmpty(t, state)
	_, err := State{}.Validate(nil)
	assert.NoError(t, err)
}

func Test_StateStore(t *testing.T) {
	store := MapStateStore{}

	// setup a fake resource server to handle the final exchange
	type tokenJSON struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
	}
	var challenge = "not received"
	resource := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.Path, "/token")
		assert.Equal(t, r.Method, "POST")
		ver := r.FormValue("code_verifier")
		w.Header().Add("Content-Type", "application/json")
		// does the hash of the code_verifier match the original hash we were given?
		sum256 := sha256.Sum256([]byte(ver))
		assert.Equal(t, challenge, base64.RawURLEncoding.EncodeToString(sum256[:]), ver)
		_ = json.NewEncoder(w).Encode(tokenJSON{
			AccessToken:  "a",
			TokenType:    "Bearer",
			RefreshToken: "r",
		})
	}))
	cfg := oauth2.Config{
		ClientID: "client_id",
		Endpoint: oauth2.Endpoint{
			AuthURL:   resource.URL + "/auth",
			TokenURL:  resource.URL + "/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
		Scopes: []string{"openid", "email", "profile"},
	}

	// create the auth URL and verify the parameters
	cookie, authUrl := AuthUrl(store, &cfg)
	parsed, err := url.Parse(authUrl)
	assert.NoError(t, err)
	challenge = parsed.Query().Get("code_challenge")
	assert.Equal(t, "S256", parsed.Query().Get("code_challenge_method"))

	// complete the exchange, letting the mock server verify the code
	req, err := http.NewRequest("GET", "https://fake/redirect", nil)
	assert.NoError(t, err)
	req.AddCookie(cookie)
	req.Form = url.Values{
		"state": []string{parsed.Query().Get("state")},
	}
	token, err := CompleteExchange(context.Background(), req, store, &cfg, "code")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "a", token.AccessToken)
}
