package mailman

import (
	"bytes"
	"fmt"
	"log"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/muesli/gomail"
	"github.com/mxk/go-imap/imap"

	"github.com/muesli/polly/api/db"
	"github.com/muesli/polly/api/utils"
)

var (
	context *db.PollyContext
)

// SetupMailmanContext sets the context
func SetupMailmanContext(ctx *db.PollyContext) {
	context = ctx
}

func sendMail(tos []string, from, subject, body, mid, contenttype, contenttypetransfer string) {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("Subject", subject)
	m.SetHeader("Message-ID", mid)
	m.SetHeader("Errors-To", context.Config.Connections.Email.Mailman.BounceAddress)
	m.SetHeader("X-BeenThere", context.Config.Connections.Email.Mailman.Address)
	m.SetAddressHeader("To", context.Config.Connections.Email.Mailman.Address, context.Config.Connections.Email.Mailman.Name)
	m.SetAddressHeader("Envelope-Sender", context.Config.Connections.Email.Mailman.BounceAddress, context.Config.Connections.Email.Mailman.Name)
	m.SetAddressHeader("List-Id", context.Config.Connections.Email.Mailman.Address, context.Config.Connections.Email.Mailman.Name)

	if len(contenttype) > 0 {
		m.SetHeader("Content-Type", contenttype)
	}
	if len(contenttypetransfer) > 0 {
		m.SetHeader("Content-Transfer-Encoding", contenttypetransfer)
	}

	m.SetRawBody(body)

	d := gomail.NewDialer(context.Config.Connections.Email.SMTP.Server, context.Config.Connections.Email.SMTP.Port,
		context.Config.Connections.Email.SMTP.User, context.Config.Connections.Email.SMTP.Password)
	s, err := d.Dial()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	err = s.Send(context.Config.Connections.Email.Mailman.Address, tos, m)
	if err != nil {
		panic(err)
	}
}

// RunProposalLoop watches freshly started proposals and sends out a reminder email
func RunProposalLoop() {
	for {
		proposals, err := context.LoadAllProposals()
		if err != nil {
			log.Println("ERROR:", err)
			time.Sleep(5 * time.Minute)
			continue
		}

		for _, p := range proposals {
			if p.Moderated && !p.StartTrigger && p.Starts.Before(time.Now()) {
				log.Printf("Proposal '%s' started. Sending out reminder emails...\n", p.Title)
				utils.SendProposalStarted(p)

				p.StartTrigger = true
				p.Update(context)
			}
		}

		time.Sleep(5 * time.Minute)
	}
}

// RunLoop fetches mail and delivers them to recipients - forever
func RunLoop() {
	var (
		c   *imap.Client
		cmd *imap.Command
		rsp *imap.Response
		err error
	)

	// Connect to the server
	if context.Config.Connections.Email.IMAP.Port == 993 {
		c, err = imap.DialTLS(context.Config.Connections.Email.IMAP.Server+":"+strconv.FormatInt(int64(context.Config.Connections.Email.IMAP.Port), 10), nil)
	} else {
		c, err = imap.Dial(context.Config.Connections.Email.IMAP.Server + ":" + strconv.FormatInt(int64(context.Config.Connections.Email.IMAP.Port), 10))
	}
	if err != nil || c == nil || c.Data == nil {
		panic(err)
	}

	// Print server greeting (first response in the unilateral server data queue)
	fmt.Println("IMAP Server says hello:", c.Data[0].Info)
	c.Data = nil

	// Enable encryption, if supported by the server
	if c.Caps["STARTTLS"] {
		c.StartTLS(nil)
	}

	// Authenticate
	if c.State() == imap.Login {
		c.Login(context.Config.Connections.Email.IMAP.User, context.Config.Connections.Email.IMAP.Password)
	}

	// List all top-level mailboxes, wait for the command to finish
	cmd, _ = imap.Wait(c.List("", "%"))
	if cmd == nil {
		return
	}

	// Print mailbox information
	fmt.Println("\nTop-level mailboxes:")
	for _, rsp = range cmd.Data {
		fmt.Println("|--", rsp.MailboxInfo().Name)
	}

	// Check for new unilateral server data responses
	for _, rsp = range c.Data {
		fmt.Println("Server data:", rsp)
	}
	c.Data = nil

	mm, _ := context.GetMailman(context.Config.Connections.Email.Mailman.Address)
	for {
		// Open a mailbox (synchronous command - no need for imap.Wait)
		c.Select("INBOX", true)
		fmt.Printf("Mailbox status: %s (msgs: %d, last-seen: %d)\n", c.Mailbox.Name, c.Mailbox.Messages, mm.LastSeen)

		if mm.LastSeen == 0 {
			mm.LastSeen = uint64(c.Mailbox.Messages)
		}

		// Fetch new mails
		startFromID := mm.LastSeen + 1
		set, _ := imap.NewSeqSet("")
		set.Add(strconv.FormatUint(startFromID, 10) + ":*")
		cmd, _ = c.UIDFetch(set, "RFC822.HEADER", "RFC822.TEXT")

		// Process responses while the command is running
		fmt.Println("\nChecking mailman INBOX")
		for cmd.InProgress() {
			// Wait for the next response (no timeout)
			c.Recv(-1)

			// Process command data
			for _, rsp = range cmd.Data {
				if uint64(rsp.MessageInfo().UID) < startFromID {
					continue
				}

				header := imap.AsBytes(rsp.MessageInfo().Attrs["RFC822.HEADER"])
				body := string(imap.AsBytes(rsp.MessageInfo().Attrs["RFC822.TEXT"]))
				if msg, _ := mail.ReadMessage(bytes.NewReader(header)); msg != nil {
					from := msg.Header.Get("From")
					subj := msg.Header.Get("Subject")
					mid := msg.Header.Get("Message-ID")
					contenttype := msg.Header.Get("Content-Type")
					contenttypetrasnfer := msg.Header.Get("Content-Transfer-Encoding")

					fmt.Println("|-- From", from)
					fmt.Println("|-- Subject", subj)
					fmt.Println("|-- ID", mid)

					users, err := context.LoadAllUsers()
					if err != nil {
						panic(err)
					}
					tos := []string{}
					for _, user := range users {
						if !user.Activated {
							continue
						}
						tos = append(tos, user.Email)
					}

					if msg.Header.Get("X-BeenThere") != context.Config.Connections.Email.Mailman.Address && !strings.Contains(body, "X-BeenThere: "+context.Config.Connections.Email.Mailman.Address) {
						if len(tos) > 0 {
							sendMail(tos, from, subj, body, mid, contenttype, contenttypetrasnfer)
						}
					} else {
						fmt.Println("IGNORING MESSAGE!")
					}
					mm.LastSeen = uint64(rsp.MessageInfo().UID)
					mm.Update(context)
				}
			}

			cmd.Data = nil

			// Process unilateral server data
			for _, rsp = range c.Data {
				// fmt.Println("Server data:", rsp)
			}
			c.Data = nil
		}

		// Check command completion status
		if rsp, err := cmd.Result(imap.OK); err != nil {
			if err == imap.ErrAborted {
				fmt.Println("Fetch command aborted")
			} else {
				fmt.Println("Fetch error:", rsp.Info)
				return
			}
		}

		time.Sleep(1 * time.Minute)
	}
}
