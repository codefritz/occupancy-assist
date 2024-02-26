package mailout

import (
	"log"
	"net/smtp"
	"os"
)

func MailOut(content string) {

	// Sender data.
	// from := "wp1058313-andre"
	from := os.Getenv("MAIL_FROM")
	password := os.Getenv("MAIL_SCRT")
	user := os.Getenv("MAIL_USER")

	// Receiver email address.
	to := []string{
		os.Getenv("MAIL_TO"),
	}

	// smtp server configuration.
	smtpHost := os.Getenv("MAIL_HOST")
	smtpPort := "25"

	// Message.
	message := []byte(headers() + intro() + content)

	// Authentication.
	auth := smtp.PlainAuth(from, user, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Email Sent Successfully!")
}

func intro() string {
	return "Der aktuelle Buchungskalender zur Ferienwohnung Strandsommer E10.\n\n"
}

func headers() string {
	return "Subject: Buchungskalender\n\n"
}
