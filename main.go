package main

import (
	"fmt"
	"log"
	"net/http"
)

func statusHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func main() {
	var appConfig AppConfig
	err := appConfig.InitFrom("app.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Running server on port: %v", appConfig.Port)

	InitServer(appConfig.Host, appConfig.Port).Add(
		RouterItem{Path: "/status", Handle: statusHandle},
	)

	err = RunServer()
	log.Fatal(err)
}
