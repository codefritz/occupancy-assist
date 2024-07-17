package mailout

import (
	"bytes"
	"github.com/codefritz/occupancy-assist/app/modules/models"
	"log"
	"net/smtp"
	"os"
	"text/template"
)

func MailOut(content models.Report) {

	// smtp server configuration.
	mailProps := smtpMailProperties()
	body := MailBody{Days: content.Days}
	buf, err := createMail(body)
	if err != nil {
		log.Println(err)
		return
	}

	// Message.
	message := []byte(buf.String() + content.Content)

	// Authentication.
	auth := smtp.PlainAuth(mailProps.from, mailProps.user, mailProps.password, mailProps.smtpHost)

	// Sending email.
	err = smtp.SendMail(mailProps.smtpHost+":"+mailProps.smtpPort, auth, mailProps.from, mailProps.to, message)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Email Sent Successfully!")
}

func createMail(body MailBody) (bytes.Buffer, error) {
	bodyTemplateContent, err := os.ReadFile("email_template.txt") // Specify the correct path
	if err != nil {
		log.Println(err)
		return bytes.Buffer{}, err
	}

	var buf bytes.Buffer
	bodyTemplate, err := template.New("mail").Parse(string(bodyTemplateContent))
	if err != nil {
		log.Println(err)
		return bytes.Buffer{}, err
	}
	err = bodyTemplate.Execute(&buf, body)
	if err != nil {
		log.Println(err)
		return bytes.Buffer{}, err
	}
	return buf, err
}

type MailBody struct {
	Days int
}

func smtpMailProperties() MailProperties {
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
