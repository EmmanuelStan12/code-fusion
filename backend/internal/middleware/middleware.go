package middleware

import "net/http"

type BaseMiddleware interface {
	Matcher(r *http.Request) bool
	Intercept(http.Handler) http.Handler
}
