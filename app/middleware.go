package app

//
//import (
//	"context"
//	_ "github.com/szucik/trade-helper/portfolio"
//	"net/http"
//)
//
//func MiddlewareLoginValid(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
//		credentials := &dao2.AuthCredentials{}
//		err := u.db.FromJSON(credentials, r.Body)
//		if err != nil {
//			u.l.Println("[ERROR] deserializing user", err)
//			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
//			return
//		}
//
//		err = u.db.Validate(credentials)
//		if err != nil {
//			u.l.Println("[ERROR] validate user", err)
//			http.Error(rw, "Unable validate user", http.StatusBadRequest)
//			return
//		}
//
//		ctx := context.WithValue(r.Context(), UserKey{}, *credentials)
//		r = r.WithContext(ctx)
//		next.ServeHTTP(rw, r)
//	})
//}
//
//func MiddlewareIsAuth(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
//		c, err := r.Cookie("Authorization")
//		if err != nil {
//			if err == http.ErrNoCookie {
//				rw.WriteHeader(http.StatusUnauthorized)
//				return
//			}
//
//			rw.WriteHeader(http.StatusBadRequest)
//			return
//		}
//
//		tknStr := c.Value
//		claims := &customClaims{}
//
//		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
//			return accessKey, nil
//		})
//
//		if err != nil {
//			if err == jwt.ErrSignatureInvalid {
//				rw.WriteHeader(http.StatusUnauthorized)
//				return
//			}
//			rw.WriteHeader(http.StatusBadRequest)
//			return
//		}
//
//		if !tkn.Valid {
//			rw.WriteHeader(http.StatusUnauthorized)
//			return
//		}
//		next.ServeHTTP(rw, r)
//	})
//}
