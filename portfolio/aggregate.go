package portfolio

type Aggregate struct {
	portfolio Portfolio
}

func (a Aggregate) Portfolio() Portfolio {
	return a.portfolio
}
