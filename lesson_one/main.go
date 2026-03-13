package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Payment struct {
	PaymentId string `json:"paymentId"`
	Amount    int    `json:"amount"`
	Status    string `json:"status"`
}

func main() {

	payments := []Payment{
		{PaymentId: "1", Amount: 100, Status: "Completed"},
		{PaymentId: "2", Amount: 200, Status: "Pending"},
		{PaymentId: "3", Amount: 150, Status: "Failed"},
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Get all payments
	router.GET("/payments", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, payments)
	})

	// Get payment by ID
	router.POST("/payments/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, payment := range payments {
			if id == payment.PaymentId {
				c.JSON(http.StatusOK, payment)
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"message": "Payment not found"})
	})

	// Create a new payment
	router.POST("/payments", func(ctx *gin.Context) {
		var newPayment Payment
		if err := ctx.BindJSON(&newPayment); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}
		payments = append(payments, newPayment)
		ctx.JSON(http.StatusCreated, newPayment)
	})

	// Serve index.html at the root route
	router.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})

	log.Println("Server is running on http://localhost:8080")

	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
