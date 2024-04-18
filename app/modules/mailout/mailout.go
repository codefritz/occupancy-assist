package mailout

import (
	"log"
	"net/smtp"
	"os"
)

const MSG_INTRO = "Der aktuelle Buchungskalender zur Ferienwohnung Strandsommer E10.\n\n"
const HEADER_SUBJECT = "Subject: Buchungskalender"
const HEADER_END = "\n\n"

func MailOut(content Report) {

	// smtp server configuration.
	mailProps := mailProperties()

	// Message.
	message := []byte(mailHeaders() + intro() + overview(content.Days) + content.Content)

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

func overview(days int) string {
	return "Belegte Tage: " + fmt.Sprint(ctx) + "\n"
}

func intro() string {
	return MSG_INTRO
}

func mailHeaders() string {
	return HEADER_SUBJECT +
		HEADER_END
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
