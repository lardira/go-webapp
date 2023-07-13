package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lardira/go-webapp/utils"
)

const (
	STATIC_PATH = "www/build"
	AUTH_HEADER = "Authorization"
	AUTH_TYPE   = "Basic"
)

var GlobalConnectionPool *sql.DB

type ApiServer struct {
	UserHandler    http.Handler
	VariantHandler http.Handler
	AuthHandler    http.Handler
	TestHandler    http.Handler
}

type DefaultResponse struct {
	Message string `json:"message"`
}

func (s *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	switch head, r.URL.Path = utils.ShiftPath(r.URL.Path); head {
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

	//just for sending built SPA and files
	go runStaticServer(appConfig.Static.Port)

	//connecting handlers with endpoints
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
