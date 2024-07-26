package db

import (
	"fmt"
	"testing"
)

func TestCheckPasswords(t *testing.T) {
	var r,result string
	var starthash string = "23546546346"
	var txt string="Detta är en teststräng."
	var err error
	db:=new(DBtype)
	db.Opendb()
	err = db.MakeHash()
	if err!=nil {
		t.Fatalf("SetHash failed %s",err.Error())
	}
	if db.hash !=starthash {
		fmt.Println("Old hash exists: ",db.hash)
	}
	r,err =db.EncryptPassword(txt)
	fmt.Println(r)
	if err!=nil {
		t.Fatalf("Encrypt failed %s",err.Error())
	}
	result,err = db.DecryptPassword(r)
	if result != txt {
		t.Fatalf("Decrypt failed %s",err.Error())
	}
}
