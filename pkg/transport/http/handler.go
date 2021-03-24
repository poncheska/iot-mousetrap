package http

import (
	"github.com/poncheska/iot-mousetrap/pkg/store"
	"net/http"
)

type Handler struct{
	Store store.Store
}

func (h Handler) GetMousetraps(w http.ResponseWriter, r *http.Request){

}

func (h Handler) Trigger(w http.ResponseWriter, r *http.Request){

}