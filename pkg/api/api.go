package api

import (
	"github.com/gorilla/mux"
	"github.com/poncheska/iot-mousetrap/pkg/store/fake"
	hp "github.com/poncheska/iot-mousetrap/pkg/transport/http"
	"github.com/poncheska/iot-mousetrap/pkg/utils"
	"io"
	"log"
	"net/http"
	"os"
)

func Start() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	h := hp.Handler{
		Store: fake.NewFakeStore(),
		Logs:  utils.NewStringLogger(),
	}

	log.SetOutput(io.MultiWriter(os.Stdout, h.Logs))

	r := mux.NewRouter()
	r.HandleFunc("/log", h.GetLog).Methods(http.MethodGet)
	r.HandleFunc("/log/clear", h.ClearLog).Methods(http.MethodGet)
	r.HandleFunc("/org/signin", h.SignIn).Methods(http.MethodPost)
	r.HandleFunc("/org/signup", h.SignUp).Methods(http.MethodPost)
	r.HandleFunc("/mousetraps", h.AuthChecker(h.GetMousetraps)).Methods(http.MethodGet)
	r.HandleFunc("/trigger/{org}/{name}/{status}", h.Trigger).Methods(http.MethodGet)

	log.Println("Server started")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
