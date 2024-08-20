package middleware

import (
    "net/http"
)

func AuthMiddleware(token string, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader != "Bearer "+token {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}