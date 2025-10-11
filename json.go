package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func responseWithJSON(w http.ResponseWriter, status int, data interface{}) {

	payload, err := json.Marshal(data)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(payload)
}


func responseWithError(w http.ResponseWriter, status int, message string) {

	if status > 499 {
		log.Println("5XX error:", message)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}
	responseWithJSON(w, status, errorResponse{Error: message})
}

