package portfolio

import (
	"context"
	"fmt"
	"net/http"
)

type Database interface {
	AddPortfolio(ctx context.Context) (Aggregate, error)
}

type Portfolios struct {
	Database     Database
	NewAggregate func(Portfolio) (Aggregate, error)
}

type PortfolioKey struct {
}

func (p Portfolios) AddPortfolio(rw http.ResponseWriter, r *http.Request) {
	portfolio := &Portfolio{}
	ctx := context.WithValue(r.Context(), PortfolioKey{}, *portfolio)
	rw.Header().Set("Content-Type", "application/json")

	//err := json.NewDecoder(r.Body).Decode(&portfolio)
	//if err != nil {
	//	http.Error(rw, err.Error(), http.StatusBadRequest)
	//	return
	//}

	//portfolio = Portfolio{
	//	UserId:  portfolio.UserId,
	//	Name:    portfolio.Name,
	//	Created: time.Time{},
	//}
	//aggregate, err := p.NewAggregate(portfolio)

	agg, _ := p.Database.AddPortfolio(ctx)
	fmt.Print(agg)
	//_, err := p.Database.Add(portfolio)
	//if err != nil {
	//	message := fmt.Sprintf("Error message: %v", err)
	//	p.l.Print(message)
	//	p.db.ToJSON(&GenericResponse{Status: false, Message: err.Error()}, rw)
	//	return
	//}
	//p.l.Print("PortfolioKey created successfully")
	//
	//p.db.ToJSON(&GenericResponse{Status: true, Message: "Portfolio created successfully"}, rw)
}
