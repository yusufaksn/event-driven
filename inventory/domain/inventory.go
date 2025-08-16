package domain

type InventoryItem struct {
	ProductID string `json:"productID"`
	OrderID   string `json:"orderID"`
	Quantity  int    `json:"quantity"`
	Message   string `json:"message"`
	EventID   string `json:"eventID"`
}
