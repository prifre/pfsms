package pfemail

import (
	"fmt"
	"strconv"

	"github.com/BrianLeishman/go-imap"
)

// func main() {
// 	var username, password, mailserver string
// 	var mailserverport int
// 	password = "1Fv2zjofj6Jb!"
// 	username = "sms@sollentunaram.se"
// 	mailserver = "mailcluster.loopia.se"
// 	mailserverport = 993
// 	e := Getonemail(username, password, mailserver, mailserverport)
// 	fmt.Println("UID:", e.UID)
// 	fmt.Println("SUBJECT:", e.Subject)
// 	fmt.Println("\r\n-------TEXT:\r\n", e.Text)
// 	fmt.Println("\r\n-------HTML:\r\n", e.HTML)
// 	fmt.Println("\r\nFlags:\r\n", e.Flags)
// 	fmt.Println(e.Flags)
// }
type Etype struct {
	uname 		string
	pword		string
	mserver 	string
	mport 		int
	// mfrequency	int
}
func (e *Etype) SetupEmail(s1,s2,s3,s4 string) {
	e.mserver=s1
	e.uname=s2
	e.pword=s3
	s0,_:=strconv.Atoi(s4)
	e.mport = int(s0)
}
func (e *Etype) Getonemail() *imap.Email {
	// get emails in inbox
	// Defaults to false. This package level option turns on or off debugging output, essentially.
	// If verbose is set to true, then every command, and every response, is printed,
	// along with other things like error messages (before the retry limit is reached)
	imap.Verbose = true

	// Defaults to 10. Certain functions retry; like the login function, and the new connection function.
	// If a retried function fails, the connection will be closed, then the program sleeps for an increasing amount of time,
	// creates a new connection instance internally, selects the same folder, and retries the failed command(s).
	// You can check out github.com/StirlingMarketingGroup/go-retry for the retry implementation being used
	// Create a new instance of the IMAP connection you want to use
	imap.RetryCount = 3
	im, err := imap.New(e.uname, e.pword, e.mserver, e.mport)
	if err != nil {
		fmt.Println("#0 Instance of IMAP could not be created."+err.Error())
		return nil
	}		
	defer im.Close()

	// Folders now contains a string slice of all the folder names on the connection
	folders, err := im.GetFolders()
	if err != nil {
		fmt.Println("#1 No email folders found."+err.Error())
		return nil
	}		
	smsfolderexists := false
	for _, f := range folders {
		if f == "INBOX.sms" {
			smsfolderexists = true
		}
	}
	if !smsfolderexists {
		cmd := "Create INBOX.sms"
		fmt.Println(cmd)
		im.Exec(cmd, true, 1, nil)
	}
	// folders = []string{
	// 	"INBOX",
	// 	"INBOX/My Folder"
	// 	"Sent Items",
	// 	"Deleted",
	// }

	// Now we can loop through those folders
	for _, f := range folders {

		// And select each folder, one at a time.
		// Whichever folder is selected last, is the current active folder.
		// All following commands will be executing inside of this folder
		if f == "INBOX" {
			err = im.SelectFolder(f)
			if err != nil {
				fmt.Println("#3 No INBOX folder found."+err.Error())
				return nil
			}		
		
			// This function implements the IMAP UID search, returning a slice of ints
			// Sending "ALL" runs the command "UID SEARCH ALL"
			// You can enter things like "*:1" to get the first UID, or "999999999:*"
			// to get the last (unless you actually have more than that many emails)
			// You can check out https://tools.ietf.org/html/rfc3501#section-6.4.4 for more
			uids, err := im.GetUIDs("ALL")
			if err != nil {
				fmt.Println("#2 No emails found."+err.Error())
				return nil
			}		
		
			// uids = []int{1, 2, 3}

			// GetEmails takes a list of ints as UIDs, and returns new Email objects.
			// If an email for a given UID cannot be found, there's an error parsing its body,
			// or the email addresses are malformed (like, missing parts of the address), then it is skipped
			// If an email is found, then an imap.Email struct slice is returned with the information from the email.
			// The Email struct looks like this:
			// type Email struct {
			// 	Flags     []string
			// 	Received  time.Time
			// 	Sent      time.Time
			// 	Size      uint64
			// 	Subject   string
			// 	UID       int
			// 	MessageID string
			// 	From      EmailAddresses
			// 	To        EmailAddresses
			// 	ReplyTo   EmailAddresses
			// 	CC        EmailAddresses
			// 	BCC       EmailAddresses
			// 	Text      string
			// 	HTML      string
			//	Attachments []Attachment
			// }
			// Where the address type fields are maps like [EmailAddress:Name EmailAddress2:Name2]
			// and an Attachment is a struct containing the Name, Content, and the MimeType (both as strings)
			emails, err := im.GetEmails(uids...)
			if err != nil {
				fmt.Println("#3 No emails found."+err.Error())
				return nil
			}		
			if len(emails) > 0 {
				for _, e := range emails {
//					if e.Subject == "sms" {
						// Should print a summary of one of the the emails
						// (yes, I said "one of", don't expect the emails to be returned in any particular order)
						//				fmt.Print(emails[0])
						im.MoveEmail(e.UID, "INBOX.sms")
						return e
//					}
				}
				//			im.MoveEmail(emails[0].UID, "INBOX/My Folder")
				// Subject: FW: FW:  FW:  New Order
				// To: Brian Leishman <brian@stumpyinc.com>
				// From: Customer Service <sales@totallylegitdomain.com>
				// Text: Hello, World!...(4.3 kB)
				// HTML: <html xmlns:v="urn:s... (35 kB)
				// 1 Attachment(s): [20180330174029.jpg (192 kB)]
			}
		}
	}
	return nil
}
