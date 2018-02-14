package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config_file string
	debug       bool
)

// Debug transport lay"github.com/spf13/cobra"er
type debugTransport struct {
	t http.RoundTripper
}

func (d debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	log.Printf("%s", reqDump)

	resp, err := d.t.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	log.Printf("%s", respDump)
	return resp, nil
}

// Define each cluster
type Cluster struct {
	Name              string
	Short_Description string
	Description       string
	Issuer            string
	Client_Secret     string
	Client_ID         string
	K8s_Master_URI    string
	K8s_Ca_URI        string

	Verifier       *oidc.IDTokenVerifier
	Provider       *oidc.Provider
	OfflineAsScope bool
	Client         *http.Client
	Redirect_URI   string
}

// Define our configuration
type Config struct {
	Clusters []Cluster
	Listen   string

	TLS_Cert   string
	TLS_Key    string
	IDP_Ca_URI string
	Logo_Uri   string
}

// Start the app
// Do some config parsing
// Setup http-clients and oidc providers
// Define per-cluster handlers
func start_app(config Config) {

	// Config validation
	listenURL, err := url.Parse(config.Listen)
	if err != nil {
		log.Fatalf("parse listen address: %v", err)
	}

	var s struct {
		ScopesSupported []string `json:"scopes_supported"`
	}

	// Generate handlers for each cluster
	for i, _ := range config.Clusters {
		cluster := config.Clusters[i]
		if debug {
			if cluster.Client == nil {
				cluster.Client = &http.Client{
					Transport: debugTransport{http.DefaultTransport},
				}
			} else {
				cluster.Client.Transport = debugTransport{cluster.Client.Transport}
			}
		} else {
			cluster.Client = http.DefaultClient
		}

		ctx := oidc.ClientContext(context.Background(), cluster.Client)
		log.Printf("Creating new provider %s", cluster.Issuer)
		provider, err := oidc.NewProvider(ctx, cluster.Issuer)

		if err != nil {
			log.Fatalf("Failed to query provider %q: %v\n", cluster.Issuer, err)
		}

		cluster.Provider = provider

		log.Printf("Verifying client %s", cluster.Client_ID)

		verifier := provider.Verifier(&oidc.Config{ClientID: cluster.Client_ID})

		cluster.Verifier = verifier

		if err := provider.Claims(&s); err != nil {
			log.Fatalf("Failed to parse provider scopes_supported: %v", err)
		}

		if len(s.ScopesSupported) == 0 {
			// scopes_supported is a "RECOMMENDED" discovery claim, not a required
			// one. If missing, assume that the provider follows the spec and has
			// an "offline_access" scope.
			cluster.OfflineAsScope = true
		} else {
			// See if scopes_supported has the "offline_access" scope.
			cluster.OfflineAsScope = func() bool {
				for _, scope := range s.ScopesSupported {
					if scope == oidc.ScopeOfflineAccess {
						return true
					}
				}
				return false
			}()
		}

		base_redirect_uri, err := url.Parse(cluster.Redirect_URI)

		if err != nil {
			fmt.Errorf("parse redirect_uri address: %v", err)
			os.Exit(1)
		}

		// Each cluster gets a different login and callback URL
		callback_uri := path.Join(base_redirect_uri.Path)
		http.HandleFunc(callback_uri, cluster.handleCallback)
		log.Printf("Registered callback handler at: %s", callback_uri)

		login_uri := path.Join("/login", cluster.Name)
		http.HandleFunc(login_uri, cluster.handleLogin)
		log.Printf("Registered login handler at: /login/%s", cluster.Name)
	}

	// Index page
	http.HandleFunc("/", config.handleIndex)

	// Serve static html assets
	fs := http.FileServer(http.Dir("html/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Determine whether to use TLS or not
	switch listenURL.Scheme {
	case "http":
		log.Printf("listening on %s", config.Listen)
		http.ListenAndServe(listenURL.Host, nil)
	case "https":
		log.Printf("listening on %s", config.Listen)
		http.ListenAndServeTLS(listenURL.Host, config.TLS_Cert, config.TLS_Key, nil)
	default:
		fmt.Errorf("listen address %q is not using http or https", config.Listen)
	}
}

var RootCmd = &cobra.Command{
	Use:   "dex-k8s-authenticator",
	Short: "Dex Kubernetes Authenticator",
	Long:  `Dex Kubernetes Authenticator provides a web-interface to generate a kubeconfig file based on a selected Kubernetes cluster. One or more clusters can be defined in the configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {

		var config Config
		err := viper.Unmarshal(&config)
		if err != nil {
			log.Fatalf("Unable to decode configuration into struct, %v", err)
		}

		// Start the app
		start_app(config)

		// Fallback if no args specified
		cmd.HelpFunc()(cmd, args)

	},
}

// Read in config file
func initConfig() {

	if config_file != "" {
		viper.SetConfigFile(config_file)

		// get the filepath
		abs, err := filepath.Abs(config_file)
		if err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}

		// get the config name
		base := filepath.Base(abs)

		// get the path
		path := filepath.Dir(abs)

		viper.SetConfigName(strings.Split(base, ".")[0])
		viper.AddConfigPath(path)

		// Find and read the config file; Handle errors reading the config file
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}

		log.Printf("Using config file:", viper.ConfigFileUsed())
	}
}

// Initialization
func init() {
	cobra.OnInitialize(initConfig)

	viper.BindPFlags(RootCmd.Flags())
	RootCmd.Flags().StringVar(&config_file, "config", "", "./config.yml")
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")
}

// Let's go!
func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
