package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/darthvadr/rss-aggregator/internal/auth"
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

func (apiConfig *apiConfig) handlerGetUserByApiKey(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, fmt.Sprintf("invalid API key: %v", err))
		return
	}

	user, err := apiConfig.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {

		log.Printf("error fetching user: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("User not found for API key:", apiKey)
			responseWithError(w, http.StatusNotFound, "user not found")
			return
		}
		responseWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responseWithJSON(w, http.StatusOK, databaseToUser(user))
}