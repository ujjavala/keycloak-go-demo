package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "strings"

    oidc "github.com/coreos/go-oidc"
    "golang.org/x/oauth2"
)

func main() {
    configURL := "http://localhost:8080/auth/realms/demo"
    ctx := context.Background()
    provider, err := oidc.NewProvider(ctx, configURL)
    if err != nil {
        panic(err)
    }

    clientID := "demo-client"
    clientSecret := "fc4502aa-e216-4fa7-b21b-158d62561537"

    redirectURL := "http://localhost:8181/demo/callback"
    // Configure an OpenID Connect aware OAuth2 client.
    oauth2Config := oauth2.Config{
        ClientID:     clientID,
        ClientSecret: clientSecret,
        RedirectURL:  redirectURL,
        // Discovery returns the OAuth2 endpoints.
        Endpoint: provider.Endpoint(),
        // "openid" is a required scope for OpenID Connect flows.
        Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
    }
    state := "somestate"

    oidcConfig := &oidc.Config{
        ClientID: clientID,
    }
    verifier := provider.Verifier(oidcConfig)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        rawAccessToken := r.Header.Get("Authorization")
        if rawAccessToken == "" {
            http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
            return
        }

        parts := strings.Split(rawAccessToken, " ")
        if len(parts) != 2 {
            w.WriteHeader(400)
            return
        }
        _, err := verifier.Verify(ctx, parts[1])

        if err != nil {
            http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
            return
        }

        w.Write([]byte("Yayy!!!Verified user id_token through keycloak"))
    })

    http.HandleFunc("/demo/callback", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Query().Get("state") != state {
            http.Error(w, "state did not match", http.StatusBadRequest)
            return
        }

        oauth2Token, err := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
        if err != nil {
            http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
            return
        }
        rawIDToken, ok := oauth2Token.Extra("id_token").(string)
        if !ok {
            http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
            return
        }
        idToken, err := verifier.Verify(ctx, rawIDToken)
        if err != nil {
            http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
            return
        }

        resp := struct {
            OAuth2Token   *oauth2.Token
            IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
        }{oauth2Token, new(json.RawMessage)}

        if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        data, err := json.MarshalIndent(resp, "", "    ")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Write(data)
    })

    log.Fatal(http.ListenAndServe("localhost:8181", nil))
}
