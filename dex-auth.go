package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/coreos/go-oidc"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

const exampleAppState = "Vgn2lp5QnymFtLntKX5dM8k773PwcM87T4hQtiESC1q8wkUBgw5D3kH0r5qJ"

func (cluster *Cluster) oauth2Config(scopes []string) *oauth2.Config {

	return &oauth2.Config{
		ClientID:     cluster.Client_ID,
		ClientSecret: cluster.Client_Secret,
		Endpoint:     cluster.Provider.Endpoint(),
		Scopes:       scopes,
		RedirectURL:  cluster.Redirect_URI,
	}
}

func (config *Config) handleIndex(w http.ResponseWriter, r *http.Request) {

	if len(config.Clusters) == 1 && r.URL.String() == config.Web_Path_Prefix {
		http.Redirect(w, r, path.Join(config.Web_Path_Prefix, "login", config.Clusters[0].Name), http.StatusSeeOther)
	} else {
		renderIndex(w, config)
	}
}

func (cluster *Cluster) handleLogin(w http.ResponseWriter, r *http.Request) {
	var scopes []string

	scopes = append(scopes, "openid", "profile", "email", "offline_access", "groups")

	log.Printf("Handling login-uri for: %s", cluster.Name)
	authCodeURL := cluster.oauth2Config(scopes).AuthCodeURL(exampleAppState, oauth2.AccessTypeOffline)
	log.Printf("Redirecting post-loginto: %s", authCodeURL)
	http.Redirect(w, r, authCodeURL, http.StatusSeeOther)
}

func (cluster *Cluster) handleScript(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling script callback for: %s", cluster.Name)
  w.Header().Add("Content-Type", "application/octet-stream")

  tokenData := cluster.renderCredentials(w, r)

  err := textTemplates.ExecuteTemplate(w, "scripttemplate", tokenData)

	if err != nil {
		log.Fatal(err)
	}
}

func (cluster *Cluster) handleCallback(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling callback for: %s", cluster.Name)

  tokenData := cluster.renderCredentials(w, r)

  err := templates.ExecuteTemplate(w, "kubeconfig.html", tokenData)

	if err != nil {
		log.Fatal(err)
	}
}

func (cluster *Cluster) renderCredentials(w http.ResponseWriter, r *http.Request) templateData {
	var (
		err      error
		token    *oauth2.Token
		IdpCaPem string
	)

	ctx := oidc.ClientContext(r.Context(), cluster.Client)
	oauth2Config := cluster.oauth2Config(nil)
	switch r.Method {
	case "GET":
		// Authorization redirect callback from OAuth2 auth flow.
		if errMsg := r.FormValue("error"); errMsg != "" {
			http.Error(w, errMsg+": "+r.FormValue("error_description"), http.StatusBadRequest)
			return templateData{}
		}
		code := r.FormValue("code")
		if code == "" {
			http.Error(w, fmt.Sprintf("No code in request: %q", r.Form), http.StatusBadRequest)
			return templateData{}
		}
		if state := r.FormValue("state"); state != exampleAppState {
			http.Error(w, fmt.Sprintf("Expected state %q got %q", exampleAppState, state), http.StatusBadRequest)
			return templateData{}
		}
		token, err = oauth2Config.Exchange(ctx, code)
	case "POST":
		// Form request from frontend to refresh a token.
		refresh := r.FormValue("refresh_token")
		if refresh == "" {
			http.Error(w, fmt.Sprintf("No refresh_token in request: %q", r.Form), http.StatusBadRequest)
			return templateData{}
		}
		t := &oauth2.Token{
			RefreshToken: refresh,
			Expiry:       time.Now().Add(-time.Hour),
		}
		token, err = oauth2Config.TokenSource(ctx, t).Token()
	default:
		http.Error(w, fmt.Sprintf("Method not implemented: %s", r.Method), http.StatusBadRequest)
		return templateData{}
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get token: %v", err), http.StatusInternalServerError)
		return templateData{}
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token in token response", http.StatusInternalServerError)
		return templateData{}
	}

	idToken, err := cluster.Verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to verify ID token: %v", err), http.StatusInternalServerError)
		return templateData{}
	}
	var claims json.RawMessage
	idToken.Claims(&claims)

	buff := new(bytes.Buffer)
	json.Indent(buff, []byte(claims), "", "  ")

	if cluster.Config.IDP_Ca_Pem != "" {
		IdpCaPem = cluster.Config.IDP_Ca_Pem
	} else if cluster.Config.IDP_Ca_Pem_File != "" {
		content, err := ioutil.ReadFile(cluster.Config.IDP_Ca_Pem_File)
		if err != nil {
			log.Fatalf("Failed to load CA from file %s, %s", cluster.Config.IDP_Ca_Pem_File, err)
		}
		IdpCaPem = cast.ToString(content)
	}

  // rawIDToken
	refreshToken,
		idpCaURI,
		idpCaPem,
		logoURI,
		webPathPrefix,
		kubectlVersion,
		claims := token.RefreshToken,
		cluster.Config.IDP_Ca_URI,
		IdpCaPem,
		cluster.Config.Logo_Uri,
		cluster.Config.Web_Path_Prefix,
		viper.GetString("kubectl_version"),
		buff.Bytes()

	var data map[string]interface{}
	err = json.Unmarshal(claims, &data)
	if err != nil {
		panic(err)
	}

	unix_username := "user"
	if data["email"] != nil {
		email := data["email"].(string)
		unix_username = strings.Split(email, "@")[0]
	}

	return templateData{
		IDToken:           rawIDToken,
		RefreshToken:      refreshToken,
		RedirectURL:       cluster.Redirect_URI,
		Claims:            string(claims),
		Username:          unix_username,
		Issuer:            data["iss"].(string),
		ClusterName:       cluster.Name,
		ShortDescription:  cluster.Short_Description,
		ClientSecret:      cluster.Client_Secret,
		ClientID:          cluster.Client_ID,
		K8sMasterURI:      cluster.K8s_Master_URI,
		K8sCaURI:          cluster.K8s_Ca_URI,
		K8sCaPem:          cluster.K8s_Ca_Pem,
		IDPCaURI:          idpCaURI,
		IDPCaPem:          idpCaPem,
		LogoURI:           logoURI,
		Web_Path_Prefix:   webPathPrefix,
		StaticContextName: cluster.Static_Context_Name,
		KubectlVersion:    kubectlVersion}
}
