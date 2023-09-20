package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var receipts map[string]Receipt = make(map[string]Receipt)

func main() {
	router := gin.Default()
	router.POST("/receipts/process", postReceipts)
	router.GET("/receipts/:id/points", getPointsByID)

	router.Run("localhost:8080")
}

// Hex IDs borrowed from sosedoff.com
func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func postReceipts(c *gin.Context) {
	var newReceipt Receipt
	if err := c.BindJSON(&newReceipt); err != nil {
		fmt.Println(err)
		return
	}

	id, _ := randomHex(10)
	receipts[id] = newReceipt
	c.IndentedJSON(http.StatusOK, ReceiptID{newReceipt, id})
}

func getPointsByID(c *gin.Context) {
	id := c.Param("id")
	if receipt, ok := receipts[id]; ok {
		c.IndentedJSON(http.StatusOK, gin.H{"points": receipt.GetPoints()})
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Receipt not found"})
}
