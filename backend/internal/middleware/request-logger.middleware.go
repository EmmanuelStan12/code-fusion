package middleware

import (
	"github.com/EmmanuelStan12/code-fusion/internal/common/utils"
	"net/http"
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

func RequestLoggerMiddleware(log *utils.Logger) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			writer := &ResponseWriterWrapper{
				ResponseWriter: w,
				StatusCode:     0,
			}
			defer func() {
				WriteLog(log, writer.StatusCode, time.Since(startTime), r)
			}()
			handler.ServeHTTP(writer, r)
		}
		return http.HandlerFunc(fn)
	}
}

func WriteLog(log *utils.Logger, status int, elapsed time.Duration, r *http.Request) {
	log.Info(" %s %s %d %d milliseconds", r.Method, r.URL.Path, status, elapsed)

	if status >= 400 {
		log.Error("%s %s returned %d in %v", r.Method, r.URL.Path, status, elapsed)
	}
}
