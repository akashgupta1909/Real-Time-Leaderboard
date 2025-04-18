package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/akashgupta1909/Real-Time-Leaderboard/internal/auth"
	"github.com/akashgupta1909/Real-Time-Leaderboard/repository"
	"github.com/akashgupta1909/Real-Time-Leaderboard/utils"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(repo *repository.UserRepository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			var userId string
			if authHeader == "" {
				utils.RespondWithError(w, http.StatusUnauthorized, "Authorization header is missing")
				return
			}

			tokenType := strings.Split(authHeader, " ")[0]
			if tokenType != "Bearer" {
				utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token type")
				return
			}

			tokenString := strings.Split(authHeader, " ")[1]
			if tokenString == "" {
				utils.RespondWithError(w, http.StatusUnauthorized, "Token is missing")
				return
			}

			token, err := auth.ValidateJWT(tokenString)
			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "Malformed jwt token")
				return
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userId, ok = claims["userId"].(string)
				if !ok {
					utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token claims")
					return
				}
			} else {
				utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			user, err := repo.FindUserByID(userId)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					utils.RespondWithError(w, http.StatusUnauthorized, "User not found")
					return
				}
				utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
			if user.ID.Hex() == "" {
				utils.RespondWithError(w, http.StatusUnauthorized, "User not found")
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, UserContextKey, user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
