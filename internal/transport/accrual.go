package transport

type Accrual interface {
	SendOrder() error
}
