package ui

import (
	"fmt"
	"testing"
)

func TestAppendtotextfile(t *testing.T) {
	Appendtotextfile("emaillog.txt","\r\nSome text1")
	Appendtotextfile("emaillog.txt","\r\nSome text2")
	Appendtotextfile("emaillog.txt","\r\nSome text3\r\nend")
}
func TestGetAllCuntries(t *testing.T) {
	s:=GetAllCountries()
	fmt.Println(s)
}
func TestFixphonenumber(t *testing.T) {
	s:=Fixphonenumber("0736290839","Sweden (+46)")
	if s!="0046736290839" {
		fmt.Println("Error #1")
		t.Fail()
	}
	s=Fixphonenumber("+46736290839","Sweden (+46)")
	if s!="0046736290839" {
		fmt.Println("Error #2")
		t.Fail()
	}
	s=Fixphonenumber("0046736290839","Sweden (+46)")
	if s!="0046736290839" {
		fmt.Println("Error #3")
		t.Fail()
	}
}
