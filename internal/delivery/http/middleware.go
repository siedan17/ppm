package http

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// MethodOverride supports PUT/PATCH/DELETE via _method form field or X-HTTP-Method-Override header
func MethodOverride(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			override := r.FormValue("_method")
			if override == "" {
				override = r.Header.Get("X-HTTP-Method-Override")
			}
			switch override {
			case "PUT", "PATCH", "DELETE":
				r.Method = override
			}
		}
		next.ServeHTTP(w, r)
	})
}
