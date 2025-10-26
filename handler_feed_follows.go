package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/darthvadr/rss-aggregator/internal/database"
	chi "github.com/go-chi/chi/v5"
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

func (apiConfig *apiConfig) handlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request) {

	feedFollowIdString := chi.URLParam(r, "feedFollowId")

	feedFollowIdUuid, err := uuid.Parse(feedFollowIdString)
	if err != nil {
		log.Println("Error parsing feed follow ID: ", fmt.Errorf("error parsing feed follow ID: %w", err))
		responseWithError(w, http.StatusBadRequest, "invalid feed follow ID")
		return
	}

	user, err := getUserFromContext(r)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// TODO: Check if the feed follow exists before attempting to delete it

	err = apiConfig.DB.DeleteFeedFollows(r.Context(), database.DeleteFeedFollowsParams{
		FeedID: uuid.NullUUID{UUID: feedFollowIdUuid, Valid: true},
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
	})
	if err != nil {
		log.Println("Error deleting feed follows: ", fmt.Errorf("error deleting feed follows: %w", err))
		responseWithError(w, http.StatusInternalServerError, "error deleting feed follows")
		return
	}
	
	responseWithJSON(w, http.StatusOK, struct{}{})
}