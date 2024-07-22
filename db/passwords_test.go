package db

import (
	"testing"
)

func TestCheckPasswords(t *testing.T) {
	var r,result string
	var starthash string = "aioprgjöaisdjgöajsdg"
	var hash string
	var txt string="Detta är en teststräng."
	var err error
	d:=new(DBtype)
	err = d.SetHash(starthash)
	if err!=nil {
		t.Fatalf("SetHash failed %s",err.Error())
	}
	hash,err = d.GetHash()
	if err!=nil {
		t.Fatalf("GetHash failed %s",err.Error())
	}
	if hash !=starthash {
		t.Fatalf("GetHash/SetHash compare failed %s",err.Error())
	}
	if d.hash !=starthash {
		t.Fatalf("GetHash/SetHash compare failed %s",err.Error())
	}
	r,err =d.EncryptPassword(txt)
	if err!=nil {
		t.Fatalf("Encrypt failed %s",err.Error())
	}
	result,err = d.DecryptPassword(r)
	if result != txt {
		t.Fatalf("Decrypt failed %s",err.Error())
	}
}
