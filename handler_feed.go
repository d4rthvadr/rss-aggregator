package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/darthvadr/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiConfig *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Title string `json:"title"` 
		Url   string `json:"url"`
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

	createdFeed, err := apiConfig.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:     uuid.New(),
		Title:  params.Title,
		Url:    params.Url,
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
	})

	if err != nil {

		log.Println("Error creating feed: ",  fmt.Errorf("error creating feed: %w", err))
		responseWithError(w, http.StatusInternalServerError, "error creating feed: ")
		return
	}

	responseWithJSON(w, http.StatusOK, databaseToFeed(createdFeed))
}

func (apiConfig *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	
	feeds, err := apiConfig.DB.GetFeeds(r.Context())
	if err != nil {
		log.Println("Error listing feeds: ", fmt.Errorf("error listing feeds: %w", err))
		responseWithError(w, http.StatusInternalServerError, "error listing feeds")
		return
	}


	mappedFeeds := databaseFeedsToFeeds(feeds)

	responseWithJSON(w, http.StatusOK, mappedFeeds)
}


