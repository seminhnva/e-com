package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendgrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *SendGridMailer) Send(tmpFile string, username, email string, data any, isSandbox bool) error {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	tmpl, err := template.ParseFS(FS, "templates/"+tmpFile)

	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "Body", data)
	if err != nil {
		return err
	}

	msg := mail.NewSingleEmail(from, subject.String(), to, "", body.String())
	msg.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	for i := 0; i < maxRetries; i++ {
		res, err := m.client.Send(msg)
		if err != nil {
			log.Printf("Fail to send email to %v, attempt %d of %d", email, i+1, maxRetries)
			log.Printf("Error: %v", err.Error())
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		log.Printf("Email send with status code %v", res.StatusCode)
		return nil
	}
	return fmt.Errorf("failed send email after %d attempts", maxRetries)
}
