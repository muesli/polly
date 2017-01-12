package utils

import (
	"io"
	"text/template"

	"github.com/muesli/polly/api/config"
	"github.com/muesli/polly/api/db"

	"github.com/go-gomail/gomail"
)

var (
	templates = make(map[string]config.EmailTemplate)
	settings  config.Data
)

// SetupEmailTemplates compiles the email templates
func SetupEmailTemplates(c config.Data) {
	settings = c
	templates["invitation"] = config.EmailTemplate{
		Subject: c.App.Templates.Invitation.Subject,
		Text:    "Hello {{.Email}}!\n\nYou've been invited to Polly!\nJoin here: " + c.Web.BaseURL + "signup/{{.AuthToken}}",
		HTML:    "Hello <b>{{.Email}}</b>!<br/><br/>You've been invited to Polly!<br/>Join here: " + c.Web.BaseURL + "signup/{{.AuthToken}}",
	}
	templates["moderation_proposal"] = config.EmailTemplate{
		Subject: c.App.Templates.ModerationProposal.Subject,
		Text:    "Hello Admin!\n\nA new proposal '{{.Title}}' has been created and awaits moderation!\nClick here: " + c.Web.BaseURL + "proposals/{{.ID}}",
		HTML:    "Hello <b>Admin</b>!<br/><br/>A new proposal <b>{{.Title}}</b> has been created and awaits moderation!<br/>Click here: " + c.Web.BaseURL + "proposals/{{.ID}}",
	}
}

// SendInvitation sends out an email, inviting a user to join polly
func SendInvitation(user *db.User) {
	tmpl := templates["invitation"]

	m := gomail.NewMessage()
	m.SetHeader("From", settings.Connections.Email.ReplyTo)
	m.SetHeader("To", settings.Connections.Email.
		AdminEmail) // FIXME: change to user.Email in production
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

	d := gomail.NewDialer(settings.Connections.Email.SMTP.Server, settings.Connections.Email.SMTP.Port,
		settings.Connections.Email.SMTP.User, settings.Connections.Email.SMTP.Password)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

// SendModerationRequest sends out an email to the admin, asking for moderation of a newly posted proposal
func SendModerationRequest(proposal *db.Proposal) {
	tmpl := templates["moderation_proposal"]

	m := gomail.NewMessage()
	m.SetHeader("From", settings.Connections.Email.ReplyTo)
	m.SetHeader("To", settings.Connections.Email.AdminEmail)
	m.SetHeader("Subject", tmpl.Subject)

	m.AddAlternativeWriter("text/plain", func(w io.Writer) error {
		t := template.Must(template.New("moderation_proposal_text").Parse(tmpl.Text))
		return t.Execute(w, *proposal)
	})
	m.AddAlternativeWriter("text/html", func(w io.Writer) error {
		t := template.Must(template.New("moderation_proposal_html").Parse(tmpl.HTML))
		return t.Execute(w, *proposal)
	})

	d := gomail.NewDialer(settings.Connections.Email.SMTP.Server, settings.Connections.Email.SMTP.Port,
		settings.Connections.Email.SMTP.User, settings.Connections.Email.SMTP.Password)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
