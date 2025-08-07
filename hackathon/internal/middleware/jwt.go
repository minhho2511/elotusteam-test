package middleware

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/minhho2511/elotusteam-test/utils"
	"net/http"
	"strings"
)

const UserContextKey = "user"

func extractToken(r *http.Request) string {
	// Try Authorization header first (standard)
	token := r.Header.Get("Authorization")
	if token == "" {
		// Fallback to lowercase (some clients might use this)
		token = r.Header.Get("authorization")
	}

	if token == "" {
		return ""
	}

	// Remove "Bearer " prefix if present
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimPrefix(token, "bearer ")

	return strings.TrimSpace(token)
}

func writeErrorResponse(w http.ResponseWriter, message utils.Message) {
	resp := utils.SetDefaultResponse(context.Background(), message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(message.Code)
	json.NewEncoder(w).Encode(resp)
}

func Authenticate(jwtSvc *utils.JWTService) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			token := extractToken(r)

			// Check if token exists
			if token == "" {
				writeErrorResponse(w, utils.Message{Code: http.StatusUnauthorized, Message: "MISSING_TOKEN"})
				return
			}

			// Validate JWT token (checks blacklist, expiration, and signature)
			claims, err := jwtSvc.ValidateJWT(token)
			if err != nil {
				writeErrorResponse(w, utils.Message{Code: http.StatusUnauthorized, Message: "TOKEN_INVALID"})
				return
			}

			// Add user info to request context
			ctx := context.WithValue(r.Context(), UserContextKey, claims.Username)
			r = r.WithContext(ctx)

			// Continue to next handler
			h.ServeHTTP(w, r)
		})
	}
}
