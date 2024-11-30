package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"fetch-rewards/internal"
	_ "fetch-rewards/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func processReceipt(c *gin.Context) {
	var receipt internal.Receipt
	if err := c.ShouldBindJSON(&receipt); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	receiptID := internal.ProcessReceipt(receipt)
	log.Printf("Processed receipt with ID: %s", receiptID)
	c.JSON(http.StatusOK, gin.H{"id": receiptID})
}

func getPoints(c *gin.Context) {
	receiptID := c.Param("receipt_id")
	log.Printf("Fetching points for receipt ID: %s", receiptID)
	points, exists := internal.GetPoints(receiptID)

	if !exists {
		log.Printf("Receipt ID %s not found", receiptID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
		return
	}

	log.Printf("Retrieved points for receipt ID %s: %d", receiptID, points)
	c.JSON(http.StatusOK, gin.H{"points": points})
}

func main() {
	log.Println("Starting the server on port 8080...")
	r := gin.Default()

	r.POST("/receipts/process", processReceipt)
	r.GET("/receipts/:receipt_id/points", getPoints)

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}