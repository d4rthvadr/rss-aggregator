package main

import (
	"encoding/json"
	"net/http"

	"github.com/darthvadr/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiConfig *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Name string `json:"name"`
	}
	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		responseWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	createdUser, err := apiConfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:   uuid.New(),
		Name: params.Name,
	})

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "error creating user: "+err.Error())
		return
	}

	responseWithJSON(w, http.StatusOK, databaseToUser(createdUser))
}

func (apiConfig *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromContext(r)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "could not find user")
		return
	}

	responseWithJSON(w, http.StatusOK, databaseToUser(user))
}