package ui

import (
	"fmt"
	"strings"
	"testing"
)

func TestAppendtotextfile(t *testing.T) {
	Appendtotextfile("emaillog.txt","\r\nSome text1\r\n")
	Appendtotextfile("emaillog.txt","\r\nSome text2\r\n")
	Appendtotextfile("emaillog.txt","\r\nSome text3\r\n")
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

func TestGetLastLines(t *testing.T) {
	var m string
	m,_ = ReadLastLineWithSeek("smslog.txt",5)
	r:=strings.Split(m,"\n")
	for i:=0;i<len(r);i++ {
		fmt.Println(i,"=>",r[i])
	}
}
