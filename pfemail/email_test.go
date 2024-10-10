package pfemail

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/BrianLeishman/go-imap"
)

func getEmailPasswords() (u string, p string) {
	var path string
	var err error
	var b0 []byte
	path, err = os.UserHomeDir()
	if err != nil {
		panic("path")
	}
	path = fmt.Sprintf("%s%c%s%c%s", path, os.PathSeparator, "pfsms", os.PathSeparator, "emailpasswords.txt")
	b0, err = os.ReadFile(path) // SQL to make tables!
	if err != nil {
		fmt.Print(err)
	}
	b := string(b0)
	u = strings.Split(strings.Split(b, "\r\n")[0], "\t")[0]
	p = strings.Split(strings.Split(b, "\r\n")[0], "\t")[1]
	return u, p
}
func TestGetallemailmovetosms(t *testing.T) {
	e0 := new(Etype)
	imap.Verbose = true
	e0.uname, e0.pword = getEmailPasswords()
	e0.mserver = "mailcluster.loopia.se"
	e0.mport = "993"
	e0.Checkemaillogin()
	e := e0.Getallsmsmail()
	e0.Moveallsmsmail()
	if e != nil {
		fmt.Println("UID:", e[0].Envelope)
		fmt.Println("SUBJECT:", e[0].Envelope.Subject)
		fmt.Println("SENDER: ", e[0].Envelope.Sender)
		fmt.Println("\r\nFlags: ", e[0].Flags)
	}
}
func TestCheckemaillogin(t *testing.T) {
	e0 := new(Etype)
	imap.Verbose = true
	e0.uname, e0.pword = getEmailPasswords()
	e0.mserver = "mailcluster.loopia.se"
	e0.mport = "993"
	err := e0.Checkemaillogin()
	if err != nil {
		t.Fail()
	}

}
