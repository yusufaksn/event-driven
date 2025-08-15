package domain

type OrderItem struct {
	ProductID   string `json:"productID"`
	OrderID     string `json:"orderID"`
	Quantity    int    `json:"quantity"`
	Total       int    `json:"total"`
	Description string `json:"description"`
	EventID     string `json:"eventID"`
}
