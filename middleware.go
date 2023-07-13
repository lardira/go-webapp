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

func setContentJson(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
}

func (s *ApiServer) EnableJsonContentType() {
	userHandler := s.UserHandler
	authHandler := s.AuthHandler
	variantHandler := s.VariantHandler
	testHandler := s.TestHandler

	s.UserHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setContentJson(&w)
		userHandler.ServeHTTP(w, r)
	})

	s.AuthHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setContentJson(&w)
		authHandler.ServeHTTP(w, r)
	})

	s.VariantHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setContentJson(&w)
		variantHandler.ServeHTTP(w, r)
	})

	s.TestHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setContentJson(&w)
		testHandler.ServeHTTP(w, r)
	})
}

func configureCorsOnRequest(w *http.ResponseWriter, r *http.Request) (shouldReturn bool) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")

	//handling preflight
	return r.Method == "OPTIONS"
}

func (s *ApiServer) EnableCors() {
	userHandler := s.UserHandler
	authHandler := s.AuthHandler
	variantHandler := s.VariantHandler
	testHandler := s.TestHandler

	s.UserHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shouldReturn := configureCorsOnRequest(&w, r)
		if shouldReturn {
			return
		}

		userHandler.ServeHTTP(w, r)
	})

	s.AuthHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shouldReturn := configureCorsOnRequest(&w, r)
		if shouldReturn {
			return
		}

		authHandler.ServeHTTP(w, r)
	})

	s.VariantHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shouldReturn := configureCorsOnRequest(&w, r)
		if shouldReturn {
			return
		}

		variantHandler.ServeHTTP(w, r)
	})

	s.TestHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shouldReturn := configureCorsOnRequest(&w, r)
		if shouldReturn {
			return
		}

		testHandler.ServeHTTP(w, r)
	})
}

func checkAuth(w http.ResponseWriter, r *http.Request) (shouldReturn bool) {
	login, password, ok := utils.ParseAuthHeader(w, r, AUTH_TYPE, AUTH_HEADER)
	if !ok {
		http_errors.NotAuthorized(w)
		return true
	}

	user, err := model.GetUserByLoginAndPasword(
		GlobalConnectionPool,
		login,
		password,
	)

	if err != nil || !user.IsAuth {
		http_errors.NotAuthorized(w)
		return true
	}
	return false
}

func (s *ApiServer) EnableSecurity() {
	variantHandler := s.VariantHandler
	testHandler := s.TestHandler

	s.VariantHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shouldReturn := checkAuth(w, r)
		if shouldReturn {
			return
		}

		variantHandler.ServeHTTP(w, r)
	})

	s.TestHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shouldReturn := checkAuth(w, r)
		if shouldReturn {
			return
		}

		testHandler.ServeHTTP(w, r)
	})
}

func logRequest(logFileName string, r *http.Request) {
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
		logRequest(logFileName, r)
		userHandler.ServeHTTP(w, r)
	})

	s.AuthHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(logFileName, r)
		authHandler.ServeHTTP(w, r)
	})

	s.VariantHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(logFileName, r)
		variantHandler.ServeHTTP(w, r)
	})

	s.TestHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(logFileName, r)
		testHandler.ServeHTTP(w, r)
	})
}
