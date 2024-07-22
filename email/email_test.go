package pfemail

import (
	"fmt"
	"testing"
)

func TestGetonemail(t *testing.T) {
	var e0 Etype
	e0.pword = ""
	e0.uname = "sms@sollentunaram.se"
	e0.mserver = "mailcluster.loopia.se"
	e0.mport = 993
	e := e0.Getonemail()
	if e!=nil {
		fmt.Println("UID:", e.UID)
		fmt.Println("SUBJECT:", e.Subject)
		fmt.Println("\r\n-------TEXT:\r\n", e.Text)
		fmt.Println("\r\n-------HTML:\r\n", e.HTML)
		fmt.Println("\r\nFlags:\r\n", e.Flags)
	}
}