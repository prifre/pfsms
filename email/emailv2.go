package pfemail

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/prifre/pfsms/db"
)

type Etype struct {
	uname 		string
	pword		string
	mserver 	string
	mport 		string
	c 			*client.Client
		// mfrequency	int
}
func (e *Etype) Connect(s1,s2,s3,s4 string) {
	e.mserver=s1
	e.uname	= s2
	e.pword	= s3
	e.mport = s4
}
func (e *Etype) Checkemaillogin() error {
	e.Login()
	if e.c==nil {
		return fmt.Errorf("Login failed")
	}
	e.c.Logout()
	return nil
}
func (e *Etype) Login() error {
	var hash string
	var err error
	if e.pword=="" {
		hash, err =db.MakeHash()
		if err!=nil {
			log.Println("buildLog MakeHash error ",err.Error())
		}
		passwdstring:=fyne.CurrentApp().Preferences().StringWithFallback("ePassword","")
		e.pword,err=db.DecryptPassword(passwdstring,hash)
		if err!=nil {
			log.Println("Decryptpassword error ",err.Error())
		}
		e.mserver=fyne.CurrentApp().Preferences().StringWithFallback("eServer","")
		e.uname = fyne.CurrentApp().Preferences().StringWithFallback("eUser","")
		e.mport = fyne.CurrentApp().Preferences().StringWithFallback("ePort","")
	}
	log.Println("Connecting to server...")
	// Connect to server
	e.c, err = client.DialTLS(fmt.Sprintf("%s:%s",e.mserver,e.mport), nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Println("Connected")
	
	// Login
	err = e.c.Login(e.uname,e.pword)
	if err!=nil {
		log.Fatal(err)
		return nil
	}
	log.Println("Logged in")
	return err
}

func (e *Etype) ListMailboxes () []string {
	var mb []string
	// List mailboxes
	e.Login()
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- e.c.List("", "*", mailboxes)
	}()
	
	log.Println("Mailboxes:")
	for m := range mailboxes {
		mb = append(mb,string(m.Name))
	}
	
	if err := <-done; err != nil {
		log.Fatal(err)
	}
	done=nil
	e.c.Logout()
	return mb	
}
func (e *Etype) Createmailboxfolder(f string) error {
	var err error
	e.Login()
	err = e.c.Create("INBOX."+f)
	if err!=nil {
		fmt.Println("#1 Create ",err.Error())
	}
	err = e.c.Expunge(nil)
	if err!=nil {
		fmt.Println("#2 Expunge ",err.Error())
	}
	e.c.Logout()
return err
}
func (e *Etype) Getallsmsmail() []*imap.Message {
	// Select INBOX
	var err error
	err = e.Login()
	if err!=nil {
		log.Println("ERROR ",err.Error())
	}
	fmt.Println(e.c)
	var imsgs []*imap.Message
	_, err = e.c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	criteria := imap.NewSearchCriteria()
	criteria.Text = []string{"TEST 123"}
	ids, err := e.c.Search(criteria)
	if err != nil {
		log.Fatal(err)
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
			imsgs = append(imsgs,msg)
			log.Println("* " + msg.Envelope.Subject)
		}
		if err := <-done; err != nil {
			log.Fatal(err)
		}
		e.c.Move(seqset,"INBOX.sms")
	}
	e.c.Logout()
	return imsgs
}

func (e *Etype) Moveallsmsmail() error {
	// Select INBOX
	e.Login()
	var err error
	_, err = e.c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	criteria := imap.NewSearchCriteria()
	criteria.Text = []string{"sms","SMS","sms*","SMS*"}
	var ids []uint32

	ids, err = e.c.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}
	if len(ids) > 0 {
		seqset := new(imap.SeqSet)
		seqset.AddNum(ids...)
		err = e.c.Move(seqset,"INBOX.sms")
	}		
	e.c.Logout()
	return err
}