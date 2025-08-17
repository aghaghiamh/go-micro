package mailer

import (
	"bytes"
	"html/template"
	"log"
	"mailer/domain"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Service struct {
	Email domain.EmailServer
}

func New(email domain.EmailServer) Service {
	return Service{
		Email: email,
	}
}

func (svc *Service) SendSMTPMessage(msg domain.Message) error {
	if msg.FromEmail == "" {
		msg.FromEmail = svc.Email.FromAddress
	}
	if msg.FromName == "" {
		msg.FromName = svc.Email.FromName
	}
	msg.DataMap = map[string]any{
		"message": msg.Data,
	}

	formattedMsg, htmlErr := svc.buildHTMLMessage(msg)
	if htmlErr != nil {
		return htmlErr
	}

	plainMsg, plainErr := svc.buildPlainMessage(msg)
	if plainErr != nil {
		return plainErr
	}

	// set-up server
	server := mail.NewSMTPClient()
	server.Host = svc.Email.Host
	server.Port = svc.Email.Port
	server.Username = svc.Email.Username
	server.Password = svc.Email.Password
	server.Encryption = svc.getEncryption()
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, smtpErr := server.Connect()
	if smtpErr != nil {
		log.Println("Err: connecting to smtp client server: ", smtpErr)
	}

	// set-up email
	email := mail.NewMSG()
	email.SetFrom(msg.FromEmail).
		AddTo(msg.ToEmail).
		SetSubject(msg.Subject)

	email.SetBody(mail.TextPlain, plainMsg)
	email.AddAlternative(mail.TextHTML, formattedMsg)

	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	err := email.Send(smtpClient)
	if err != nil {
		log.Println("Err: sending email using smtp client")
		return err
	}

	return nil
}

func (svc *Service) buildHTMLMessage(msg domain.Message) (string, error) {
	templateToRender := "./templates/mail.html.gohtml"

	t, tErr := template.New("email-html").ParseFiles(templateToRender)
	if tErr != nil {
		log.Println("Err: parse template: ", tErr)
		return "", tErr
	}

	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		log.Println("Err: executing template: ", err)
		return "", err
	}

	formattedMessage := tpl.String()
	formattedMessage, fErr := svc.inlineCss(formattedMessage)
	if fErr != nil {
		log.Println("Err: inline css formatting: ", fErr)
		return "", fErr
	}

	return formattedMessage, nil
}

func (svc *Service) inlineCss(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}
	prem, premErr := premailer.NewPremailerFromString(s, &options)
	if premErr != nil {
		log.Println("Err: inline css premailer formatting: ", premErr)
		return "", premErr
	}

	html, tErr := prem.Transform()
	if tErr != nil {
		log.Println("Err: transforming prem into html: ", tErr)
		return "", tErr
	}

	return html, nil
}

func (svc *Service) buildPlainMessage(msg domain.Message) (string, error) {
	templateToRender := "./templates/mail.plain.gohtml"

	t, tErr := template.New("email-plain").ParseFiles(templateToRender)
	if tErr != nil {
		log.Println("Err: parse template: ", tErr)
		return "", tErr
	}

	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		log.Println("Err: executing template: ", err)
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

func (svc *Service) getEncryption() mail.Encryption {
	switch svc.Email.Encryption {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
