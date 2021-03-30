package http

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	authHeader = "Authorization"
	orgIdHeader = "OrgId"
)
func (h Handler) AuthChecker(handlerFunc http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authHeader)

		if header == "" {
			WriteJSONError(w, "auth header is empty", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			WriteJSONError(w,  "invalid auth header", http.StatusUnauthorized)
			return
		}

		if headerParts[1] == "" {
			WriteJSONError(w,  "token is empty", http.StatusUnauthorized)
			return
		}

		id, err := h.tokenService.ParseToken(headerParts[1])
		if err != nil {
			WriteJSONError(w,  err.Error(), http.StatusUnauthorized)
			return
		}

		log.Printf("user: %v authorized", id)

		r.Header.Set(orgIdHeader, strconv.FormatInt(id, 10))
		handlerFunc(w, r)
	}
}