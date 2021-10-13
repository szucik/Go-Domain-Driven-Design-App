package data

type Currency struct {
	ID            string
	Name          string
	Price         float32
	PurchasePrice float32
	PurchaseDate  string
}

var currencyList = []*Currency{
	&Currency{
		ID:            "1",
		Name:          "BTC",
		Price:         51324,
		PurchasePrice: 36000,
	},
	&Currency{
		ID:            "1",
		Name:          "BTC",
		Price:         3400,
		PurchasePrice: 2431,
	},
}

func GetCurrency() []*Currency {
	return currencyList
}
