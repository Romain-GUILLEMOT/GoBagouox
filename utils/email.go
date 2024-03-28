package utils

import (
	"bytes"
	"crypto/tls"
	"embed"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

//go:embed views/*
var t embed.FS

func SendEmail(recipient string, subject string, templateName string, data map[string]string) {
	smtpHost := os.Getenv("MAIL_HOST")
	smtpPort := os.Getenv("MAIL_PORT")
	sender := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	tmpl, err := template.ParseFS(t, "views/emails/"+templateName+".html")
	if err != nil {
		Error("Can't load an emails template.", err, 3)
		return
	}
	Debug("New Email view parsed", 3)

	// Form the body of emails
	var body bytes.Buffer
	data["Name"] = subject
	// "Execute" the template and write its output to our bytes buffer.
	err = tmpl.Execute(&body, data)
	if err != nil {
		Error("Can't convert emails template into bytes.", err, 3)
		return
	}
	Debug("New emails view transformed", 3)

	// Form MIME headers
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	mime = fmt.Sprintf("From: %s<%s>\nTo: %s\nSubject: %s\n%s\n\n", os.Getenv("MAIL_FROM_NAME"), sender, recipient, subject, mime)
	message := []byte(mime + "\r\n" + body.String())

	// Form the authentication
	auth := smtp.PlainAuth("", sender, password, smtpHost)

	// Connecter au serveur SMTP
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsconfig)
	if err != nil {
		Error("Can't connect to the SMTP server.", err, 3)
		return
	}

	c, errClient := smtp.NewClient(conn, smtpHost)
	if errClient != nil {
		Error("Can't create the SMTP client.", errClient, 3)
		return
	}

	// Authentifier
	if err := c.Auth(auth); err != nil {
		Error("SMTP server rejected authentication.", err, 3)
		return
	}

	// Set the sender and recipient
	if err := c.Mail(sender); err != nil {
		Error("Failed while setting the sender.", err, 3)
		return
	}
	if err := c.Rcpt(recipient); err != nil {
		Error("Failed while setting the recipient.", err, 3)
		return
	}

	// Send the emails body
	wc, err := c.Data()
	if err != nil {
		Error("Failed while preparing to send the emails body.", err, 3)
		return
	}
	_, err = wc.Write(message)
	if err != nil {
		Error("Failed while writing the emails body.", err, 3)
		return
	}
	err = wc.Close()
	if err != nil {
		Error("Failed while closing the connection after sending the emails.", err, 3)
		return
	}

	// Send the QUIT command and close the connection
	err = c.Quit()
	if err != nil {
		Error("Failed while quitting the SMTP session.", err, 3)
		return
	}

	Info("Email sent with subject: "+subject+" to "+recipient, 3)
}
