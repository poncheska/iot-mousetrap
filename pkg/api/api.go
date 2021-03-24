package api

import (
	"github.com/gorilla/mux"
	"github.com/poncheska/iot-mousetrap/pkg/store/fake"
	hp "github.com/poncheska/iot-mousetrap/pkg/transport/http"
	"net/http"
)

func Start(){
	h := hp.Handler{
		Store: fake.NewFakeStore(),
	}
	r := mux.NewRouter()
	r.HandleFunc("/mousetraps", h.GetMousetraps).Methods(http.MethodGet)
	r.HandleFunc("/trigger/{org}/{name}", h.GetMousetraps).Methods(http.MethodPut)

}