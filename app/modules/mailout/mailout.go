package mailout

import (
	"bytes"
	"crypto/tls"
	"fmt"
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

	// Create properly formatted email message
	message := createEmailMessage(mailProps, buf.String()+content.Details)

	// Send email using Outlook SMTP with TLS
	err = sendEmailViaOutlook(mailProps, message)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
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

// createEmailMessage creates a properly formatted email message with headers
func createEmailMessage(mailProps MailProperties, body string) []byte {
	// Extract subject from body (if it starts with "Subject:")
	subject := "Buchungskalender"
	bodyContent := body

	// Check if body starts with Subject: line
	if len(body) > 8 && body[:8] == "Subject:" {
		lines := bytes.Split([]byte(body), []byte("\n"))
		if len(lines) > 0 {
			subject = string(lines[0][9:]) // Remove "Subject: " prefix
			bodyContent = string(bytes.Join(lines[1:], []byte("\n")))
		}
	}

	// Create email headers
	headers := make(map[string]string)
	headers["From"] = mailProps.from
	headers["To"] = mailProps.to[0]
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=utf-8"

	// Build message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + bodyContent

	return []byte(message)
}

// sendEmailViaOutlook sends email using Outlook SMTP with proper TLS configuration
func sendEmailViaOutlook(mailProps MailProperties, message []byte) error {
	// Outlook/Office365 SMTP settings
	smtpServer := mailProps.smtpHost
	if smtpServer == "" {
		smtpServer = "smtp.office365.com"
	}

	smtpPort := mailProps.smtpPort
	if smtpPort == "" {
		smtpPort = "587"
	}

	// Create TLS config
	tlsConfig := &tls.Config{
		ServerName: smtpServer,
	}

	// Connect to SMTP server
	conn, err := tls.Dial("tcp", smtpServer+":"+smtpPort, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, smtpServer)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Close()

	// Authenticate using PLAIN auth (works for App Passwords)
	auth := smtp.PlainAuth("", mailProps.user, mailProps.password, smtpServer)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	// Set sender
	if err = client.Mail(mailProps.from); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}

	// Set recipients
	for _, to := range mailProps.to {
		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("failed to set recipient %s: %v", to, err)
		}
	}

	// Send message
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %v", err)
	}

	_, err = writer.Write(message)
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	return nil
}
