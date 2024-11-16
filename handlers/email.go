package handlers

import (
	"log"
	"net/http"
	"os"

	"portfolio-backend/models"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

func SendEmail(c *gin.Context) {
	var details models.EmailDetails
	if err := c.ShouldBindJSON(&details); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", os.Getenv("SMTP_EMAIL_FOR_PORTFOLIO"))
	mail.SetHeader("To", os.Getenv("SMTP_EMAIL_FOR_PORTFOLIO"))
	mail.SetHeader("Subject", "New Contact Form Submission")
	mail.SetBody("text/plain",
		"New contact form Portfolio Submission:\n"+
			"First Name: "+details.FirstName+"\n"+
			"Last Name: "+details.LastName+"\n"+
			"Email: "+details.Email+"\n"+
			"Phone: "+details.Phone+"\n"+
			"Message: "+details.Message,
	)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("SMTP_EMAIL_FOR_PORTFOLIO"), os.Getenv("SMTP_PASSWORD_FOR_PORTFOLIO"))
	if err := dialer.DialAndSend(mail); err != nil {
		log.Printf("Error sending email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}
