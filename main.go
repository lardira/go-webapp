package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lardira/go-webapp/model"
	http_errors "github.com/lardira/go-webapp/utils/errors"
)

const STATIC_PATH = "www/build"

// func apiPath(path string) string {
// 	return fmt.Sprintf("/api/%v%v", API_VERSION, path)
// }

var GlobalConnectionPool *sql.DB

func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

type ApiServer struct {
	UserHandler http.Handler
	AuthHandler http.Handler
}

type UserHandler struct {
}

type AuthHandler struct {
}

func (uh *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)

	id, err := strconv.Atoi(head)
	if len(head) > 0 && err != nil {
		http_errors.BadRequest(w, err)
		return
	}

	switch r.Method {

	case http.MethodPost:
		if id != 0 {
			http_errors.BadRequest(w, errors.New("id provided"))
			return
		}

		var request model.UserRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http_errors.BadRequest(w, err)
			return
		}

		err = request.Validate()
		if err != nil {
			http_errors.BadRequest(w, err)
			return
		}

		user, err := model.CreateUser(GlobalConnectionPool, request.Login, request.Password)
		if err != nil {
			http_errors.Error(w, errors.New("could not create user"))
			return
		}

		w.WriteHeader(http.StatusCreated)

		response, err := json.Marshal(user)
		if err != nil {
			http_errors.Error(w, errors.New("error occured"))
			return
		}

		fmt.Fprint(w, string(response))
	case http.MethodPut:
		type AuthRequest struct {
			Auth bool `json:"is_auth"`
		}

		var request AuthRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http_errors.BadRequest(w, err)
			return
		}

		err = model.SetUserAuth(GlobalConnectionPool, int64(id), request.Auth)
		if err != nil {
			http_errors.Error(w, errors.New("could not update user"))
			return
		}

		fmt.Fprint(w, "ok")

	default:
		http_errors.MethodNotAllowed(w)
	}
}

func (ah *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)

	id, err := strconv.Atoi(head)
	if len(head) > 0 && err != nil {
		http_errors.BadRequest(w, err)
		return
	}

	switch r.Method {

	case http.MethodPost:
		type AuthResponse struct {
			Key string `json:"key"`
		}

		if id != 0 {
			http_errors.BadRequest(w, errors.New("id provided"))
			return
		}

		var request model.UserRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http_errors.BadRequest(w, err)
			return
		}

		err = request.Validate()
		if err != nil {
			http_errors.BadRequest(w, err)
			return
		}

		err = model.Authorize(
			GlobalConnectionPool,
			request.Login,
			request.Password,
		)

		if err != nil {
			http_errors.NotFound(w, errors.New("no user with such credentials"))
			return
		}

		w.WriteHeader(http.StatusOK)

		response, err := json.Marshal(
			AuthResponse{
				Key: fmt.Sprintf("Basic %s:%s", request.Login, request.Password),
			},
		)

		if err != nil {
			http_errors.Error(w, errors.New("could not respond"))
			return
		}

		fmt.Fprint(w, string(response))

	default:
		http_errors.MethodNotAllowed(w)
	}
}

func (s *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	switch head, r.URL.Path = ShiftPath(r.URL.Path); head {

	case "auth":
		s.AuthHandler.ServeHTTP(w, r)
	case "users":
		s.UserHandler.ServeHTTP(w, r)
	case "status":
		fmt.Fprint(w, "ok")

	default:
		http_errors.MethodNotAllowed(w)
	}
}

func (s *ApiServer) EnableJsonContentType() {
	userHandler, authHandler := s.UserHandler, s.AuthHandler

	s.UserHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		userHandler.ServeHTTP(w, r)
	})

	s.AuthHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		authHandler.ServeHTTP(w, r)
	})
}

func (s *ApiServer) EnableSecurity() {
	const HEADER_AUTH = "Authorization"
	const AUTH_TYPE = "Basic"

	userHandler := s.UserHandler

	s.UserHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(HEADER_AUTH)
		if len(authHeader) < 1 {
			http_errors.NotAuthorized(w)
			return
		}

		creds := strings.Split(
			authHeader[len(AUTH_TYPE)+1:],
			":",
		)
		login := creds[0]
		password := creds[1]

		user, err := model.GetUserByLoginAndPasword(
			GlobalConnectionPool,
			login,
			password,
		)

		if err != nil || !user.IsAuth {
			http_errors.NotAuthorized(w)
			return
		}

		userHandler.ServeHTTP(w, r)
	})
}

func (s *ApiServer) EnableRequestLogging() {
	logFileNameFormat := "logs/%s-%v-%v_log.txt"
	timeNow := time.Now()
	logFileName := fmt.Sprintf(logFileNameFormat, "API", timeNow.Day(), timeNow.Year())

	userHandler, authHandler := s.UserHandler, s.AuthHandler

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
}

func runStaticServer(port string) {
	staticServer := http.NewServeMux()

	staticServer.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fileServer := http.FileServer(http.Dir(STATIC_PATH))
		fileMatcher := regexp.MustCompile(`\.[a-zA-Z]*$`)

		if !fileMatcher.MatchString(r.URL.Path) {
			indexPath := fmt.Sprint(STATIC_PATH, "/index.html")
			http.ServeFile(w, r, indexPath)
		} else {
			fileServer.ServeHTTP(w, r)
		}
	})

	log.Printf("Running static server on port %v\n", port)

	address := fmt.Sprint(":", port)
	err := http.ListenAndServe(address, staticServer)
	log.Fatal(err)
}

func main() {
	var appConfig AppConfig
	err := appConfig.InitFrom("app.json")
	if err != nil {
		log.Fatal(err)
	}

	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		appConfig.DbConfig.User,
		appConfig.DbConfig.Password,
		appConfig.DbConfig.Host,
		appConfig.DbConfig.Port,
		appConfig.DbConfig.Name,
	)

	connPool, err := sql.Open("pgx", dbUrl)
	if err != nil {
		log.Fatalf("Could not connect to db: %v\n", err)
	}
	defer connPool.Close()

	//checking if db connection is valid
	if err := connPool.PingContext(context.Background()); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Connected to db %v", appConfig.DbConfig.Name)
		GlobalConnectionPool = connPool
	}

	//just for sending built SPA
	go runStaticServer(appConfig.Static.Port)

	apiServer := &ApiServer{
		UserHandler: new(UserHandler),
		AuthHandler: new(AuthHandler),
	}

	//middleware
	apiServer.EnableRequestLogging()
	apiServer.EnableJsonContentType()
	apiServer.EnableSecurity()

	log.Printf("Running API server on port %v\n", appConfig.Port)

	address := fmt.Sprint(":", appConfig.Port)
	http.ListenAndServe(address, apiServer)
}
