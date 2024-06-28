package database

import (
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
