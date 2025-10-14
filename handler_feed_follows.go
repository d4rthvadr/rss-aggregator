package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/darthvadr/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiConfig *apiConfig) handlerCreateFeedFollows(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"` 
	}

	user, err := getUserFromContext(r)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Println("Error decoding request body: ", fmt.Errorf("error decoding request body: %w", err))
		responseWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	createdFeedFollows, err := apiConfig.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:     uuid.New(),
		FeedID: uuid.NullUUID{UUID: params.FeedId, Valid: true},
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
	})

	if err != nil {

		log.Println("Error creating feed: ",  fmt.Errorf("error creating feed: %w", err))
		responseWithError(w, http.StatusInternalServerError, "error creating feed")
		return
	}

	responseWithJSON(w, http.StatusOK, databaseToFeedFollows(createdFeedFollows))
}


func (apiConfig *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request) {

	user, err := getUserFromContext(r)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}


	feedFollows, err := apiConfig.DB.GetFeedFollows(r.Context(), uuid.NullUUID{UUID: user.ID, Valid: true})
	if err != nil {
		log.Println("Error listing feed follows: ", fmt.Errorf("error listing feed follows: %w", err))
		responseWithError(w, http.StatusInternalServerError, "error listing feed follows")
		return
	}
	mappedFeedFollows := databaseFeedFollowsToFeedFollows(feedFollows)

	responseWithJSON(w, http.StatusOK, mappedFeedFollows)
}
