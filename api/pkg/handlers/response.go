package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response struct{}

func (resp *Response) ErrorMessage(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	log.Printf("%s", err.Error())
	fmt.Fprintf(w, "%s", err.Error())
}

func (resp *Response) JsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
