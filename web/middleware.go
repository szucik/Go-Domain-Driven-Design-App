package web

import (
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/szucik/trade-helper/apperrors"
)

func MiddlewareIsAuth(store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if r.RequestURI == "/signup" || r.RequestURI == "/signin" {
				next.ServeHTTP(rw, r)
				return
			}

			session, err := store.Get(r, "X-Auth")
			if err != nil || session.Values["username"] == nil {
				apperrors.Error("unauthorized", "Unauthorized", http.StatusUnauthorized).JSONError(rw)
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}
