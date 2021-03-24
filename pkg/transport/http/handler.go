package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/poncheska/iot-mousetrap/pkg/models"
	"github.com/poncheska/iot-mousetrap/pkg/store"
	"github.com/poncheska/iot-mousetrap/pkg/utils"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	Store store.Store
	Logs  *utils.StringLogger
}

func (h Handler) GetMousetraps(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func (h Handler) Trigger(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	orgName := vars["org"]
	tm := time.Now()
	log.Printf("mousetrap %v/%v triggered: %v", orgName, name, tm)
}

func (h Handler) GetLog(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.Logs.Logs))
}

func (h Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func (h Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var cred *models.Credentials
	err := json.NewDecoder(r.Body).Decode(cred)
	if err != nil {
		log.Printf("signup: error: %v", err)
	}
}

func (h Handler) ClearLog(w http.ResponseWriter, r *http.Request) {
	h.Logs.Clear()
	w.Write([]byte(h.Logs.Logs))
}
