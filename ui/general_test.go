package ui

import (
	"testing"
)

func TestAppendtotextfile(t *testing.T) {
	Appendtotextfile("emaillog.txt","\r\nSome text1")
	Appendtotextfile("emaillog.txt","\r\nSome text2")
	Appendtotextfile("emaillog.txt","\r\nSome text3\r\nend")
}