package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	th_data "github.com/szucik/go-simple-rest-api/internal/dao"
)

type Portfolios struct {
	l  *log.Logger
	db *th_data.Dao
}

func NewPortfolios(l *log.Logger, db *th_data.Dao) *Portfolios {
	return &Portfolios{l, db}
}

type PortfolioKey struct{}

func (p *Portfolios) AddPortfolio(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rc := r.Context().Value(PortfolioKey{}).(th_data.Portfolio)

	portfolio := &th_data.Portfolio{
		UserId:  rc.UserId,
		Name:    rc.Name,
		Created: time.Time{},
	}

	_, err := p.db.Add(portfolio)
	if err != nil {
		message := fmt.Sprintf("Error message: %v", err)
		p.l.Print(message)
		p.db.ToJSON(&GenericResponse{Status: false, Message: err.Error()}, rw)
		return
	}
	p.l.Print("PortfolioKey created successfully")

	p.db.ToJSON(&GenericResponse{Status: true, Message: "Portfolio created successfully"}, rw)
}
