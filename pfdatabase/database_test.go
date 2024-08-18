package pfdatabase

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
	d.Closedatabase()
}

func TestImportdata(t *testing.T) {
	db:=new(DBtype)
	// below removes database !!!!!
	db.Opendb()
	db.Closedatabase()
	os.Remove(db.Databasepath)
	db.ImportCustomers(strings.Replace(db.Databasepath,"pfsms.db","KUNDER.txt",-1))
}
func TestShowCustomers(t *testing.T) {
	r:=new(DBtype).ShowCustomers()
	fmt.Println(r)
}
func TestShowGroupnames(t *testing.T) {
	r:=new(DBtype).ShowGroups()
	fmt.Println(r)
}
