package handlers

import (
	"encoding/json"
	"github.com/szucik/go-simple-rest-api/data"
	"log"
	"net/http"
)

type Currency struct {
	l *log.Logger
}

func NewCurrency(l *log.Logger) *Currency {
	return &Currency{l}
}

func (c *Currency) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	lc := data.GetCurrency()
	d, err := json.Marshal(lc)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
	rw.Write(d)
}
