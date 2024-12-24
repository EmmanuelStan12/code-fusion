package middleware

import (
	"github.com/EmmanuelStan12/code-fusion/configs"
	"log"
	"net/http"
	"strings"
	"time"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	StatusCode   int
	BytesWritten int
}

func (rw *ResponseWriterWrapper) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseWriterWrapper) Write(data []byte) (int, error) {
	bytes, err := rw.ResponseWriter.Write(data)
	rw.BytesWritten += bytes
	return bytes, err
}

func RequestLoggerMiddleware(config configs.LoggerConfig) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			writer := &ResponseWriterWrapper{
				ResponseWriter: w,
				StatusCode:     0,
			}
			handler.ServeHTTP(writer, r)
			defer func() {
				WriteLog(config, time.Now(), writer.StatusCode, time.Since(startTime), r)
			}()
		}
		return http.HandlerFunc(fn)
	}
}

func WriteLog(config configs.LoggerConfig, currentTime time.Time, status int, elapsed time.Duration, r *http.Request) {
	if config.Level == configs.LogLevelInfo || config.Level == configs.LogLevelDebug || config.Level == configs.LogLevelError {
		log.Printf("[%s] %s %s %d %d", currentTime.Format(time.RFC3339), r.Method, r.URL.Path, status, elapsed)
	}

	if config.Level == configs.LogLevelDebug {
		for name, values := range r.Header {
			log.Printf("Header: %s=%s", name, strings.Join(values, ","))
		}
	}

	if config.Level == configs.LogLevelError && status >= 400 {
		log.Printf("[%s] ERROR: %s %s returned %d in %v", currentTime.Format(time.RFC3339), r.Method, r.URL.Path, status, elapsed)
	}
}
