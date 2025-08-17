package domain

type PaymentItem struct {
	OrderID    string  `json:"orderID"`
	PaymentID  string  `json:"paymentID"`
	Message    string  `json:"message"`
	EventID    string  `json:"eventID"`
	TotalPrice float32 `json:"totalPrice"`
}
