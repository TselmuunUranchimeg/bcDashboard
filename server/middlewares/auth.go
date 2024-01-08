package middlewares

import (
	"context"
	"net/http"

	"bcDashboard/services"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err == nil {
			t, err := services.VerifyToken(token.Value)
			if err == nil {
				var name services.Username = "username"
				ctx := context.WithValue(r.Context(), name, t.Subject())
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}
