package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetApiKey retrieves the API key from the provided HTTP headers.
// It looks for the "X-API-Key" header and returns its value.
// If the header is missing or empty, it returns an error indicating that the API key is missing.
func GetApiKey(headers http.Header) (string, error) {
	apiKey := headers.Get("X-API-Key")
	if apiKey == "" {
		return "", errors.New("missing API key")
	}

	vals := strings.Split(apiKey, " ")
	if len(vals) < 2 {
		return "", errors.New("malformed API key")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("invalid API key format")
	}
	return vals[1], nil
}