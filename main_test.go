package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"fetch-rewards/internal"
)

func TestProcessReceipt(t *testing.T) {
	router := setupRouter()

	receipt := internal.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []internal.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: 6.49},
			{ShortDescription: "Emils Cheese Pizza", Price: 12.25},
		},
		Total: 35.35,
	}

	jsonValue, _ := json.Marshal(receipt)
	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response["id"])
	log.Printf("TestProcessReceipt: Received response ID: %s", response["id"])
}

func TestGetPoints(t *testing.T) {
	router := setupRouter()

	// Assuming you have a way to add a receipt to the store for testing
	receiptID := "test-receipt-id"
	internal.ReceiptsStore.Lock()
	internal.ReceiptsStore.Store[receiptID] = 100
	internal.ReceiptsStore.Unlock()

	req, _ := http.NewRequest("GET", "/receipts/"+receiptID+"/points", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]int
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 100, response["points"])
	log.Printf("TestGetPoints: Retrieved points for receipt ID %s: %d", receiptID, response["points"])
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/receipts/process", internal.ProcessReceiptHandler)
	r.GET("/receipts/:receipt_id/points", internal.GetPointsHandler)
	return r
} 