package mailout

import (
	"bytes"
	"github.com/codefritz/occupancy-assist/app/modules/models"
	"log"
	"net/smtp"
	"os"
	"text/template"
)

type loginAuth struct {
    username, password string
}

func LoginAuth(username, password string) smtp.Auth {
    return &loginAuth{username, password}
}

func MailOut(content models.Report) {

	// smtp server configuration.
	mailProps := smtpMailProperties()

  // Authentication.
  conn, err := net.Dial("tcp", mailProps.smtpHost ":" mailProps.smtpPort)
  if err != nil {
      log.Println(err)
  }

  c, err := smtp.NewClient(conn, mailProps.smtpHost)
  if err != nil {
      println(err)
  }

  tlsconfig := &tls.Config{
      ServerName: mailProps.smtpHost,
  }

  if err = c.StartTLS(tlsconfig); err != nil {
      log.Println(err)
  }

  auth := LoginAuth(mailProps.user, mailProps.password)

  if err = c.Auth(auth); err != nil {
      log.Println(err)
  }

	body := MailBody{Days: content.Days}
	buf, err := createMail(body)
	if err != nil {
		log.Println(err)
		return
	}

	// Message.
	message := []byte(buf.String() + content.Details)

	// Sending email.
	err = smtp.SendMail(mailProps.smtpHost+":"+mailProps.smtpPort, auth, mailProps.from, mailProps.to, message)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Email Sent Successfully!")
}

func createMail(body MailBody) (bytes.Buffer, error) {
	bodyTemplateContent, err := os.ReadFile("modules/mailout/email_template.txt") // Specify the correct path
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
		smtpPort: os.Getenv("MAIL_PORT"),
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
