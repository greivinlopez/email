package email

import (
	"bytes"
	"log"
	"net/smtp"
	"os"
	"text/template"
)

// EmailUser store user parameters for email server authentication
type EmailUser struct {
	Username    string
	Password    string
	EmailServer string
	Port        string
}

// SmtpTemplateData store the basic information for templating an e-mail
type SmtpTemplateData struct {
	From    string
	To      string
	Subject string
	Body    string
}

var (
	// For Gmail use: &EmailUser{'yourGmailUsername', 'password', 'smtp.gmail.com', 587}
	emailUser *EmailUser
)

// GetSenderInfo returns the authentication settings for the email sender user
// The settings are retrieved from system environment variables
func GetSenderInfo() *EmailUser {
	if emailUser == nil {
		emailUser = &EmailUser{
			os.Getenv("EMAIL_USER"),
			os.Getenv("EMAIL_PASSWORD"),
			os.Getenv("EMAIL_SERVER"),
			os.Getenv("EMAIL_PORT")}
	}
	return emailUser
}

// ConnectToSmtpServer returns a smtp.Auth by authenticating with an
// SMTP Server with the info provided in the parameter
func ConnectToSmtpServer(emailUser EmailUser) smtp.Auth {
	auth := smtp.PlainAuth("", emailUser.Username, emailUser.Password, emailUser.EmailServer)
	return auth
}

// SendMail parses the given template against the given context.
// Context could be a struct or a map such as context := map[string]interface{}{"SenderName": "Greivin"}
// It uses the given smtp.Auth previously created. See "ConnectToSmtpServer".
func SendMail(auth smtp.Auth, emailTemplate string, context interface{}, recipients []string) error {
	var err error
	var doc bytes.Buffer
	//context := &SmtpTemplateData{from, to, subject, body}
	t := template.New("emailTemplate")
	if t, err = t.Parse(emailTemplate); err != nil {
		log.Print("error trying to parse mail template ", err)
	}
	if err = t.Execute(&doc, context); err != nil {
		log.Print("error trying to execute mail template ", err)
	}
	err = smtp.SendMail(emailUser.EmailServer+":"+emailUser.Port,
		auth,
		emailUser.Username,
		//[]string{"nathanleclaire@gmail.com"},
		recipients,
		doc.Bytes())
	if err != nil {
		log.Print("ERROR: attempting to send a mail ", err)
	}

	return nil
}
