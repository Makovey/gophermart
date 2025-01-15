package model

import "github.com/Makovey/gophermart/internal/types"

type OrderDetails struct {
	Order string       `json:"order"`
	Goods []OrderGoods `json:"goods"`
}

type OrderGoods struct {
	Description string             `json:"description"`
	Price       types.FloatDecimal `json:"price"`
}
