package mailer

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/resend/resend-go/v3"
)

type ResendMailer struct {
	fromEmail string
	client    *resend.Client
}

func NewResendMailer(apiKey, fromEmail string) *ResendMailer {
	client := resend.NewClient(apiKey)

	return &ResendMailer{
		fromEmail: fromEmail,
		client:    client,
	}
}

func (m *ResendMailer) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {

	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return -1, err
	}

	// development testing mode
	from := fmt.Sprintf("%s <%s>", FromName, m.fromEmail)
	recipient := email

	if isSandbox {
		// Domain gratisan dari Resend untuk dev mode
		from = fmt.Sprintf("%s <onboarding@resend.dev>", FromName)

		// Paksa kirim ke email pribadi yang terdaftar di akun Resend milikmu
		recipient = "rafiahsiraprayoga@gmail.com"
	}
	params := &resend.SendEmailRequest{
		From:    from,
		To:      []string{recipient},
		Subject: subject.String(),
		Html:    body.String(),
	}

	// Production Mode

	// params := &resend.SendEmailRequest{
	// 	From:    fmt.Sprintf("%s <%s>", FromName, m.fromEmail),
	// 	To:      []string{email},
	// 	Subject: subject.String(),
	// 	Html:    body.String(),
	// }

	sent, err := m.client.Emails.Send(params)
	if err != nil {
		return -1, fmt.Errorf("resend: failed to send email: %w", err)
	}

	// Resend mengembalikan ID pesan jika sukses
	_ = sent.Id
	return 200, nil
}
