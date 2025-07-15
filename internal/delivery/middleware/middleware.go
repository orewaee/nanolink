package middleware

import "net/http"

type Middleware interface {
	Use(next http.Handler) http.Handler
}
