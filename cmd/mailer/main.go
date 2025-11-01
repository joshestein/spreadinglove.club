package main

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"spreadlove/internal/database"
)

func main() {
	db, queries, err := database.Setup("./love.db")
	if err != nil {
		log.Fatal("Failed to setup database:", err)
	}
	defer db.Close()

	ctx := context.Background()
	msgs, err := queries.ListPendingMessages(ctx)
	if err != nil {
		log.Fatal("Error fetching pending messages:", err)
	}

	if len(msgs) == 0 {
		fmt.Println("No pending messages to review")
		return
	}

	from := os.Getenv("EMAIL_SENDER")
	recipient := os.Getenv("EMAIL_RECIPIENT")
	password := os.Getenv("EMAIL_PASSWORD")

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := "587"

	to := []string{recipient}

	body := "<html><body>Go to <a href='https://spreadlove.club/admin'>the admin page</a> to check them out.</body></html>"
	message := []byte("From: " + from + "\n" +
		"Subject: There are pending love messages for you to review\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		body + "\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return
	}

	fmt.Println("Email sent successfully!")
}
