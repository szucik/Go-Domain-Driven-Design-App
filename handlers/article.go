package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Article struct {
	l *log.Logger
}

func (h *Article) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Logger in articles is working")
	d, error := ioutil.ReadAll(r.Body)
	if error != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(rw, "Hello %s", d)
}

func NewArticle(l *log.Logger) *Article {
	return &Article{l}
}
