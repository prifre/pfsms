package db

import (
	"fmt"
	"testing"
)

func TestDatabase(t *testing.T) {
	d:=new(dbtype)
	d.Setupdb()
	d.Opendb()
}

func TestImportdata(t *testing.T) {
	d:=new(dbtype)
	d.Setupdb()
	d.ImportCustomers("KUNDER.txt")
	d.AddMessage("testmeddelande","Detta är ett speciellt innehåll")
}

func TestShowCustomers(t *testing.T) {
	d:=new(dbtype)
	d.Setupdb()
	r,err:=d.ShowCustomers(10,30)
	if err!=nil {
		t.Fatalf("ShowCustomer failed %s",err.Error())
	}
	fmt.Println(r)
}
