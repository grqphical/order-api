package main

type Status int

const (
	OrderRecieved       Status = 0
	OrderProcessing     Status = 1
	OrderOutForDelivery Status = 2
	OrderShipped        Status = 3
)

type Item struct {
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Quantity int     `json:"quantity"`
}

type Order struct {
	ID          string `json:"id"`
	Active      bool   `json:"active"`
	Items       []Item `json:"items"`
	Address     string `json:"address"`
	Recipient   string `json:"recipient"`
	OrderStatus Status `json:"orderStatus"`
}
