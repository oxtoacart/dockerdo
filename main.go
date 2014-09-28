package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	SESSION_KEY = "SESSION_KEY"
)

var (
	store = sessions.NewCookieStore([]byte(getEnvRequired("SESSION_SECRET")))
)

func main() {
	initOAuth()
	initRoutes()

	http.ListenAndServe("localhost:16000", nil)
}

func initRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/oauth/callback", oauthCallback)
	http.Handle("/", r)
}

func getSession(req *http.Request) *sessions.Session {
	session, _ := store.Get(req, SESSION_KEY)
	return session
}

func getEnvRequired(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Errorf("No environment variable %s found", key))
	}
	return val
}

func redirect(resp http.ResponseWriter, url string) {
	resp.Header().Set("Location", url)
	resp.WriteHeader(302)
}
