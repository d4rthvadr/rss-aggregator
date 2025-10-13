package main

import (
	"encoding/json"
	"fmt"
	"log"
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
		log.Println("Error decoding request body: ", fmt.Errorf("error decoding request body: %w", err))
		responseWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	createdUser, err := apiConfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:   uuid.New(),
		Name: params.Name,
	})

	if err != nil {
		log.Println("Error creating user: ", fmt.Errorf("error creating user: %w", err))
		responseWithError(w, http.StatusInternalServerError, "error creating user")
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

	responseWithJSON(w, http.StatusOK, user)
}