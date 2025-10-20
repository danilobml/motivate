package services

import (
	"net/smtp"
	"os"

	"github.com/danilobml/motivate/internal/errs"
)

func SendMail(to []string, subject string, body string) error {
	if os.Getenv("FROM_EMAIL") == "" ||
		os.Getenv("FROM_EMAIL_PASSWORD") == "" ||
		os.Getenv("FROM_EMAIL_SMTP") == "" ||
		os.Getenv("SMTP_ADDR") == "" {
		return errs.ErrMailServiceDisabled
	}

	auth := smtp.PlainAuth(
		"",
		os.Getenv("FROM_EMAIL"),
		os.Getenv("FROM_EMAIL_PASSWORD"),
		os.Getenv("FROM_EMAIL_SMTP"),
	)

	message := "Subject: " + subject + "\n" + body

	err := smtp.SendMail(
		os.Getenv("SMTP_ADDR"),
		auth,
		os.Getenv("FROM_EMAIL"),
		to,
		[]byte(message),
	)

	return err
}
