package db

import (
	"fmt"
	"testing"
)

func TestDatabase(t *testing.T) {
	d:=new(DBtype)
	d.Opendb()
	d.Setupdb()
}

func TestImportdata(t *testing.T) {
	d:=new(DBtype)
	d.Setupdb()
	d.ImportCustomers("KUNDER.txt")
	d.AddMessage("testmeddelande","Detta är ett speciellt innehåll")
}

func TestShowCustomers(t *testing.T) {
	d:=new(DBtype)
	d.Setupdb()
	r,err:=d.ShowCustomers(10,30)
	if err!=nil {
		t.Fatalf("ShowCustomer failed %s",err.Error())
	}
	fmt.Println(r)
}
