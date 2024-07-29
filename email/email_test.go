package pfemail

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/BrianLeishman/go-imap"
)

func getEmailPasswords() (u string,p string) {
	var path string
	var err error
	var b0 []byte
	path, err = os.UserHomeDir()
	if err!=nil {
		panic("path")
	}
	path = fmt.Sprintf("%s%c%s",path ,os.PathSeparator, "pfsms")
	frfile:=path+string(os.PathSeparator)+"emailpasswords.txt"
    b0, err = os.ReadFile(frfile) // SQL to make tables!
    if err != nil {
        fmt.Print(err)
    }
	b:=string(b0)
	u=strings.Split(strings.Split(b,"\r\n")[0],"\t")[0]
	p=strings.Split(strings.Split(b,"\r\n")[0],"\t")[1]
	return u,p
}
func TestGetallemailmovetosms(t *testing.T) {
	e0 := new(Etype)
	imap.Verbose=true
	e0.uname,e0.pword = getEmailPasswords()
	e0.mserver = "mailcluster.loopia.se"
	e0.mport = 993
	e := e0.Getallmailmovetosmsfolder()
	if e!=nil {
		fmt.Println("UID:", e.UID)
		fmt.Println("SUBJECT:", e.Subject)
		fmt.Println("\r\n-------TEXT: ", strings.Replace(e.Text,"\r\n","",-1))
		fmt.Println("\r\n-------HTML: ", strings.Replace(e.HTML,"\r\n","",-1))
		fmt.Println("\r\nFlags: ", e.Flags)
	}
}
func TestCheckemaillogin(t *testing.T) {
	e0 := new(Etype)
	imap.Verbose=true
	e0.uname,e0.pword = getEmailPasswords()
	e0.mserver = "mailcluster.loopia.se"
	e0.mport = 993
	err := e0.Checkemaillogin()
    if err != nil {
        t.Fail()
    }

}