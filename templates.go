// FIXME: Dislike this file a bit - what's the take on referencing
// viper config values (treat it as a global, or pass values around?)
package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
)

func renderIndex(w http.ResponseWriter, config *Config) {
	t, _ := template.ParseFiles("./templates/index.html")
	t.Execute(w, config)
}

type tokenTmplData struct {
	IDToken          string
	RefreshToken     string
	RedirectURL      string
	Claims           string
	Username         string
	Issuer           string
	ClusterName      string
	ShortDescription string
	ClientSecret     string
	ClientID         string
	K8sMasterURI     string
	K8sCaURI         string
	K8sCaPem         string
	IDPCaURI         string
	LogoURI          string
}

func (cluster *Cluster) renderToken(w http.ResponseWriter,
	idToken,
	refreshToken string,
	idpCaURI string,
	logoURI string,
	claims []byte) {

	var data map[string]interface{}
	err := json.Unmarshal(claims, &data)
	if err != nil {
		panic(err)
	}

	email := data["email"].(string)
	unix_username := strings.Split(email, "@")[0]

	t, _ := template.ParseFiles("./templates/kubeconfig.html")

	token_data := tokenTmplData{
		IDToken:          idToken,
		RefreshToken:     refreshToken,
		RedirectURL:      cluster.Redirect_URI,
		Claims:           string(claims),
		Username:         unix_username,
		Issuer:           data["iss"].(string),
		ClusterName:      cluster.Name,
		ShortDescription: cluster.Short_Description,
		ClientSecret:     cluster.Client_Secret,
		ClientID:         cluster.Client_ID,
		K8sMasterURI:     cluster.K8s_Master_URI,
		K8sCaURI:         cluster.K8s_Ca_URI,
		K8sCaPem:         cluster.K8s_Ca_Pem,
		IDPCaURI:         idpCaURI,
		LogoURI:          logoURI}

	t.Execute(w, token_data)

}
