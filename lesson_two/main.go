package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Payment struct {
	PaymentId string `json:"paymentId" validate:"required"`
	Amount    int    `json:"amount" validate:"required,gt=0"`
	Status    string `json:"status" validate:"required,oneof=Pending Completed Failed"`
}

type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

var validate *validator.Validate

func main() {

	validate = validator.New(validator.WithRequiredStructEnabled())

	payments := []Payment{
		{PaymentId: "1", Amount: 100, Status: "Completed"},
		{PaymentId: "2", Amount: 200, Status: "Pending"},
		{PaymentId: "3", Amount: 150, Status: "Failed"},
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Get all payments
	router.GET("/payments", func(c *gin.Context) {
		c.JSON(http.StatusOK, payments)
	})

	// Get payment by ID
	router.GET("/payments/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, payment := range payments {
			if id == payment.PaymentId {
				c.JSON(http.StatusOK, payment)
				return
			}
		}
		c.Header("Content-Type", "application/problem+json")
		c.JSON(http.StatusNotFound, ProblemDetails{
			Type:     "/problem/payment-not-found",
			Title:    "Payment Not Found",
			Detail:   fmt.Sprintf("No payment found with ID '%s'", id),
			Instance: fmt.Sprintf("/payments/%s", id),
		})
	})

	// Create a new payment
	router.POST("/payments", func(ctx *gin.Context) {
		var newPayment Payment
		var err error

		if err := ctx.BindJSON(&newPayment); err != nil {
			ctx.Header("Content-Type", "application/problem+json")
			ctx.JSON(http.StatusBadRequest, ProblemDetails{
				Type:     "/problem/invalid-json",
				Title:    "Invalid JSON",
				Detail:   "Failed to parse request body as JSON",
				Instance: "/payments",
			})
			return
		}
		err = validate.Struct(newPayment)
		if err != nil {
			ctx.Header("Content-Type", "application/problem+json")
			ctx.JSON(http.StatusBadRequest, ProblemDetails{
				Type:     "/problem/validation-error",
				Title:    "Validation Error",
				Detail:   err.Error(),
				Instance: "/payments",
			})
			return
		}
		// Check for duplicate PaymentId
		for _, payment := range payments {
			if payment.PaymentId == newPayment.PaymentId {
				ctx.Header("Content-Type", "application/problem+json")
				ctx.JSON(http.StatusConflict, ProblemDetails{
					Type:     "/problem/duplicate-payment-id",
					Title:    "Duplicate Payment ID",
					Detail:   fmt.Sprintf("A payment with ID '%s' already exists", newPayment.PaymentId),
					Instance: "/payments",
				})
				return
			}
		}
		payments = append(payments, newPayment)
		ctx.JSON(http.StatusCreated, newPayment)
	})

	log.Println("Server is running on http://localhost:3000")

	err := router.Run(":3000")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
