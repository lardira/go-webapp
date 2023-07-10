package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

const API_VERSION = "v1"
const STATIC_PATH = "www/build"

func apiPath(path string) string {
	return fmt.Sprintf("/api/%v%v", API_VERSION, path)
}

func main() {
	var appConfig AppConfig
	err := appConfig.InitFrom("app.json")
	if err != nil {
		log.Fatal(err)
	}

	staticServer := InitServer(
		"STATIC",
		appConfig.Static.Port,
		Endpoint{
			Path: "/",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				fileServer := http.FileServer(http.Dir(STATIC_PATH))
				fileMatcher := regexp.MustCompile(`\.[a-zA-Z]*$`)

				if !fileMatcher.MatchString(r.URL.Path) {
					indexPath := fmt.Sprint(STATIC_PATH, "/index.html")
					http.ServeFile(w, r, indexPath)
				} else {
					fileServer.ServeHTTP(w, r)
				}
			},
		},
	)

	staticServer.EnableLogging()

	apiServer := InitServer(
		"API",
		appConfig.Port,
		Endpoint{
			Path: "/status",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "ok")
			},
		},
		Endpoint{
			Path: apiPath("/test"),
			Handler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "ok")
			},
		},
	)

	apiServer.EnableLogging()

	go RunServer(&staticServer)
	RunServer(&apiServer)
}
