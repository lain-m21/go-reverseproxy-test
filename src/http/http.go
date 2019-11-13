package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type (
	sayRequest struct {
		Message string
	}
	sayResponse struct {
		Message string
		Flags   []string
	}
)

func NewBackendServer(endpoint string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle(endpoint, http.HandlerFunc(sayHandler))
	server := &http.Server{
		Handler: mux,
	}
	return server
}

func sayHandler(w http.ResponseWriter, r *http.Request) {
	var request sayRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
	}
	flags := r.Header.Get("flags")
	flagSet := strings.Split(flags, ",")

	response := sayResponse{
		Message: fmt.Sprintf("Test service received a message: %s", request.Message),
		Flags:   flagSet,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
