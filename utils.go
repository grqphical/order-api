package main

import (
	"encoding/json"
	"os"
)

// Removes an item from a slice based on an index
func remove(slice []Order, s int) []Order {
	return append(slice[:s], slice[s+1:]...)
}

// Saves the curent JSON data to the "database" which is just a JSON file
func saveDatabase(orders []Order) {
	bytes, err := json.Marshal(orders)

	if err != nil {
		panic(err)
	}

	os.WriteFile("orders.json", bytes, 0644)
}
