package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestIndex(t *testing.T) {
	response := "OK"
	r := SetUpRouter()
	r.GET("/", index)

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, err := io.ReadAll(w.Body)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, response, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}
