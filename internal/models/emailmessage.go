package models

type EmailData struct {
	Recipient    string
	TemplateFile string
	Data         any
}
