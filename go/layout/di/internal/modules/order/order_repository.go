package order

func newRepositoryOrder() repositoryOrder {
	return &orderMongo{}
}

type repositoryOrder interface {
	GetOrderID(account string) string
}

type orderMongo struct {
	_ struct{}
}

func (om *orderMongo) GetOrderID(account string) string {
	return "orderID-001"
}
