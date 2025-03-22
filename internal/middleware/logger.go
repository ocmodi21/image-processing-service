package middleware

import (
	"log"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request URI and method
		log.Printf("Request URI: %s, Method: %s", r.RequestURI, r.Method)
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
