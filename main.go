package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"fmt"

	"os"

    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"

    // docs "example/order-api/docs"
)

var orders []Order

// swagger:model
type IndexResponse struct {
    DocsUrl string `json:"documentationUrl"`
}

// Index godoc
//
// @Summary Base Route
// @Schemes http https
// @Produce plain
// @Success 200 {object} IndexResponse
// @Router / [get]
func index(c *gin.Context) {
	c.JSON(http.StatusOK, IndexResponse { DocsUrl: "/swagger/index.html"})
}

// AddOrder godoc
//
// @Summary Adds an order to the system
// @Schemes http https
// @Accept json
// @Produce json
// @Param order body Order true "Order"
// @Success 201 {object} Order
// @Failure 500 {string} string "Failed to parse JSON"
// @Router /add-order [post]
func addOrder(c *gin.Context) {
	var newOrder Order

	if err := c.BindJSON(&newOrder); err != nil {
		c.String(500, "Failed to parse JSON")
		return
	}

	orders = append(orders, newOrder)

	saveDatabase(orders)

	c.JSON(http.StatusCreated, newOrder)
}

// GetOrder godoc
//
// @Summary Adds an order to the system
// @Param   id  query    int true "Order ID"
// @Schemes http https
// @Produce json
// @Success 200 {object} Order
// @Failure 404 {string} string "Order with ID 'X' not found"
// @Router /get-order [get]
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

// UpdateOrderStatus godoc
//
// @Summary Updates an order's status
// @Param   id      query    int    true    "Order ID"
// @Param   status  query    Status true    "Order Status"
// @Schemes http https
// @Produce json
// @Success 202 {object} Order
// @Failure 423 {string} string "Order is no longer active"
// @Failure 404 {string} string "Order with id 'X' not found"
// @Router /update-order-status [patch]
func updateOrderStatus(c *gin.Context) {
	id := c.Query("id")
	status := c.PostForm("status")

	for i := range orders {
		if orders[i].ID == id {
			if !orders[i].Active {
				c.String(http.StatusLocked, "Order is no longer active")
				return
			}

			orders[i].OrderStatus = Status(status)

			saveDatabase(orders)

			c.JSON(http.StatusAccepted, orders[i])
			return
		}
	}

	errMsg := fmt.Sprintf("Order with id '%s' not found", id)

	c.String(http.StatusNotFound, errMsg)
}

// RemoveOrder godoc
//
// @Summary Removes an order from the system
// @Param   id  query    int true "Order ID"
// @Schemes http https
// @Produce json
// @Success 200 {object} Order
// @Failure 404 {string} string "Order with ID 'X' not found"
// @Router /remove-order [delete]
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

// CompleteOrder godoc
//
// @Summary Deactivates an order and archives it 
// @Param   id  query    int true "Order ID"
// @Schemes http https
// @Produce json
// @Success 200 {object} Order
// @Failure 404 {string} string "Order with ID 'X' not found"
// @Router /complete-order [patch]
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

// EditOrder godoc
//
// @Summary Removes an order from the system
// @Param   id          query       int     true    "Order ID"
// @Param   address     formData    string  true    "Address"
// @Param   recipient   formData    string  true    "Recipient"
// @Schemes http https
// @Accept x-www-form-urlencoded
// @Produce json
// @Success 200 {object} Order
// @Failure 404 {string} string "Order with ID 'X' not found"
// @Router /edit-order [patch]
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

//  @title Order API
//  @version 1.0
//  @description A simple Order tracking API for an ecommerce site. View source code here: https://github.com/grqphical07/order-api
//  @license.name MIT
//  @license.url https://github.com/grqphical07/order-api/blob/main/LICENSE
//  @BasePath /
// @Schemes http https
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

    router.StaticFile("/docs/swagger.json", "docs/swagger.json")

    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/swagger.json")))

	router.GET("/", index)
	router.POST("/add-order", addOrder)
	router.GET("/get-order", getOrder)
	router.PATCH("/update-order-status", updateOrderStatus)
	router.DELETE("/remove-order", removeOrder)
	router.PATCH("/complete-order", completeOrder)
	router.PATCH("/edit-order", editOrder)
	router.Run("localhost:6969")
}
