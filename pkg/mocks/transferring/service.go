package transferring

type MockService struct {
	Err error
}

func (m *MockService) BalanceBetweenAccounts(originBalance float64, destinationBalance float64, _ float64) (_ float64, _ float64, _ error) {
	return originBalance, destinationBalance, m.Err
}
