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
	var r int
	db:=new(DBtype)
	// below removes database !!!!!
	db.Opendb()
	db.Closedatabase()
	os.Remove(db.Databasepath)
	err = db.ImportCustomers(strings.Replace(db.Databasepath,"pfsms.db","KUNDER.txt",-1))
	if err!=nil {
		t.Fatalf("ImportCustomers failed %s",err.Error())
	}
	r,err =db.AddMessage("testmeddelande","Detta är ett speciellt innehåll")
	if err!=nil {
		t.Fatalf("AddMessage failed %s",err.Error())
	}
	if r<0 {
		t.Fatalf("AddMessage failed 2 %s",err.Error())		
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
func TestMessageroutines(t *testing.T) {
	d:=new(DBtype)
	r1,_:=d.ShowMessages()
	fmt.Println(r1)
	r2,_:=d.AddMessage("Specialmessage","Hej hopp du glade åäö")
	err:=d.DeleteMessage(r2)
	if err!=nil {
		t.Fatalf("DeleteMessage failed %s",err.Error())
	} else {
		fmt.Println("DeleteMessage ok")
	}
	d.AddMessage("Akvarellkurs","Välkommen till akvarellkurs i höst\r\nVi startar den 5:e Augusti klockan 18.00.\r\nPeter Sollentuna Ram")
	r3,_:=d.AddMessage("Galleriinbjudan","Tips att Edsvik har vernissage 12/8 kl. 13 med konstnären Lin!")
	d.UpdateMessage(r3,"Galleriinbjudan","Litet tips att Edsvik har vernissage 12/8 kl. 13 med konstnären Lin!")
	r1,_=d.ShowMessages()
	fmt.Println(r1)
}