package internal

import (
	"log"
	"strings"
	"sync"
	"time"
	"math"
	"github.com/google/uuid"
)

type Receipt struct {
	Retailer     string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
	Items        []Item  `json:"items"`
	Total        float64 `json:"total,string"`
}

type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            float64 `json:"price,string"`
}

var ReceiptsStore = struct {
	sync.RWMutex
	Store map[string]int
}{Store: make(map[string]int)}


func CalculatePoints(receipt Receipt) int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name
	retailerPoints := 0
	for _, c := range receipt.Retailer {
		if isAlnum(c) {
			retailerPoints++
		}
	}
	points += retailerPoints
	log.Printf("Retailer points: %d (retailer name: %s)\n", retailerPoints, receipt.Retailer)

	// Rule 2: 50 points if the total is a round dollar amount with no cents
	if math.Mod(receipt.Total, 1) == 0 {
		points += 50
		log.Printf("Added 50 points for round dollar total: %.2f\n", receipt.Total)
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if math.Mod(receipt.Total, 0.25) == 0 {
		points += 25
		log.Printf("Added 25 points for total being a multiple of 0.25: %.2f\n", receipt.Total)
	}

	// Rule 4: 5 points for every two items on the receipt
	itemPairPoints := (len(receipt.Items) / 2) * 5
	points += itemPairPoints
	log.Printf("Added %d points for %d items (5 points per pair)\n", itemPairPoints, len(receipt.Items))

	// Rule 5: Points for item descriptions that are multiples of 3
	for _, item := range receipt.Items {
		description := strings.TrimSpace(item.ShortDescription)
		if len(description)%3 == 0 {
			itemPoints := int(math.Ceil(item.Price * 0.2))
			points += itemPoints
			log.Printf("Added %d points for item '%s' (price: %.2f, length: %d)\n", itemPoints, description, item.Price, len(description))
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd
	if date, err := time.Parse("2006-01-02", receipt.PurchaseDate); err == nil {
		if date.Day()%2 == 1 {
			points += 6
			log.Printf("Added 6 points for odd purchase day: %d\n", date.Day())
		}
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm
	if t, err := time.Parse("15:04", receipt.PurchaseTime); err == nil {
		if t.Hour() == 14 {
			points += 10
			log.Printf("Added 10 points for purchase time between 2:00 PM and 4:00 PM: %s\n", receipt.PurchaseTime)
		}
	}

	log.Printf("Total points: %d\n", points)
	return points
}

func isAlnum(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func ProcessReceipt(receipt Receipt) string {
	receiptID := uuid.New().String()
	points := CalculatePoints(receipt)

	ReceiptsStore.Lock()
	ReceiptsStore.Store[receiptID] = points
	ReceiptsStore.Unlock()

	log.Printf("Stored receipt ID %s with points %d", receiptID, points)
	return receiptID
}

func GetPoints(receiptID string) (int, bool) {
	ReceiptsStore.RLock()
	points, exists := ReceiptsStore.Store[receiptID]
	ReceiptsStore.RUnlock()

	if exists {
		log.Printf("Found points for receipt ID %s: %d", receiptID, points)
	} else {
		log.Printf("No points found for receipt ID %s", receiptID)
	}

	return points, exists
} 