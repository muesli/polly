package main

import (
	"io"
	"text/template"

	"github.com/go-gomail/gomail"
)

// EmailTemplate holds all values of an email template
type EmailTemplate struct {
	Subject string
	Text    string
	HTML    string
}

var (
	templates = make(map[string]EmailTemplate)
)

func setupEmailTemplates() {
	templates["invitation"] = EmailTemplate{
		Subject: config.Templates.Invitation.Subject,
		Text:    "Hello {{.Email}}!\n\nYou've been invited to Polly!\nJoin here: " + config.Web.BaseURL + "signup/{{.AuthToken}}",
		HTML:    "Hello <b>{{.Email}}</b>!<br/><br/>You've been invited to Polly!<br/>Join here: " + config.Web.BaseURL + "signup/{{.AuthToken}}",
	}
	templates["moderation_proposal"] = EmailTemplate{
		Subject: config.Templates.ModerationProposal.Subject,
		Text:    "Hello Admin!\n\nA new proposal '{{.Title}}' has been created and awaits moderation!\nClick here: " + config.Web.BaseURL + "proposals/{{.Id}}",
		HTML:    "Hello <b>Admin</b>!<br/><br/>A new proposal <b>{{.Title}}</b> has been created and awaits moderation!<br/>Click here: " + config.Web.BaseURL + "proposals/{{.Id}}",
	}
}

func sendInvitation(user *DbUser) {
	tmpl := templates["invitation"]

	m := gomail.NewMessage()
	m.SetHeader("From", config.Connections.Email.ReplyTo)
	m.SetHeader("To", "muesli@gmail.com") // user.Email)
	m.SetHeader("Subject", tmpl.Subject)
	//	m.SetAddressHeader("Cc", "foo@foobar.com", "Joe")
	//	m.Attach("/tmp/attachment.jpg")

	m.AddAlternativeWriter("text/plain", func(w io.Writer) error {
		t := template.Must(template.New("invitation_text").Parse(tmpl.Text))
		return t.Execute(w, *user)
	})
	m.AddAlternativeWriter("text/html", func(w io.Writer) error {
		t := template.Must(template.New("invitation_html").Parse(tmpl.HTML))
		return t.Execute(w, *user)
	})

	d := gomail.NewDialer(config.Connections.Email.SMTP.Server, config.Connections.Email.SMTP.Port, config.Connections.Email.SMTP.User, config.Connections.Email.SMTP.Password)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func sendModerationRequest(proposal *DbProposal) {
	tmpl := templates["moderation_proposal"]

	m := gomail.NewMessage()
	m.SetHeader("From", config.Connections.Email.ReplyTo)
	m.SetHeader("To", config.Connections.Email.AdminEmail)
	m.SetHeader("Subject", tmpl.Subject)

	m.AddAlternativeWriter("text/plain", func(w io.Writer) error {
		t := template.Must(template.New("moderation_proposal_text").Parse(tmpl.Text))
		return t.Execute(w, *proposal)
	})
	m.AddAlternativeWriter("text/html", func(w io.Writer) error {
		t := template.Must(template.New("moderation_proposal_html").Parse(tmpl.HTML))
		return t.Execute(w, *proposal)
	})

	d := gomail.NewDialer(config.Connections.Email.SMTP.Server, config.Connections.Email.SMTP.Port, config.Connections.Email.SMTP.User, config.Connections.Email.SMTP.Password)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
