package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"fmt"

	"strconv"

	"os"
)

var orders []Order

// Used to get status of API
func index(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

// Adds an order to the order processing system
func addOrder(c *gin.Context) {
	var newOrder Order

	if err := c.BindJSON(&newOrder); err != nil {
		c.String(400, "Failed to parse JSON")
		return
	}

	if err := ValidateStruct(newOrder); err != nil {
		c.String(400, "Bad Request")
		return
	}

	orders = append(orders, newOrder)

	saveDatabase(orders)

	c.JSON(http.StatusCreated, newOrder)
}

// Returns an order based on a given ID
func getOrder(c *gin.Context) {
	id := c.Query("id")

	for i := range orders {
		if orders[i].ID == id {
			c.JSON(http.StatusOK, orders[i])
			return
		}
	}

	errMsg := fmt.Sprintf("Order with id '%s' not found", id)

	c.String(http.StatusNotFound, errMsg)
}

// Updates an order's status
func updateOrderStatus(c *gin.Context) {
	id := c.Query("id")
	status := c.PostForm("status")

	for i := range orders {
		if orders[i].ID == id {
			if !orders[i].Active {
				c.String(http.StatusLocked, "Order is no longer active")
				return
			}

			status_int, err := strconv.Atoi(status)

			if err != nil {
				c.JSON(http.StatusBadRequest, "Error Parsing ID")
				return
			}

			orders[i].OrderStatus = Status(status_int)

			saveDatabase(orders)

			c.JSON(http.StatusAccepted, orders[i])
			return
		}
	}

	errMsg := fmt.Sprintf("Order with id '%s' not found", id)

	c.String(http.StatusNotFound, errMsg)
}

// Removes an order from the system
func removeOrder(c *gin.Context) {
	id := c.Query("id")

	for i := range orders {
		if orders[i].ID == id {
			orders = remove(orders, i)

			saveDatabase(orders)

			c.JSON(http.StatusOK, orders[i])
			return
		}
	}

	errMsg := fmt.Sprintf("Order with id '%s' not found", id)

	c.String(http.StatusNotFound, errMsg)
}

// Changes an order's status to completed and deactivates it
// all past orders are archived
func completeOrder(c *gin.Context) {
	id := c.Query("id")

	for i := range orders {
		if orders[i].ID == id {
			orders[i].Active = false
			orders[i].OrderStatus = OrderShipped

			saveDatabase(orders)

			c.JSON(http.StatusOK, orders[i])
			return
		}
	}

	errMsg := fmt.Sprintf("Order with id '%s' not found", id)

	c.String(http.StatusNotFound, errMsg)
}

func editOrder(c *gin.Context) {
	id := c.Query("id")
	address := c.PostForm("address")
	recipient := c.PostForm("recipient")

	for i := range orders {
		if orders[i].ID == id {
			if address != "" {
				orders[i].Address = address
			}

			if recipient != "" {
				orders[i].Recipient = recipient
			}

			saveDatabase(orders)

			c.JSON(http.StatusOK, orders[i])
			return
		}
	}

	errMsg := fmt.Sprintf("Order with id '%s' not found", id)

	c.String(http.StatusNotFound, errMsg)

}

func main() {
	// Read our database file
	data, err := os.ReadFile("orders.json")

	if err != nil {
		panic(err)
	}

	// Decode the JSON into a list of Order structs
	err = json.Unmarshal(data, &orders)

	if err != nil {
		panic(err)
	}

	// Setup our API webserver
	router := gin.Default()

	router.GET("/", index)
	router.POST("/add-order", addOrder)
	router.GET("/get-order", getOrder)
	router.PATCH("/update-order-status", updateOrderStatus)
	router.DELETE("/remove-order", removeOrder)
	router.PATCH("/complete-order", completeOrder)
	router.PATCH("/edit-order", editOrder)
	router.Run("localhost:6969")
}
