package mailer

import (
	"bytes"
	"embed"
	"log"
	"sync"
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

func (m Mailer) Send(recipient, templateFile string, data any) {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)

	if err != nil {
		m.MailerErrorChan <- err
		return
	}

	subject := new(bytes.Buffer)
	// 	we can execute just the parts in the template file
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		m.MailerErrorChan <- err
		return
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		m.MailerErrorChan <- err
		return
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		m.MailerErrorChan <- err
		return
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.Sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	err = m.Dialer.DialAndSend(msg)
	if err != nil {
		m.MailerErrorChan <- err
		return
	}

}

func (m *Mailer) ListenForMail(wg *sync.WaitGroup) {

	defer wg.Done()

	for {
		select {
		case data := <-m.MailerDataChan:
			go m.Send(data.Recipient, data.TemplateFile, data.Data)
		case err := <-m.MailerErrorChan:
			log.Println(err)
		case <-m.MailerDoneChan:
			return
		}
	}

}
