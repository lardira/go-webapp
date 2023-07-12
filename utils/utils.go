package utils

import (
	"net/http"
	"strings"
)

func ParseAuthHeader(w http.ResponseWriter, r *http.Request, authType, headerName string) (string, string, bool) {
	authHeader := r.Header.Get(headerName)
	if len(authHeader) < 1 {
		return "", "", false
	}

	creds := strings.Split(
		authHeader[len(authType)+1:],
		":",
	)
	login := creds[0]
	password := creds[1]
	return login, password, true
}
