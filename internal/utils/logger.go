package utils

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Log is an instance of logger.
var Log = log.New()

func init() {
	Log.SetLevel(log.InfoLevel)
	Log.SetFormatter(&log.JSONFormatter{})
}

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// WithLogging is an HTTP handler adding incoming request logging.
func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		Log.WithFields(log.Fields{
			"uri":         r.RequestURI,
			"method":      r.Method,
			"status":      responseData.status,
			"duration_μs": duration.Microseconds(),
			"size":        responseData.size,
		}).Info("Served HTTP request")
	})
}
