package api

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/poncheska/iot-mousetrap/docs"
	"github.com/poncheska/iot-mousetrap/pkg/store/sql"
	hp "github.com/poncheska/iot-mousetrap/pkg/transport/http"
	"github.com/poncheska/iot-mousetrap/pkg/utils"
	httpSwagger "github.com/swaggo/http-swagger"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func Start() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatalln("dsn is empty")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	h := hp.Handler{
		Store: sql.NewMySQLStore(db),
		Logs:  utils.NewStringLogger(),
		PubSub: &utils.PubSub{
			Streamers: []*utils.Streamer{},
			SMutex:    &sync.Mutex{},
		},
	}

	log.SetOutput(io.MultiWriter(os.Stdout, h.Logs))

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./front/"))))
	r.HandleFunc("/", h.MainPage).Methods(http.MethodGet)
	r.HandleFunc("/log", h.GetLog).Methods(http.MethodGet)
	r.HandleFunc("/log/clear", h.ClearLog).Methods(http.MethodGet)
	r.HandleFunc("/org/sign-in", h.SignIn).Methods(http.MethodPost)
	r.HandleFunc("/org/sign-up", h.SignUp).Methods(http.MethodPost)
	r.HandleFunc("/mousetraps", h.AuthChecker(h.GetMousetraps)).Methods(http.MethodGet)
	r.HandleFunc("/mousetraps/ws", h.AuthChecker(h.GetMousetrapsWS)).Methods(http.MethodGet)
	r.HandleFunc("/trigger/{name}/{status:[01]}", h.AuthChecker(h.Trigger)).Methods(http.MethodGet)
	r.PathPrefix("/swagger/").HandlerFunc(httpSwagger.Handler()).Methods(http.MethodGet)

	log.Println("Server started")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
