package model

import "github.com/shopspring/decimal"

type OrderDetails struct {
	Order string     `json:"order"`
	Goods OrderGoods `json:"goods"`
}

type OrderGoods struct {
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
}
