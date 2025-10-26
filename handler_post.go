package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/darthvadr/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiConfig *apiConfig) handlerGetPostForUser(w http.ResponseWriter, r *http.Request) {

	
	user, err := getUserFromContext(r)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	posts, err := apiConfig.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
		Limit:  10,
	})

	if err != nil {

		log.Println("Error getting posts: ",  fmt.Errorf("error getting posts: %w", err))
		responseWithError(w, http.StatusInternalServerError, "error getting posts: ")
		return
	}

	mappedPosts := databasePostsToPosts(posts)
	responseWithJSON(w, http.StatusOK, mappedPosts)
}



