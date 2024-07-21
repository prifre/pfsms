package pfemail

import (
	"fmt"
	"testing"
)

func TestGetonemail(t *testing.T) {
	var username, password, mailserver string
	var mailserverport int
	password = "1Fv2zjofj6Jb!"
	username = "sms@sollentunaram.se"
	mailserver = "mailcluster.loopia.se"
	mailserverport = 993
	e := Getonemail(username, password, mailserver, mailserverport)
	if e!=nil {
		fmt.Println("UID:", e.UID)
		fmt.Println("SUBJECT:", e.Subject)
		fmt.Println("\r\n-------TEXT:\r\n", e.Text)
		fmt.Println("\r\n-------HTML:\r\n", e.HTML)
		fmt.Println("\r\nFlags:\r\n", e.Flags)
	}
}