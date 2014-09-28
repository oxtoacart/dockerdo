package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"code.google.com/p/go-uuid/uuid"
	"github.com/golang/oauth2"
)

var (
	oauth *oauth2.Config
)

func initOAuth() {
	var err error
	oauth, err = oauth2.NewConfig(&oauth2.Options{
		ClientID:     getEnvRequired("OAUTH_CLIENT_ID"),
		ClientSecret: getEnvRequired("OAUTH_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:16000/oauth/callback",
		Scopes:       []string{"read", "write"},
	},
		"https://cloud.digitalocean.com/v1/oauth/authorize",
		"https://cloud.digitalocean.com/v1/oauth/token")
	if err != nil {
		log.Fatalf("Unable to configure OAuth: %s", err)
	}
}

func oauthCallback(resp http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	if params["error"] != nil {
		resp.Write([]byte("<html><body>"))
		resp.Write([]byte(params["error_description"][0]))
		resp.Write([]byte("</body></html>"))
		return
	}
	session := getSession(req)
	code := params["code"][0]
	oauthState := params["state"][0]
	expectedState := session.Values["oauthState"]
	if oauthState != expectedState {
		panic(fmt.Errorf("State didn't match!  Expected: %s\nGot: %s", expectedState, oauthState))
	}
	token, err := oauth.Exchange(code)
	if err != nil {
		panic(err)
	}
	session.Values["authToken"] = tokenToJson(token)
	session.Save(req, resp)
	redirect(resp, "/")
}

func authenticated(resp http.ResponseWriter, req *http.Request) bool {
	session := getSession(req)
	authToken := session.Values["authToken"]
	if authToken == nil {
		oauthState := uuid.NewRandom().String()
		session.Values["oauthState"] = oauthState
		session.Save(req, resp)
		authURL := oauth.AuthCodeURL(oauthState, "offline", "always")
		redirect(resp, authURL)
		return false
	}
	return true
}

func tokenToJson(token *oauth2.Token) string {
	bytes, _ := json.Marshal(token)
	return base64.StdEncoding.EncodeToString(bytes)
}

func tokenFromJson(tokenJson string) *oauth2.Token {
	token := &oauth2.Token{}
	jsonBytes, _ := base64.StdEncoding.DecodeString(tokenJson)
	json.Unmarshal(jsonBytes, token)
	return token
}
