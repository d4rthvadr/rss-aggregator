package main

import (
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	responseWithJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}
