package pfdatabase

import (
	"fmt"
	"testing"
)

func TestCheckPasswords(t *testing.T) {
	var r,result string
	var starthash string = "23546546346"
	var h string
	var txt string="Detta är en teststräng."
	var err error
	h, err = MakeHash()
	if err!=nil {
		t.Fatalf("SetHash failed %s",err.Error())
	}
	if h !=starthash {
		fmt.Println("Old hash exists: ",h)
	}
	r,err =EncryptPassword(txt,h)
	if err!=nil {
		t.Fatalf("Encrypt failed %s",err.Error())
	}
	fmt.Println("Original text:",txt)
	fmt.Println("Encryptedtext:",r)
	result,err = DecryptPassword(r,h)
	if result != txt {
		t.Fatalf("Decrypt failed %s",err.Error())
	}
	fmt.Println("Decryptedtext:",result)
}
