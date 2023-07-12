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
	"github.com/lardira/go-webapp/utils"
	http_errors "github.com/lardira/go-webapp/utils/errors"
)

const (
	STATIC_PATH = "www/build"
	HEADER_AUTH = "Authorization"
	AUTH_TYPE   = "Basic"
)

var GlobalConnectionPool *sql.DB

func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

type DefaultResponse struct {
	Message string `json:"message"`
}

type ApiServer struct {
	UserHandler    http.Handler
	VariantHandler http.Handler
	AuthHandler    http.Handler
	TestHandler    http.Handler
}

type UserHandler struct {
}

type VariantHandler struct {
}

type AuthHandler struct {
}

type TestHandler struct {
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

	default:
		http_errors.MethodNotAllowed(w)
	}
}

func returAvailableTasks(w http.ResponseWriter, id int) {
	tasks, err := model.GetAllTasksByVariantId(GlobalConnectionPool, int64(id))
	if err != nil {
		http_errors.Error(w, err)
		return
	}

	response, err := json.Marshal(tasks)
	if err != nil {
		http_errors.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(response))
}

func returnAllVariants(w http.ResponseWriter) {
	variants, err := model.GetAllVariants(GlobalConnectionPool)
	if err != nil {
		http_errors.Error(w, err)
		return
	}

	response, err := json.Marshal(variants)
	if err != nil {
		http_errors.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(response))
}

func returnTaskById(w http.ResponseWriter, taskId int, varId int) {
	task, err := model.GetTask(GlobalConnectionPool, int64(taskId), int64(varId))
	if err != nil {
		http_errors.Error(w, err)
		return
	}

	response, err := json.Marshal(task)
	if err != nil {
		http_errors.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(response))
}

func (uh *VariantHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)

	varId, err := strconv.Atoi(head)
	if len(head) > 0 && err != nil {
		http_errors.BadRequest(w, err)
		return
	}

	variantIdPresented := len(head) > 0

	switch r.Method {

	case http.MethodGet:
		if variantIdPresented {
			head, r.URL.Path = ShiftPath(r.URL.Path)

			//return variant's task data if id of task
			taskIdPresented := len(head) > 0
			if taskIdPresented {
				taskId, err := strconv.Atoi(head)
				if len(head) > 0 && err != nil {
					http_errors.BadRequest(w, err)
					return
				}
				returnTaskById(w, taskId, varId)

			} else {
				returAvailableTasks(w, varId)
			}

		} else {
			returnAllVariants(w)
		}

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

		response, err := json.Marshal(
			AuthResponse{
				Key: fmt.Sprintf("Basic %s:%s", request.Login, request.Password),
			},
		)

		if err != nil {
			http_errors.Error(w, errors.New("could not respond"))
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(response))

	case http.MethodPut:

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

		err = model.LogOutUser(
			GlobalConnectionPool,
			request.Login,
			request.Password,
		)
		if err != nil {
			http_errors.NotFound(w, errors.New("no user with such credentials"))
			return
		}

		w.WriteHeader(http.StatusOK)

		response, _ := json.Marshal(DefaultResponse{
			Message: "ok",
		})
		fmt.Fprint(w, string(response))

	default:
		http_errors.MethodNotAllowed(w)
	}
}

func (th *TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)

	varId, err := strconv.Atoi(head)
	if len(head) > 0 && err != nil {
		http_errors.BadRequest(w, err)
		return
	}

	idPresented := len(head) > 0

	switch r.Method {

	case http.MethodPost:
		if idPresented {
			login, password, ok := utils.ParseAuthHeader(w, r, AUTH_TYPE, HEADER_AUTH)
			if !ok {
				http_errors.BadRequest(w, err)
				return
			}

			user, err := model.GetUserByLoginAndPasword(GlobalConnectionPool, login, password)
			if err != nil {
				http_errors.BadRequest(w, err)
				return
			}

			test, err := model.CreateTest(GlobalConnectionPool, user.Id, int64(varId))
			if err != nil {
				http_errors.Error(w, err)
				return
			}

			response, err := json.Marshal(model.TestResponse{Id: test.Id})
			if err != nil {
				http_errors.Error(w, err)
				return
			}

			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, string(response))
		}

	case http.MethodPut:
		if idPresented {
			var request model.TestAnswerRequest

			err := json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				http_errors.BadRequest(w, err)
				return
			}

			err = model.AddTestAnswer(
				GlobalConnectionPool,
				request.TestId,
				request.Answer,
			)

			if err != nil {
				http_errors.Error(w, err)
				return
			}

			w.WriteHeader(http.StatusCreated)
			response, _ := json.Marshal(DefaultResponse{Message: "ok"})
			fmt.Fprint(w, string(response))
		}

	case http.MethodGet:
		if idPresented {
			head, r.URL.Path = ShiftPath(r.URL.Path)

			testId, err := strconv.Atoi(head)
			if len(head) > 0 && err != nil {
				http_errors.BadRequest(w, err)
				return
			}

			testResult, err := model.GetTestResult(
				GlobalConnectionPool,
				int64(testId),
				int64(varId),
			)

			if err != nil {
				http_errors.Error(w, err)
				return
			}

			response, err := json.Marshal(testResult)
			if err != nil {
				http_errors.Error(w, err)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, string(response))
		}

	default:
		http_errors.MethodNotAllowed(w)
	}
}

func (s *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string

	switch head, r.URL.Path = ShiftPath(r.URL.Path); head {
	case "status":
		fmt.Fprint(w, "ok")
	case "auth":
		s.AuthHandler.ServeHTTP(w, r)
	case "users":
		s.UserHandler.ServeHTTP(w, r)
	case "variants":
		s.VariantHandler.ServeHTTP(w, r)
	case "tests":
		s.TestHandler.ServeHTTP(w, r)
	default:
		http.NotFound(w, r)
	}
}

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
		login, password, ok := utils.ParseAuthHeader(w, r, AUTH_TYPE, HEADER_AUTH)
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
		login, password, ok := utils.ParseAuthHeader(w, r, AUTH_TYPE, HEADER_AUTH)
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
		UserHandler:    new(UserHandler),
		AuthHandler:    new(AuthHandler),
		VariantHandler: new(VariantHandler),
		TestHandler:    new(TestHandler),
	}

	//middleware
	apiServer.EnableRequestLogging()
	apiServer.EnableJsonContentType()
	apiServer.EnableSecurity()
	apiServer.EnableCors()

	log.Printf("Running API server on port %v\n", appConfig.Port)

	address := fmt.Sprint(":", appConfig.Port)
	http.ListenAndServe(address, apiServer)
}
