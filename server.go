package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var Session *http.Server

type Handle func(http.ResponseWriter, *http.Request)

type Router struct {
	mux map[string]Handle
}

type RouterItem struct {
	Path   string
	Handle Handle
}

func (rt *Router) Add(items ...RouterItem) {
	for _, item := range items {
		rt.mux[item.Path] = item.Handle
	}
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	head := GetHeader(r.URL.Path)
	handle, ok := rt.mux[head]
	if ok {
		handle(w, r)
		return
	}
	http.NotFound(w, r)
}

func InitServer(host, port string) *Router {
	router := &Router{
		mux: make(map[string]Handle),
	}

	Session = &http.Server{
		Addr:           fmt.Sprintf("%s:%s", host, port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        router,
	}

	return router
}

func GetHeader(url string) string {
	sl := strings.Split(url, "/")
	return fmt.Sprintf("/%s", sl[1])
}

func RunServer() error {
	return Session.ListenAndServe()
}
