package utils

import (
	"net/http"
	"path"
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

func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
