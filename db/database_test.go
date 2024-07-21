package db

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestDatabase(t *testing.T) {
	d:=new(DBtype)
	d.Opendb()
	d.Setupdb()
}

func TestImportdata(t *testing.T) {
	var err error
	db:=new(DBtype)
	// below removes database !!!!!
	db.Opendb()
	db.Closedatabase()
	os.Remove(db.Databasepath)
	err = db.ImportCustomers(strings.Replace(db.Databasepath,"pfsms.db","KUNDER2.txt",-1))
	if err!=nil {
		t.Fatalf("ImportCustomers failed %s",err.Error())
	}
	err =db.AddMessage("testmeddelande","Detta är ett speciellt innehåll")
	if err!=nil {
		t.Fatalf("AddMessage failed %s",err.Error())
	}		
}
func TestShowCustomers(t *testing.T) {
	d:=new(DBtype)
	r,err:=d.ShowCustomers(10,30)
	if err!=nil {
		t.Fatalf("ShowCustomer failed %s",err.Error())
	}
	fmt.Println(r)
}
func TestShowGroupnames(t *testing.T) {
	d:=new(DBtype)
	r,err:=d.ShowGroupnames()
	if err!=nil {
		t.Fatalf("ShowCustomer failed %s",err.Error())
	}
	fmt.Println(r)
}
func TestAddMessage(t *testing.T) {
	d:=new(DBtype)
	d.AddMessage("Specialmessage","Hej hopp du glade åäö")
}