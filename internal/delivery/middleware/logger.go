package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/phuslu/log"
)

type LoggerMiddleware struct{}

func NewLoggerMiddleware() Middleware {
	return &LoggerMiddleware{}
}

func (middleware *LoggerMiddleware) Use(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		traceId := uuid.New().String()
		now := time.Now()

		log.Info().
			Str("path", request.RequestURI).
			Str("user_agent", request.Header.Get("User-Agent")).
			Str("trace_id", traceId).
			Msg("request processing...")

		next.ServeHTTP(writer, request)

		log.Info().
			Str("elapsed", time.Since(now).String()).
			Str("trace_id", traceId).
			Msg("request processed")
	})
}
