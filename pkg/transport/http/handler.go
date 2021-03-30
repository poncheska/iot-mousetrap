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
	tokenService utils.TokenService
}

func WriteJSONError(w http.ResponseWriter, errStr string, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": errStr,
	})
}

func (h Handler) GetMousetraps(w http.ResponseWriter, r *http.Request) {
	orgId, err:= strconv.ParseInt(r.Header.Get(orgIdHeader),10,64)
	if err != nil {
		log.Printf("getmousetraps: error: %v", err)
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	mt, err := h.Store.Mousetrap.GetAll(orgId)
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
	var status bool
	if vars["status"] == "0" {
		status = false
	} else if vars["status"] == "1" {
		status = true
	} else {
		log.Printf("mousetrap %v/%v triggered: %v with invalid status", orgName, name, tm)
		return
	}

	if mt, err := h.Store.Mousetrap.GetByName(name, orgName); err != nil {
		id, err := h.Store.Mousetrap.Create(models.Mousetrap{
			Name: name,
			OrgName: orgName,
			Status: status,
			LastTrigger: tm,
		})
		if err != nil {
			log.Printf("mousetrap %v/%v triggered with error: %v", orgName, name, err)
			return
		}
		log.Printf("mousetrap %v/%v created with id = %v", orgName, name, id)
		log.Printf("mousetrap %v/%v triggered: %v", orgName, name, tm)

	} else {
		mt.LastTrigger = tm
		mt.Status = status
		err := h.Store.Mousetrap.Update(mt)
		if err != nil {
			log.Printf("mousetrap %v/%v triggered with error: %v", orgName, name, err)
			return
		}
		log.Printf("mousetrap %v/%v triggered: %v", orgName, name, tm)
	}
}

func (h Handler) GetLog(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.Logs.Logs))
}

func (h Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	cred := &models.Credentials{}
	if err := json.NewDecoder(r.Body).Decode(cred); err != nil {
		log.Println("signin: request decode error: " + err.Error())
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := cred.CheckNotEmpty(); err != nil{
		log.Println("signin: " + err.Error())
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	org, err := h.Store.Organisation.GetByCredentials(cred.Name, cred.Password)
	if err != nil {
		log.Println("signin: store error: " + err.Error())
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.tokenService.CreateToken(org.Id)
	if err != nil {
		log.Println("signin: token create error: " + err.Error())
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
	if err != nil {
		log.Println("signin: response encode error: " + err.Error())
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func (h Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	cred := &models.Credentials{}
	err := json.NewDecoder(r.Body).Decode(cred)
	if err != nil {
		log.Println("signup: request decode error: " + err.Error())
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := cred.CheckNotEmpty(); err != nil{
		log.Println("signin: " + err.Error())
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.Store.Organisation.Create(models.Organisation{
		Name:     cred.Name,
		Password: cred.Password,
	})
	if err != nil {
		log.Println("signup: store error: " + err.Error())
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("signup: organisation created with id = %v", id)
}

func (h Handler) ClearLog(w http.ResponseWriter, r *http.Request) {
	h.Logs.Clear()
	w.Write([]byte("logs cleared"))
}
