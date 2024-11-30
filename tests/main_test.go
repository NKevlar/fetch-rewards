package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"fetch-rewards/internal"
)

func TestProcessReceipt(t *testing.T) {
	router := gin.Default()
	router.POST("/receipts/process", func(c *gin.Context) {
		var receipt internal.Receipt
		if err := c.ShouldBindJSON(&receipt); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		receiptID := internal.ProcessReceipt(receipt)
		c.JSON(http.StatusOK, gin.H{"id": receiptID})
	})

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
}

func TestGetPoints(t *testing.T) {
	router := gin.Default()
	router.GET("/receipts/:receipt_id/points", func(c *gin.Context) {
		receiptID := c.Param("receipt_id")
		points, exists := internal.GetPoints(receiptID)

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"points": points})
	})

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
}


