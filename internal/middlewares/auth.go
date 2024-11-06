package middlewares

import (
	"context"
	"net/http"

	"git.sr.ht/~rehandaphedar/tilawah-hub/internal/db"
	"github.com/go-chi/render"
	"github.com/spf13/viper"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_token")
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, render.M{
				"message": "Missing session token",
				"error":   err.Error(),
			})
			return
		}
		sessionToken := sessionCookie.Value

		session, err := db.Queries.AuthSelectSession(context.Background(), sessionToken)
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, render.M{
				"message": "Invalid session token",
				"error":   err.Error(),
			})
			return
		}

		if !viper.GetBool("disable_csrf_checks") {
			csrfToken := r.Header.Get("X-CSRF-TOKEN")
			if session.CsrfToken != csrfToken {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, render.M{
					"messages": "Invalid csrf token",
					"error":    "Discrepancy between session token and csrf token",
				})
				return
			}
		}

		ctx := context.WithValue(r.Context(), "username", session.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
