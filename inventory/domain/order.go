package domain

type OrderItem struct {
	ProductID   string  `json:"productID"`
	OrderID     string  `json:"orderID"`
	Quantity    int     `json:"quantity"`
	Description string  `json:"description"`
	EventID     string  `json:"eventID"`
	Price       float32 `json:"price"`
}
