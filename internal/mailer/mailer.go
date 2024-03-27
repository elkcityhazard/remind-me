package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"text/template"
	"time"

	"github.com/elkcityhazard/remind-me/internal/models"
	"github.com/go-mail/mail/v2"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	Dialer          *mail.Dialer
	Sender          string
	MailerErrorChan chan error
	MailerDataChan  chan *models.EmailData
	MailerDoneChan  chan bool
}

// New returns a new mailer with all necessary config passed into it
func New(host string, port int, username, password, sender string) Mailer {

	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		Dialer:          dialer,
		Sender:          sender,
		MailerErrorChan: make(chan error),
		MailerDataChan:  make(chan *models.EmailData),
		MailerDoneChan:  make(chan bool),
	}
}

func (m Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)

	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	// 	we can execute just the parts in the template file
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.Sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	err = m.Dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil

}

func (m *Mailer) ListenForMail() {

	for {
		select {
		case data := <-m.MailerDataChan:
			log.Println("Received a message...")
			go m.Send(data.Recipient, data.TemplateFile, data.Data)
		case err := <-m.MailerErrorChan:
			log.Println(err)
		case <-m.MailerDoneChan:
			fmt.Println("mailer done signal")
			return
		}
	}

}
