package main

// swagger:enum Status
type Status string

const (
	OrderRecieved       Status = "OrderRecieved"
	OrderProcessing     Status = "OrderProcessing"
	OrderOutForDelivery Status = "OrderOutForDelivery"
	OrderShipped        Status = "OrderShipped"
)

type Item struct {
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Quantity int     `json:"quantity"`
}

// swagger:model
type Order struct {
	ID          string `json:"id"`
	Active      bool   `json:"active"`
	Items       []Item `json:"items"`
	Address     string `json:"address"`
	Recipient   string `json:"recipient"`
    OrderStatus Status `json:"orderStatus"`
}
