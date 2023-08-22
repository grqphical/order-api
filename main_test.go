package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var router *gin.Engine = gin.Default()

func setupSuite(tb testing.TB) func(tb testing.TB) {

	err := os.WriteFile("orders.json", []byte("[]"), 0644)

	if err != nil {
		panic(err)
	}

	// Return a function to teardown the test
	return func(tb testing.TB) {
		e := os.Remove("orders.json")
		if e != nil {
			log.Fatal(e)
		}
	}
}

func TestIndex(t *testing.T) {

	router.GET("/", index)

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetOrder(t *testing.T) {
	orders = append(orders,
		Order{ID: "1",
			Active: true, Address: "123 Example Street",
			Items:       []Item{{Name: "Laptop", Quantity: 1}},
			Recipient:   "John Doe",
			OrderStatus: OrderRecieved})

	router.GET("/get-order", getOrder)

	req, err := http.NewRequest("GET", "/get-order?id=1", nil)

	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	responseData, err := io.ReadAll(w.Body)

	if err != nil {
		panic(err)
	}

	var order Order

	json.Unmarshal(responseData, &order)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, order, orders[0])
}

func TestAddOrder(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	orderToAdd := Order{ID: "2",
		Active: true, Address: "125 Example Street",
		Items:       []Item{{Name: "Jeans", Quantity: 2}},
		Recipient:   "Jean Doe",
		OrderStatus: OrderProcessing}

	router.POST("/add-order", addOrder)

	data, err := json.Marshal(orderToAdd)

	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", "/add-order", bytes.NewReader(data))

	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	responseData, err := io.ReadAll(w.Body)

	if err != nil {
		panic(err)
	}

	var order Order

	json.Unmarshal(responseData, &order)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, order, orders[1])

}

func TestUpdateOrder(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	router.PATCH("/update-order-status", updateOrderStatus)

	orders = append(orders,
		Order{ID: "1",
			Active: true, Address: "123 Example Street",
			Items:       []Item{{Name: "Laptop", Quantity: 1}},
			Recipient:   "John Doe",
			OrderStatus: OrderRecieved})

	form_data := url.Values{
		"status": {"3"},
	}

	req, err := http.NewRequest("PATCH", "/update-order-status?id=1", strings.NewReader(form_data.Encode()))

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	responseData, err := io.ReadAll(w.Body)

	if err != nil {
		panic(err)
	}

	var order Order

	json.Unmarshal(responseData, &order)

	assert.Equal(t, http.StatusAccepted, w.Code)
	assert.Equal(t, order.OrderStatus, orders[0].OrderStatus)
}

func TestCompleteOrder(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	router.PATCH("/complete-order", completeOrder)

	req, err := http.NewRequest("PATCH", "/complete-order?id=1", nil)

	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	responseData, err := io.ReadAll(w.Body)

	if err != nil {
		panic(err)
	}

	var order Order

	json.Unmarshal(responseData, &order)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, order.Active, false)
}

func TestEditOrder(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	router.PATCH("/edit-order", editOrder)

	// orders = append(orders,
	// 	Order{ID: "1",
	// 		Active: true, Address: "123 Example Street",
	// 		Items:       []Item{{Name: "Laptop", Quantity: 1}},
	// 		Recipient:   "John Doe",
	// 		OrderStatus: OrderRecieved})

	form_data := url.Values{
		"address":   {"240 Park Street"},
		"recipient": {"Jane Doe"},
	}

	req, err := http.NewRequest("PATCH", "/edit-order?id=1", strings.NewReader(form_data.Encode()))

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	responseData, err := io.ReadAll(w.Body)

	if err != nil {
		panic(err)
	}

	var order Order

	json.Unmarshal(responseData, &order)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, order.Address, form_data["address"][0])
	assert.Equal(t, order.Recipient, form_data["recipient"][0])
}

func TestGetOrderError(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	req, err := http.NewRequest("GET", "/get-order?id=5", nil)

	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddOrderErrors(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	data := []byte("{870sasdfals: iudsaiofh} /")

	req, err := http.NewRequest("POST", "/add-order", bytes.NewReader(data))

	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	body, err := io.ReadAll(w.Body)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Failed to parse JSON", string(body[:]))
}

func TestUpdateOrderStatusErrors(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	orders = []Order{{ID: "1",
		Active: false, Address: "123 Example Street",
		Items:       []Item{{Name: "Laptop", Quantity: 1}},
		Recipient:   "John Doe",
		OrderStatus: OrderRecieved}}

	form_data := url.Values{
		"status": {"3"},
	}

	req, err := http.NewRequest("PATCH", "/update-order-status?id=1", strings.NewReader(form_data.Encode()))

	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusLocked, w.Code)

	// Test what happens when server recieves an id it doesnt have

	orders = []Order{{ID: "1",
		Active: true, Address: "123 Example Street",
		Items:       []Item{{Name: "Laptop", Quantity: 1}},
		Recipient:   "John Doe",
		OrderStatus: OrderRecieved}}

	req, err = http.NewRequest("PATCH", "/update-order-status?id=foo", strings.NewReader(form_data.Encode()))

	if err != nil {
		panic(err)
	}

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

}
