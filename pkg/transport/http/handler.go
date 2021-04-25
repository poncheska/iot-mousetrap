package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/poncheska/iot-mousetrap/pkg/models"
	"github.com/poncheska/iot-mousetrap/pkg/store"
	"github.com/poncheska/iot-mousetrap/pkg/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	Store        store.Store
	Logs         *utils.StringLogger
	tokenService utils.TokenService
	PubSub       *utils.PubSub
}

type errorResponse struct {
	Message string `json:"message"`
}

func WriteJSONError(w http.ResponseWriter, errStr string, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(errorResponse{
		Message: errStr,
	})
}

func (h Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadFile("./front/index.html")
	if err != nil {
		log.Fatalln(err.Error())
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(bs)
}

// @Summary Get Mousetraps
// @Security ApiKeyAuth
// @Tags organisation
// @Description get mousetraps info by org id
// @ID get-mousetraps
// @Accept  json
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /mousetraps [get]
func (h Handler) GetMousetraps(w http.ResponseWriter, r *http.Request) {
	orgId, err := strconv.ParseInt(r.Header.Get(orgIdHeader), 10, 64)
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

func (h Handler) GetMousetrapsWS(w http.ResponseWriter, r *http.Request) {
	orgId, err := strconv.ParseInt(r.Header.Get(orgIdHeader), 10, 64)
	if err != nil {
		log.Printf("getmousetrapsws: error: %v", err)
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	mt, err := h.Store.Mousetrap.GetAll(orgId)
	if err != nil {
		log.Printf("getmousetrapsws: error: %v", err)
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(mt) == 0 {
		log.Println("getmousetrapsws: no such mousetrap")
		WriteJSONError(w, "no such mousetrap", http.StatusBadRequest)
		return
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("chat: socket upgrader error: " + err.Error())
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = conn.WriteJSON(mt)
	if err != nil {
		log.Printf("getmousetrapsws: %v", err)
		return
	}

	streamer := h.PubSub.GetStreamer(orgId)

	go MousetrapsStream(conn, streamer)
}

func MousetrapsStream(conn *websocket.Conn, s *utils.Streamer) {
	s.Subscribe()
	log.Printf("add subscriber with id = %v (total %v)", s.Id, s.SubCounter)
	defer func() {
		s.Unsubscribe()
		log.Printf("del subscriber with id = %v (total %v)", s.Id, s.SubCounter)
		conn.Close()
	}()
	for {
		msg := <-s.Ch
		log.Printf("mtstreamer new message: %v", msg)
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Printf("getmousetrapsws: %v", err)
			return
		}
	}
}

// @Summary Trigger mousetrap
// @Security ApiKeyAuth
// @Tags mousetrap
// @Description update mousetrap status
// @ID trigger-mousetrap
// @Accept  json
// @Produce  json
// @Param name path string true "Mousetrap name"
// @Param status path int true "Mousetrap status 0=off 1=on"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /trigger/{name}/{status} [get]
func (h Handler) Trigger(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	tm := time.Now()
	orgId, err := strconv.ParseInt(r.Header.Get(orgIdHeader), 10, 64)
	if err != nil {
		log.Printf("getmousetraps: error: %v", err)
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	var status bool
	if vars["status"] == "0" {
		status = false
	} else if vars["status"] == "1" {
		status = true
	} else {
		log.Printf("mousetrap %v/%v triggered: %v with invalid status", orgId, name, tm)
		WriteJSONError(w, "invalid status", http.StatusBadRequest)
		return
	}

	if mt, err := h.Store.Mousetrap.GetByName(name, orgId); err != nil {
		id, err := h.Store.Mousetrap.Create(models.Mousetrap{
			Name:        name,
			OrgId:       orgId,
			Status:      status,
			LastTrigger: tm,
		})
		if err != nil {
			log.Printf("mousetrap %v/%v triggered with error: %v", orgId, name, err)
			WriteJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("mousetrap %v/%v created with id = %v", orgId, name, id)
		log.Printf("mousetrap %v/%v triggered: %v", orgId, name, tm)
	} else {
		mt.LastTrigger = tm
		mt.Status = status
		err := h.Store.Mousetrap.Update(mt)
		if err != nil {
			log.Printf("mousetrap %v/%v triggered with error: %v", orgId, name, err)
			return
		}
		log.Printf("mousetrap %v/%v triggered: %v", orgId, name, tm)
	}
	err = h.trigNotify(orgId)
	if err != nil {
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("trigger notified")
}

func (h Handler) trigNotify(orgId int64) error {
	mt, err := h.Store.Mousetrap.GetAll(orgId)
	if err != nil {
		log.Printf("getmousetraps: error: %v", err)
		return err
	}
	if len(mt) == 0 {
		log.Println("getmousetraps: no such mousetrap")
		return err
	}
	bs, err := json.Marshal(mt)
	if err != nil {
		log.Printf("getmousetraps: error: %v", err)
		return err
	}
	h.PubSub.Notify(orgId, string(bs))
	return nil
}

func (h Handler) GetLog(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.Logs.Logs))
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body models.Credentials true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /org/sign-in [post]
func (h Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	cred := &models.Credentials{}
	if err := json.NewDecoder(r.Body).Decode(cred); err != nil {
		log.Println("signin: request decode error: " + err.Error())
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := cred.CheckNotEmpty(); err != nil {
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

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body models.Credentials true "credentials"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /org/sign-up [post]
func (h Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	cred := &models.Credentials{}
	err := json.NewDecoder(r.Body).Decode(cred)
	if err != nil {
		log.Println("signup: request decode error: " + err.Error())
		WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := cred.CheckNotEmpty(); err != nil {
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
