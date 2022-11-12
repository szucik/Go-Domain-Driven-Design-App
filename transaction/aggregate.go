package transaction

type Aggregate struct {
	tranasaction Transaction
}

func (a Aggregate) Transaction() Transaction {
	return a.tranasaction
}
