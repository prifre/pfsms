package ui

import (
	"fmt"
	"strings"
	"testing"
)

func TestAppendtotextfile(t *testing.T) {
	Appendtotextfile("emaillog.txt", "\r\nSome text1\r\n")
	Appendtotextfile("emaillog.txt", "\r\nSome text2\r\n")
	Appendtotextfile("emaillog.txt", "\r\nSome text3\r\n")
}
func TestGetAllCuntries(t *testing.T) {
	s := GetAllCountries()
	fmt.Println(s)
}
func TestFixphonenumber(t *testing.T) {
	s := Fixphonenumber("0736290839", "Sweden (+46)")
	if s != "0046736290839" {
		fmt.Println("Error #1")
		t.Fail()
	}
	s = Fixphonenumber("+46736290839", "Sweden (+46)")
	if s != "0046736290839" {
		fmt.Println("Error #2")
		t.Fail()
	}
	s = Fixphonenumber("+33736290839", "Sweden (+46)")
	if s != "0033736290839" {
		fmt.Println("Error #3 wrong country")
		t.Fail()
	}
	s = Fixphonenumber("+181736290839", "Sweden (+46)")
	if s != "00181736290839" {
		fmt.Println("Error #4 wrong country")
		t.Fail()
	}
	s = Fixphonenumber("0046736290839", "Sweden (+46)")
	if s != "0046736290839" {
		fmt.Println("Error #5")
		t.Fail()
	}
}
func TestGetLastLines(t *testing.T) {
	m := ""
	// for i := 0; i < 30; i++ {
	// 	m += fmt.Sprintf("LINE %d\r\n", i)
	// }
	fn := "pfsms.log"
	// f, _ := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// defer f.Close()
	// _, _ = f.WriteString(m)
	// f.Close()
	m = ReadLastLineWithSeek(fn, 25)
	r := strings.Split(m, "\n")
	for i := 1; i < len(r); i++ {
		fmt.Println(i, "=>", r[i])
	}
}
