package main

import (
	"net/http"
)

func home(resp http.ResponseWriter, req *http.Request) {
	if authenticated(resp, req) {
		resp.Write([]byte("<html><body>Authenticated!</body></html>"))
	}
}
