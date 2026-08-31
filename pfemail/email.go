package pfemail

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type Etype struct {
	uname   string
	pword   string
	mserver string
	mport   string
	c       *client.Client
}

func (e *Etype) Connect(s1, s2, s3, s4 string) {
	e.mserver = s1
	e.uname = s2
	e.pword = s3
	e.mport = s4
}

func (e *Etype) Checkemaillogin() error {
	err := e.Login()
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}
	e.c.Logout()
	return nil
}

func (e *Etype) Login() error {
	if e.pword == "" {
		e.pword = fyne.CurrentApp().Preferences().StringWithFallback("epassword", "")
		e.mserver = fyne.CurrentApp().Preferences().StringWithFallback("eserver", "")
		e.uname = fyne.CurrentApp().Preferences().StringWithFallback("euser", "")
		e.mport = fyne.CurrentApp().Preferences().StringWithFallback("eport", "993")
		if e.mport == "" {
			e.mport = "993"
		}
	}

	// Anslut till servern
	var err error
	e.c, err = client.DialTLS(fmt.Sprintf("%s:%s", e.mserver, e.mport), nil)
	if err != nil {
		log.Println("#1 Login DialTLS error:", err)
		return err
	}
	log.Println("Connected to server, checking login")

	// Logga in
	err = e.c.Login(e.uname, e.pword)
	if err != nil {
		log.Println("#2 Login failed:", err)
		e.c.Close()
		return err
	}
	log.Println("Logged in successfully")
	return nil
}

func (e *Etype) ListMailboxes() []string {
	if err := e.Login(); err != nil {
		return nil
	}
	defer e.c.Logout()

	var mb []string
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)

	go func() {
		done <- e.c.List("", "*", mailboxes)
	}()

	for m := range mailboxes {
		mb = append(mb, m.Name)
	}

	if err := <-done; err != nil {
		log.Println("#1 ListMailboxes error:", err)
		return nil
	}

	return mb
}

func (e *Etype) Createmailboxfolder(f string) error {
	if err := e.Login(); err != nil {
		return err
	}
	defer e.c.Logout()

	err := e.c.Create("INBOX." + f)
	if err != nil {
		log.Println("#1 Create mailbox error:", err)
		return err
	}

	err = e.c.Expunge(nil)
	if err != nil {
		log.Println("#2 Expunge error:", err)
		return err
	}

	return nil
}

func (e *Etype) Getallsmsmail() []*imap.Message {
	if err := e.Login(); err != nil {
		log.Println("#0 Getallsmsmail Login error:", err)
		return nil
	}
	defer e.c.Logout()

	var imsgs []*imap.Message
	_, err := e.c.Select("INBOX", false)
	if err != nil {
		log.Println("#1 Getallsmsmail Select error:", err)
		return nil
	}

	criteria := imap.NewSearchCriteria()
	criteria.Text = []string{"TEST 123"}
	ids, err := e.c.Search(criteria)
	if err != nil {
		log.Println("#2 Getallsmsmail Search error:", err)
		return nil
	}

	log.Println("IDs found:", ids)
	if len(ids) > 0 {
		seqset := new(imap.SeqSet)
		seqset.AddNum(ids...)
		items := []imap.FetchItem{imap.FetchEnvelope}
		messages := make(chan *imap.Message, 10)
		done := make(chan error, 1)

		go func() {
			done <- e.c.Fetch(seqset, items, messages)
		}()

		for msg := range messages {
			imsgs = append(imsgs, msg)
			if msg.Envelope != nil {
				log.Println("* " + msg.Envelope.Subject)
			}
		}

		if err := <-done; err != nil {
			log.Println("#3 Getallsmsmail Fetch error:", err)
			return nil
		}

		if err := e.c.Move(seqset, "INBOX.sms"); err != nil {
			log.Println("#4 Getallsmsmail Move error:", err)
		}
	}

	return imsgs
}

func (e *Etype) Moveallsmsmail() error {
	if err := e.Login(); err != nil {
		return err
	}
	defer e.c.Logout()

	_, err := e.c.Select("INBOX", false)
	if err != nil {
		log.Println("#1 Moveallsmsmail Select error:", err)
		return err
	}

	criteria := imap.NewSearchCriteria()
	criteria.Text = []string{"sms", "SMS", "sms*", "SMS*"}

	ids, err := e.c.Search(criteria)
	if err != nil {
		log.Println("#2 Moveallsmsmail Search error:", err)
		return err
	}

	if len(ids) > 0 {
		seqset := new(imap.SeqSet)
		seqset.AddNum(ids...)
		err = e.c.Move(seqset, "INBOX.sms")
		if err != nil {
			log.Println("#3 Moveallsmsmail Move error:", err)
			return err
		}
	}

	return nil
}
