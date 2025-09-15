package middleware

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

type LoggingMiddleware struct {
	logger *slog.Logger
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func NewLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func NewLoggingMiddleware(logger *slog.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

func (m *LoggingMiddleware) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := newResponseWriter(w)

		next.ServeHTTP(rw, r)
		duration := time.Since(start)

		switch {
		case rw.statusCode >= 500:
			m.logger.Error("request", "method", r.Method, "path", r.URL.Path, "status", rw.statusCode, "duration", duration)
		case rw.statusCode >= 400:
			m.logger.Warn("request", "method", r.Method, "path", r.URL.Path, "status", rw.statusCode, "duration", duration)
		default:
			m.logger.Info("request", "method", r.Method, "path", r.URL.Path, "status", rw.statusCode, "duration", duration)
		}

	})
}
