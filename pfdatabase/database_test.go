package pfdatabase

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

const databasepath = "..\\pfsmsdata\\pfsms.db"
func TestCheckGUI(t *testing.T) {
	checkGUI()
}

func TestDatabase(t *testing.T) {
	db:=new(DBtype)
	db.Databasepath =databasepath
	db.Opendb()
	db.Closedatabase()
}
func TestSetupDatabase(t *testing.T) {
	// Removes database and creates new one
	db:=new(DBtype)
	db.Databasepath = databasepath
	db.Closedatabase()
	os.Remove(db.Databasepath)
	db.Opendb()
	db.Setupdb()
	db.Closedatabase()
	fmt.Println("Database setup complete")
}

func TestImportdata(t *testing.T) {
	fmt.Println("Importing data into database")
	db:=new(DBtype)
	db.Databasepath = databasepath
	db.Opendb()
	// below removes database !!!!!
	// db.Closedatabase()
	// os.Remove(db.Databasepath)
	fmt.Println("Importing Customers")
	db.ImportCustomers(strings.ReplaceAll(db.Databasepath,"pfsms.db","customers.txt"))
	fmt.Println("Importing Groups")
	db.ImportGroups(strings.ReplaceAll(db.Databasepath,"pfsms.db","groups.txt"))
	fmt.Println("Importing History")
	db.ImportHistory(strings.ReplaceAll(db.Databasepath,"pfsms.db","history.txt"))
	fmt.Println("Importing Queue")
	db.ImportQueue(strings.ReplaceAll(db.Databasepath,"pfsms.db","queue.txt"))
	fmt.Println("Finished Importing data into database")
}
func TestExportdata(t *testing.T) {
	fmt.Println("Exporting data from database")
	db:=new(DBtype)
	db.Databasepath = databasepath
	db.Opendb()
	fmt.Println("Exporting Customers")
	db.ExportCustomers(strings.ReplaceAll(db.Databasepath,"pfsms.db","customers.txt"))
	fmt.Println("Exporting Groups")
	db.ExportGroups(strings.ReplaceAll(db.Databasepath,"pfsms.db","groups.txt"))
	fmt.Println("Exporting History")
	db.ExportHistory(strings.ReplaceAll(db.Databasepath,"pfsms.db","history.txt"))
	fmt.Println("Exporting Queue")
	db.ExportQueue(strings.ReplaceAll(db.Databasepath,"pfsms.db","que.txt"))
	fmt.Println("Finished Exporting data from database")
}

func TestShowCustomers(t *testing.T) {
	db:=new(DBtype)
	db.Databasepath = databasepath
	r:=db.ShowCustomers()
	fmt.Println(r)
}
func TestShowGroupnames(t *testing.T) {
	db:=new(DBtype)
	db.Databasepath = databasepath
	r:=db.ShowGroups()
	fmt.Println(r)
}
