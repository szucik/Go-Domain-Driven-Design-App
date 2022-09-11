package handlers

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/szucik/go-simple-rest-api/internal/dao"
)

func (u *Users) MiddlewareLoginValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		credentials := &dao.AuthCredentials{}
		err := u.db.FromJSON(credentials, r.Body)
		if err != nil {
			u.l.Println("[ERROR] deserializing user", err)
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		err = u.db.Validate(credentials)
		if err != nil {
			u.l.Println("[ERROR] validate user", err)
			http.Error(rw, "Unable validate user", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey{}, *credentials)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}

func (u *Users) MiddlewareIsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("Authorization")
		if err != nil {
			if err == http.ErrNoCookie {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		tknStr := c.Value
		claims := &customClaims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return accessKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(rw, r)
	})
}

func (u *Users) MiddlewareUserValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		usr := &dao.User{}
		err := u.db.FromJSON(usr, r.Body)
		if err != nil {
			u.l.Println("[ERROR] deserializing user", err)
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		err = u.db.Validate(usr)
		if err != nil {
			u.l.Println("[ERROR] validate user", err)
			http.Error(rw, "Unable validate user", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey{}, *usr)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}

func (p *Portfolios) MiddlewarePortfolioValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		portfolio := &dao.Portfolio{}
		err := p.db.FromJSON(portfolio, r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing portfolio", err)
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		err = p.db.Validate(portfolio)
		if err != nil {
			p.l.Println("[ERROR] validate portfolio", err)
			http.Error(rw, "Unable validate portfolio", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), PortfolioKey{}, *portfolio)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}

func (t *Transactions) MiddlewareTransactionValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		transaction := &dao.Transaction{}
		err := t.db.FromJSON(transaction, r.Body)
		if err != nil {
			t.l.Println("[ERROR] deserializing Transaction", err)
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		err = t.db.Validate(transaction)
		if err != nil {
			t.l.Println("[ERROR] validate Transaction", err)
			http.Error(rw, "Unable validate Transaction", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), TransactionKey{}, *transaction)

		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
