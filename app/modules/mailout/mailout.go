package mailout

import (
	"log"
	"net/smtp"
	"os"
)

func MailOut(content string) {

	// smtp server configuration.
	mailProps := mailProperties()

	// Message.
	message := []byte(headers() + intro() + content)

	// Authentication.
	auth := smtp.PlainAuth(mailProps.from, mailProps.user, mailProps.password, mailProps.smtpHost)

	// Sending email.
	err := smtp.SendMail(mailProps.smtpHost+":"+mailProps.smtpPort, auth, mailProps.from, mailProps.to, message)
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

func mailProperties() MailProperties {
	return MailProperties{
		from:     os.Getenv("MAIL_FROM"),
		password: os.Getenv("MAIL_SCRT"),
		user:     os.Getenv("MAIL_USER"),
		to:       []string{os.Getenv("MAIL_TO")},
		smtpHost: os.Getenv("MAIL_HOST"),
		smtpPort: "25",
	}
}

type MailProperties struct {
	from     string
	password string
	user     string
	to       []string
	smtpHost string
	smtpPort string
}
