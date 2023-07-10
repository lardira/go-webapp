package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	Handler http.Handler
	Address string
	Name    string
}

type Endpoint struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

func InitServer(name, port string, endpoints ...Endpoint) Server {
	address := fmt.Sprint(":", port)
	server := http.NewServeMux()

	for _, endpoint := range endpoints {
		server.HandleFunc(endpoint.Path, endpoint.Handler)
	}

	return Server{
		Name:    name,
		Handler: server,
		Address: address,
	}
}

func (s *Server) EnableLogging() {
	logFileNameFormat := "logs/%s-%v-%v_log.txt"
	timeNow := time.Now()
	logFileName := fmt.Sprintf(logFileNameFormat, s.Name, timeNow.Day(), timeNow.Year())

	handler := s.Handler

	s.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logFile, err := os.OpenFile(
			logFileName,
			os.O_CREATE|os.O_APPEND|os.O_RDWR,
			0666,
		)

		if err != nil {
			panic(err)
		}

		multiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multiWriter)

		logMsg := fmt.Sprintf(
			"%v - %v",
			r.Method,
			r.RequestURI,
		)
		log.Println(logMsg)
		handler.ServeHTTP(w, r)

		logFile.Close()
	})
}

func RunServer(server *Server) error {
	log.Printf("Running server on address %v", server.Address)
	return http.ListenAndServe(server.Address, server.Handler)
}
