package model

type OrderDetails struct {
	Order string       `json:"order"`
	Goods []OrderGoods `json:"goods"`
}

type OrderGoods struct {
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
