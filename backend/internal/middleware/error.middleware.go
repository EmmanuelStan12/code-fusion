package middleware

import (
	"encoding/json"
	"github.com/EmmanuelStan12/code-fusion/configs"
	"github.com/EmmanuelStan12/code-fusion/internal/common/errors"
	"github.com/EmmanuelStan12/code-fusion/internal/common/utils"
	"log"
	"net/http"
)

func ErrorMiddleware(localeConfig *configs.LocaleConfig) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err != nil {
					var response utils.ApiResponse[any]
					switch e := err.(type) {
					case *errors.CustomError:
						response = utils.BuildResponse[any](nil, false, e.Code, e.Type, localeConfig)
					case error:
						response = utils.BuildResponse[any](nil, false, http.StatusInternalServerError, e.Error(), localeConfig)
					default:
						response = utils.BuildResponse[any](nil, false, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", localeConfig)
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(response.StatusCode)
					if err := json.NewEncoder(w).Encode(response); err != nil {
						log.Printf("Error encoding response: %v", err)
					}
				}
			}()
			handler.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
