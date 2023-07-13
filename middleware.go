package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lardira/go-webapp/model"
	"github.com/lardira/go-webapp/utils"
	http_errors "github.com/lardira/go-webapp/utils/errors"
)

func (s *ApiServer) EnableJsonContentType() {
	userHandler := s.UserHandler
	authHandler := s.AuthHandler
	variantHandler := s.VariantHandler
	testHandler := s.TestHandler

	s.UserHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		userHandler.ServeHTTP(w, r)
	})

	s.AuthHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		authHandler.ServeHTTP(w, r)
	})

	s.VariantHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		variantHandler.ServeHTTP(w, r)
	})

	s.TestHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		testHandler.ServeHTTP(w, r)
	})
}

func (s *ApiServer) EnableCors() {
	userHandler := s.UserHandler
	authHandler := s.AuthHandler
	variantHandler := s.VariantHandler
	testHandler := s.TestHandler

	s.UserHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")

		//handling preflight
		if r.Method == "OPTIONS" {
			return
		}

		userHandler.ServeHTTP(w, r)
	})

	s.AuthHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")

		//handling preflight
		if r.Method == "OPTIONS" {
			return
		}

		authHandler.ServeHTTP(w, r)
	})

	s.VariantHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")

		//handling preflight
		if r.Method == "OPTIONS" {
			return
		}

		variantHandler.ServeHTTP(w, r)
	})

	s.TestHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")

		//handling preflight
		if r.Method == "OPTIONS" {
			return
		}

		testHandler.ServeHTTP(w, r)
	})
}

func (s *ApiServer) EnableSecurity() {
	variantHandler := s.VariantHandler
	testHandler := s.TestHandler

	s.VariantHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, password, ok := utils.ParseAuthHeader(w, r, AUTH_TYPE, AUTH_HEADER)
		if !ok {
			http_errors.NotAuthorized(w)
			return
		}

		user, err := model.GetUserByLoginAndPasword(
			GlobalConnectionPool,
			login,
			password,
		)

		if err != nil || !user.IsAuth {
			http_errors.NotAuthorized(w)
			return
		}

		variantHandler.ServeHTTP(w, r)
	})

	s.TestHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, password, ok := utils.ParseAuthHeader(w, r, AUTH_TYPE, AUTH_HEADER)
		if !ok {
			http_errors.NotAuthorized(w)
			return
		}

		user, err := model.GetUserByLoginAndPasword(
			GlobalConnectionPool,
			login,
			password,
		)

		if err != nil || !user.IsAuth {
			http_errors.NotAuthorized(w)
			return
		}

		testHandler.ServeHTTP(w, r)
	})
}

func (s *ApiServer) EnableRequestLogging() {
	userHandler := s.UserHandler
	authHandler := s.AuthHandler
	variantHandler := s.VariantHandler
	testHandler := s.TestHandler

	logFileNameFormat := "logs/%s-%v-%v_log.txt"
	timeNow := time.Now()
	logFileName := fmt.Sprintf(logFileNameFormat, "API", timeNow.Day(), timeNow.Year())

	s.UserHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logFile, err := os.OpenFile(
			logFileName,
			os.O_CREATE|os.O_APPEND|os.O_RDWR,
			0666,
		)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()

		multiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multiWriter)

		logMsg := fmt.Sprintf(
			"%v - %v",
			r.Method,
			r.RequestURI,
		)
		log.Println(logMsg)
		userHandler.ServeHTTP(w, r)
	})

	s.AuthHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logFile, err := os.OpenFile(
			logFileName,
			os.O_CREATE|os.O_APPEND|os.O_RDWR,
			0666,
		)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()

		multiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multiWriter)

		logMsg := fmt.Sprintf(
			"%v - %v",
			r.Method,
			r.RequestURI,
		)
		log.Println(logMsg)
		authHandler.ServeHTTP(w, r)
	})

	s.VariantHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logFile, err := os.OpenFile(
			logFileName,
			os.O_CREATE|os.O_APPEND|os.O_RDWR,
			0666,
		)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()

		multiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multiWriter)

		logMsg := fmt.Sprintf(
			"%v - %v",
			r.Method,
			r.RequestURI,
		)
		log.Println(logMsg)
		variantHandler.ServeHTTP(w, r)
	})

	s.TestHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logFile, err := os.OpenFile(
			logFileName,
			os.O_CREATE|os.O_APPEND|os.O_RDWR,
			0666,
		)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()

		multiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multiWriter)

		logMsg := fmt.Sprintf(
			"%v - %v",
			r.Method,
			r.RequestURI,
		)
		log.Println(logMsg)
		testHandler.ServeHTTP(w, r)
	})
}
