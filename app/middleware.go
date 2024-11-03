package app

import (
	"github.com/szucik/trade-helper/apperrors"
	_ "github.com/szucik/trade-helper/portfolio"
	"net/http"
)

func MiddlewareIsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/signup" && r.RequestURI != "/signin" {
			_, err := r.Cookie("X-Auth")
			if err != nil {
				if err == http.ErrNoCookie {
					apperrors.Error("Invalid user", "BadRequest", 400).JSONError(rw)
					rw.WriteHeader(http.StatusUnauthorized)
					return
				}
				rw.WriteHeader(http.StatusBadRequest)
				apperrors.Error("Bad data", "BadRequest", 400).JSONError(rw)
				return
			}
		}

		next.ServeHTTP(rw, r)
	})
}
