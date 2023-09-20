package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            float64 `json:"price,string"`
}

type Receipt struct {
	Retailer     string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
	Items        []Item  `json:"items"`
	Total        float64 `json:"total,string"`
}

type ReceiptID struct {
	Receipt Receipt `json:"-"`
	ID      string  `json:"id"`
}

func (receipt Receipt) GetPurchaseDateTime() time.Time {
	dt, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", receipt.PurchaseDate, receipt.PurchaseTime))
	if err != nil {
		// In an actual application, you'd definitely want to handle this properly
		// But at least for a quick demo API, I'm just using time.Now() as a default
		return time.Now()
	}
	return dt
}

func (receipt Receipt) GetTotal() float64 {
	var sum float64 = 0
	for _, item := range receipt.Items {
		sum += item.Price
	}

	return sum
}

func (receipt Receipt) GetPoints() int {
	// Add 1 point for each alphanumeric character in the retailer's name
	points := CountAlphanumeric(receipt.Retailer)

	// Add 50 points if the total is a whole number
	if receipt.GetTotal() == math.Floor(receipt.GetTotal()) {
		points += 50
	}

	// Add 25 points if the total is a multiple of 0.25
	if int(receipt.GetTotal()*100)%25 == 0 {
		points += 25
	}

	// Add 5 points for every 2 items
	points += 5 * (len(receipt.Items) / 2)

	// Add 6 points if the purchase date is odd
	if receipt.GetPurchaseDateTime().Day()%2 == 1 {
		points += 6
	}

	// Add 10 points if purchased between 2:00 and 4:00 PM
	if receipt.GetPurchaseDateTime().Hour() >= 14 && receipt.GetPurchaseDateTime().Hour() < 17 {
		points += 10
	}

	// For each item, if the length of the trimmed name is a multiple of 3,
	// add ceil(0.2*price) points
	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			points += int(math.Ceil(0.2 * item.Price))
		}
	}

	return points
}

func CountAlphanumeric(s string) int {
	count := 0
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			count++
		}
	}
	return count
}
