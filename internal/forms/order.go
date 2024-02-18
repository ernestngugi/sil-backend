package forms

type CreateOrderForm struct {
	Amount float64 `json:"amount"`
	Item string `json:"item"`
}
