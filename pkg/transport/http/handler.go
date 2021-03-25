package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/poncheska/iot-mousetrap/pkg/models"
	"github.com/poncheska/iot-mousetrap/pkg/store"
	"github.com/poncheska/iot-mousetrap/pkg/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	Store store.Store
	Logs  *utils.StringLogger
}

func WriteJSONError(w http.ResponseWriter, errStr string, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": errStr,
	})
}

func (h Handler) GetMousetraps(w http.ResponseWriter, r *http.Request) {
	//TODO заменить на чтение из хэдера после авторизации
	orgId, err := strconv.Atoi(r.URL.Query()["orgId"][0])
	if err != nil {
		log.Printf("getmousetraps: error: %v", err)
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	mt, err := h.Store.Mousetrap.GetAll(int64(orgId))
	if err != nil {
		log.Printf("getmousetraps: error: %v", err)
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(mt) == 0 {
		log.Println("getmousetraps: no such mousetrap")
		WriteJSONError(w, "no such mousetrap", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(mt)
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
	var cred models.Credentials
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		log.Printf("signup: error: %v", err)
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.Store.Organisation.Create(models.Organisation{
		Name:     cred.Name,
		Password: cred.Password,
	})
	if err != nil {
		log.Printf("signup: error: %v", err)
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h Handler) ClearLog(w http.ResponseWriter, r *http.Request) {
	h.Logs.Clear()
	w.Write([]byte("logs cleared"))
}
