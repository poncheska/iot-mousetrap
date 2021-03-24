package http

import "net/http"

func (h Handler) AuthChecker(handlerFunc http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO auth
		handlerFunc(w,r)
	}
}
