package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/darthvadr/rss-aggregator/internal/auth"
	"github.com/darthvadr/rss-aggregator/internal/database"
)

type contextKey string

const userContextKey contextKey = "user"

// handleUserError handles errors encountered during user-related operations.
// It checks if the error is due to a missing user (sql.ErrNoRows) and responds
// with a 404 Not Found status and a relevant message. For all other errors,
// it logs the error and responds with a 500 Internal Server Error.
func handleUserError(w http.ResponseWriter, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("User not found")
		responseWithError(w, http.StatusNotFound, "user not found")
		return
	}
	log.Printf("error fetching user: %v", err)
	responseWithError(w, http.StatusInternalServerError, "internal server error")
}

func (apiConfig *apiConfig) middlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			responseWithError(w, http.StatusUnauthorized, fmt.Sprintf("invalid API key: %v", err))
			return
		}

		user, err := apiConfig.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			handleUserError(w, err)
			return
		}

		// Store user in context
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getUserFromContext retrieves the authenticated user from the request context
// and returns an error if the user is not found.
// This function is useful for handlers that need to access the authenticated user's information.
// It is request scoped and should be called within the context of an HTTP request.
func getUserFromContext(r *http.Request) (User, error) {
	user, ok := r.Context().Value(userContextKey).(database.User)
	if !ok {
		return User{}, errors.New("user not found in context")
	}

	mappedUser := databaseToUser(user)
	return mappedUser, nil
}