package pfdatabase

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	_ "github.com/mattn/go-sqlite3"
)

type DBtype struct {
	conn         *sql.DB
	statement    *sql.Stmt
	reply        sql.Result
	Databasepath string
}
func (db *DBtype) Opendb() {
	var err error
	// var temp fyne.URI
	if db.conn != nil {
		return // allready opened!
	}
	if _, err := os.Stat(db.Databasepath); errors.Is(err, os.ErrNotExist) {
		db.Setupdb()
	}
	db.conn, err = sql.Open("sqlite3", db.Databasepath) // Open the created SQLite File
	if err != nil {
		log.Fatal("setupdatabase storage.Child error", err.Error())
		db.Setupdb()
		db.conn, err = sql.Open("sqlite3", db.Databasepath) // Open the created SQLite File
		if err != nil {
			log.Println("#1 Opendb ",err.Error())
		}
	}
}
func (db *DBtype) Setupdb() error {
	var err error
	db.Databasepath = fyne.CurrentApp().Preferences().String("pfsmsdb")
	if _, err = os.Stat(db.Databasepath); err != nil {
		log.Println("#1 Setupdb database not found, creating new db: " + db.Databasepath)
		var file *os.File
		file, err = os.Create(db.Databasepath) // Create SQLite file
		if err != nil {
			log.Println("#2 Setupdb Failed to create db", err.Error())
			return err
		}
		file.Close()
		err = db.Createtables() // Create Database Tables
		if err != nil {
			log.Println("#3 Setupdb Could not create tables!", err.Error())
			return err
		} else {
			log.Println("Database tables created")
		}
	}
	return err
}
func (db *DBtype) Closedatabase() error {
	var err error
	if db.conn!=nil {
		err = db.conn.Close()
		db.conn = nil
	} else {
		log.Println("#1 Closedatabase error. Closed already!")
	}
	return err
}
func (db *DBtype) Createtables() error {
	var err error
	// check if table exists
	db.conn, err = sql.Open("sqlite3", db.Databasepath) // Open the created SQLite File
	if err != nil {
		log.Fatal("#1 Createtables sql.Open ", err.Error())
	}
	_, table_check := db.conn.Query("SELECT * FROM tblCustomers;")
	if table_check == nil {
		return nil
		//table tblMain exists, so probably all is well...
	}
	//create tables...
	var s = "CREATE TABLE tblCustomers (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, "
   		s += "phone VARCHAR(20), firstname VARCHAR(100), lastname VARCHAR(100), note TEXT);"

		s += "CREATE TABLE tblGroups (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, "
   		s += "groupname VARCHAR(100), phone VARCHAR(100));"

		s += "CREATE TABLE tblHistory (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, "
   		s += "tstamp VARCHAR(20), groupname VARCHAR(100), phone VARCHAR(20), message TEXT);"

		s += "CREATE TABLE tblHashtable (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, "
   		s += "hash VARCHAR(100));"
	sq :=strings.Split(s,";")
	for i:=0;i<len(sq);i++ {
		if len(sq[i])<10 {
			continue
		}
		db.statement, err = db.conn.Prepare(sq[i]) // Prepare SQL Statement
		if err != nil {
			if err.Error() == "table tblCustomers already exists" {
				err = nil
				return err
			}
			log.Println("#1 CreateTables: '",sq[i],"'", err.Error())
		}
		db.reply, err = db.statement.Exec() // Execute SQL Statements
		if err != nil {
			log.Println("#2 CreateTables failed: ", sq[i], " ", err.Error(), db.reply)
			return err
		}
	}
	db.conn.Close()
	return err
}
func (db *DBtype) ShowCustomers()  [][]string {
	var phone,firstname,lastname string
	var data [][]string
	var err error
	db.Opendb()
	sq:="SELECT phone,firstname,lastname FROM tblCustomers ORDER BY phone ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ShowCustomers Query error:", err.Error())
		return nil
	}
	for rows.Next() {
		err = rows.Scan(&phone,&firstname,&lastname)
		if err != nil {
			log.Println("#2 ShowCustomers Scan error:", err.Error())
			return nil
		}
		data = append(data,[]string{phone,firstname,lastname})
	}
	db.Closedatabase()
	return data
}
func (db *DBtype) ImportCustomers(frfile string) {
	// customer textfile should be phone <<tab>> firstname <<tab>> lastname <<tab>> note <<cr>> <<lf>>
	var err error
	var b0 []byte
	var c int
	var sq,phone,firstname,lastname,note string
    b0, err = os.ReadFile(frfile) // SQL to make tables!
    if err != nil {
        log.Println("#1 ImportCustomers ReadFile", err.Error())
		return
    }
	b:=string(b0)
	//fix bad characters
	b = Removebadsqlcharacters(b)
	b = strings.Replace(b,"\n","",-1)
	// allcurrent := db.ShowCustomers()
	db.Opendb()
	for i:=0;i<len(strings.Split(b,"\r"));i++ {
		b1:=strings.Split(b,"\r")[i]
		b2:=strings.Split(b1,"\t")
		firstname=""
		lastname=""
		note=""
		phone=""
		if len(b2)>0 {
			phone = Fixphonenumber(b2[0])
		} else {
			continue
		}
		if len(b2)>1 {
			firstname = b2[1]
		}
		if len(b2)>2 {
			lastname = b2[2]
		}
		if len(b2)>3 {
			note = b2[3]
		}
		if len(phone)<5 {
			continue
		}
		// check if phonenumber already exists
		var cnt string
		r, err := db.conn.Query("SELECT COUNT(*) AS cnt FROM tblCustomers WHERE phone = '"+phone+"'")
		if err!=nil {
			log.Println("#2 ImportCustomers Query failed: ", sq, " ", err.Error())
		} else {
			for r.Next() {
				r.Scan(&cnt)
			}
		}
		if cnt>"0" {
			continue
		}
		sq ="INSERT INTO tblCustomers (phone,firstname,lastname,note)"
		sq += fmt.Sprintf(" VALUES ('%s','%s','%s','%s')",phone,firstname,lastname,note)
		db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
		if err != nil {
			log.Println("#3 ImportCustomers prepare failed: ", sq, " ", err.Error())
			return
		}
		db.reply, err = db.statement.Exec() // Execute SQL Statements
		if err != nil {
			log.Println("#4 ImportCustomers Exec failed: ", sq, " ", err.Error())
			return
		}
		c++
	}
	log.Printf("Imported %d customers.\r\n",c)
	db.Closedatabase()
}
func (db *DBtype) ExportCustomers(tofile string) {
	var err error
	var sq,txt,phone,firstname,lastname,note string
	db.Opendb()
	sq ="SELECT phone,firstname,lastname,note FROM tblCustomers ORDER BY phone ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ExportCustomers Query ", err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&phone,&firstname,&lastname,&note)
		if err!=nil {
			log.Println("#2 ExpotCustomers Scan ",err.Error())
		}
		if len(phone)>0 {
			txt +=fmt.Sprintf("%s\t%s\t%s\t%s\r\n",phone,firstname,lastname,note)
		}
	}
	db.Closedatabase()
	if txt=="" {
		// export sample data!!!
		txt ="+46736290839\tPeter\tFreund\r\n"
		txt +="087543169\tLin\tZhang\r\n"
		txt +="0890510\tFröken\tUr\tClock\r\n"
		txt +="004690200\tTelia Support\tUtomlands\r\n"
		txt +="0046771735311\tTre Support\tCompany\r\n"
		txt +="+46708 222 222\tTelenor Support\tCompany\r\n"
		txt +="+12024561111\tWhite\tHouse\tin USA\r\n"
	}
	err = os.WriteFile(tofile, []byte(txt), 0644)
	if err != nil {
		log.Println("#3 ExportCustomers WriteFile ",err.Error())
	}
}
func (db *DBtype) ShowGroups() []string {
	// should show all available groupnames
	var data []string
	var err error
	var groupname string
	db.Opendb()
	sq:="SELECT DISTINCT groupname FROM tblGroups ORDER BY groupname ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ShowGroups Query error:", err.Error())
		return nil
	}
	for rows.Next() {
		err = rows.Scan(&groupname)
		if err != nil {
			log.Println("#2 ShowGroups Scan error:", err.Error())
			return nil
		}
		data=append(data,groupname)
	}
	db.Closedatabase()
	return data
}
func (db *DBtype) ShowAllGroups() [][]string {
	// should show all available groupnames
	var data [][]string
	var err error
	var groupname,phone string
	db.Opendb()
	sq:="SELECT groupname,phone FROM tblGroups ORDER BY groupname ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ShowGropAllnames Query error:", err.Error())
		return nil
	}
	for rows.Next() {
		err = rows.Scan(&groupname,&phone)
		if err!=nil {
			log.Println("#2 ShowAllGroups Scan ", err.Error())
		}
		data=append(data,[]string{groupname,phone})
	}
	db.Closedatabase()
	return data
}
func (db *DBtype) SaveGroup(groupname string,phones string) {
	var err error
	var sq string
	groupname=Removebadsqlcharacters(groupname)
	p:=strings.Split(phones,",")
	for i:=0;i<len(p);i++ {
		p[i]=Fixphonenumber(p[i])
	}
	db.Opendb()
	sq = fmt.Sprintf("DELETE FROM tblGroups WHERE groupname = '%s'",groupname)
	_,err = db.conn.Exec(sq)
	if err!=nil {
		log.Println("#1 SaveGroup DELETE failed ", err.Error())
		return
	}
	for i:=0;i<len(p);i++ {
		sq =fmt.Sprintf("INSERT INTO tblGroups (groupname,phone) VALUES ('%s','%s')",groupname,p[i])
		db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
		if err != nil {
			log.Println("#2 SaveGroup prepare", err.Error())
			return 
		}
		db.reply, err = db.statement.Exec() // Execute SQL Statements
		if err != nil {
			log.Println("#3 SaveGroups Exec", err.Error())
			return
		}
	}
	db.Closedatabase()
}
func (db *DBtype) DeleteGroup(g string) {
	var err error
	var sq string
	db.Opendb()
	sq = fmt.Sprintf("DELETE FROM tblGroups WHERE groupname = '%s'",g)
	db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
	if err != nil {
		log.Println("#1 DeleteGroup Prepare failed: ", sq, " ", err.Error())
		return
	}
	db.reply, err = db.statement.Exec() // Execute SQL Statements
	if err != nil {
		log.Println("#2 DeleteGroup Exec failed: ", sq, " ", err.Error())
		return
	}
	db.Closedatabase()
}
func (db *DBtype) ImportGroups(b string) {
	// b should be tab-separated \r separated string with -> grouname <tab> phone <CR>
	// since import is done directory from field in Messages... (not from file)
	var sq string
	var err error
	db.Opendb()
	sq="DELETE FROM tblGroups;"
	db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
	if err != nil {
		log.Println("#1 ImportGroups Prepare DELETE", err.Error())
		return
	}
	_, err = db.statement.Exec() // Execute SQL Statements
	if err!=nil {
		log.Println("#2 ImportGroups Failed DELETE tblGroups")
	}
	b=Removebadsqlcharacters(b)
	for i:=0;i<len(strings.Split(b,"\r"));i++ {
		b1:=strings.Split(b,"\r")[i]
		b2:=strings.Split(b1,"\t")
		if len(b2)<2 {
			continue
		}
		b2[1]=Fixphonenumber(b2[1])
		if len(b2[1])<5 {
			continue
		}
		sq ="INSERT INTO tblGroups (groupname,phone)"
		sq += fmt.Sprintf(" VALUES ('%s','%s')",b2[0],b2[1])
		db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
		if err != nil {
			log.Println("#3 ImportGroups prepare", err.Error())
			return
		}
		db.reply, err = db.statement.Exec() // Execute SQL Statements
		if err != nil {
			log.Println("#4 ImportGroups Exec ", err.Error())
			return
		}
	}
	db.Closedatabase()
}
func (db *DBtype) ExportGroups(tofile string) {
	var err error
	var sq,txt,groupname,phone string
	db.Opendb()
	sq ="SELECT groupname,phone FROM tblGroups ORDER BY groupname ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ExportGroups Query ", err.Error())
		return
	}
	for rows.Next() {
		err = rows.Scan(&groupname,&phone)
		if err!=nil {
			log.Println("#2 ExpotGroups Scan ,",err.Error())
		}
		if len(phone)>0 {
			txt +=fmt.Sprintf("%s\t%s\r\n",groupname,phone)
		}
	}
	db.Closedatabase()
	if txt=="" {
		// export sample data!!!
		txt ="Sample\t0046736290839\r\n"
		txt +="Sample\t004687543169\r\n"
		txt +="Sample\t04690510\r\n"
		txt +="Sample\t0012024561111\r\n"
	}
	err = os.WriteFile(tofile, []byte(txt), 0644)
	if err != nil {
		log.Println("#3 ExportGroups WriteFile,",err.Error())
	}
}
func (db *DBtype) ExportHistory(tofile string) {
	// s += "CREATE TABLE tblHistory (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, "
	// s += "tstamp VARCHAR(20), phone VARCHAR(20), groupname VARCHAR(100), message TEXT);"
	var err error
	var sq,txt,tstamp,groupname,phone,message string
	db.Opendb()
	sq ="SELECT tstamp,groupname,phone,message FROM tblHistory ORDER BY tstamp ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ExportHistory Query ", err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&tstamp,&groupname,&phone,&message)
		if err!=nil {
			log.Println("#2 ExportHistory Scan ",err.Error())
		}
		txt +=fmt.Sprintf("%s\t%s\t%s\t%s\r\n",tstamp,groupname,phone,message)
	}
	db.Closedatabase()
	if txt=="" {
		// export sample History!!!
		txt ="20240801135055\ttest\t0046736290839\tThis is a test message 1\r\n"
		txt +="20240801135055\ttest\t0046736290839\tThis is a test message 2\r\n"
		txt +="20240801135055\ttest\t0046736290839\tThis is a test message 3\r\n"
		txt +="20240801135055\ttest\t0046736290839\tThis is a test message 4\r\n"
	}
	err = os.WriteFile(tofile, []byte(txt), 0644)
	if err != nil {
		log.Println("#3 ExportHistory WriteFile ",err.Error())
	}
}
func (db *DBtype) SaveHistory(result [][]string) {
	// resulting string with history from ariassms = tstamp,phone,message
	// message in \"\"
	var sq string
	var err error
	db.Opendb()
	for _,r := range result {
		// tstamp := r[0]
		// phone := r[1]
		// groupname := r[2]
		// message := r[3]
		sq ="INSERT INTO tblHistory (tstamp,groupname,phone,message)"
		sq += fmt.Sprintf(" VALUES ('%s','%s','%s','%s')",r[0],r[1],r[2],r[3])
		db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
		if err != nil {
			log.Println("#1 SaveHistory Prepare", err.Error())
			return
		}
		db.reply, err = db.statement.Exec() // Execute SQL Statements
		if err != nil {
			log.Println("#2 SaveHistory Exec ", err.Error())
			return
		}
	}
	db.Closedatabase()
}
func (db *DBtype) GetFname(phone string) string {
	var firstname string
	db.Opendb()
	rows, err := db.conn.Query("SELECT firstname FROM tblCustomers WHERE phone = '"+phone+"'")
	if err != nil {
		fmt.Println("#2 ShowCustomers Query error:", err.Error())
		return ""
	}
	for rows.Next() {
		err = rows.Scan(&firstname)
		if err!=nil {
			log.Println("rows.Scan failed in GetFname")
		}
	}
	db.Closedatabase()
	return firstname
}
func (db *DBtype) GetLname(phone string) string {
	var lastname string
	db.Opendb()
	rows, err := db.conn.Query("SELECT lastname FROM tblCustomers WHERE phone = '"+phone+"'")
	if err!=nil {
		return ""
	}
	for rows.Next() {
		if err != nil {
			fmt.Println("#2 ShowCustomers Query error:", err.Error())
			return ""
		}
		err = rows.Scan(&lastname)
		if err!=nil {
			log.Println("rows.Scan failed in GetFname")
		}
	}
	db.Closedatabase()
	return lastname
}
func Fixphonenumber(pn string) string {
	// pn phonenumber  cc coutrycode
	// Sweden (+46) converts to 0046
	cc := fyne.CurrentApp().Preferences().StringWithFallback("mobileCountry", "Sweden(+46)")
	var cci string ="00"
	for i:=0;i<len(cc);i++ {
		if strings.Index("0123456789",string(cc[i]))>0 {
			cci +=string(cc[i])
		}
	}
	var cc2,pn2 string 
	for i:=0;i<len(cc);i++ {
		if strings.Index("0123456789",string(cc[i]))>0 {
			cc2 +=string(cc[i])
		}
	}
	for i:=0;i<len(pn);i++ {
		if strings.Index("+0123456789",string(pn[i]))>0 {
			pn2 +=string(pn[i])
		}
	}
	if len(pn2)<6 {
		return ""
	}
	if pn2[0:2]==string("00") {
		return pn2
	}
	if string(pn2[0])=="0" {
		return  "00"+cc2+pn2[1:]
	}
	if string(pn2[0])=="+" {
		return "00"+pn2[1:]
	}
	return "00"+cc2+pn2
}
func Removebadsqlcharacters(v string) string {
	v=strings.Replace(v,"´","",-1)
	v=strings.Replace(v,"`","",-1)
	v=strings.Replace(v,"@","",-1)
	v=strings.Replace(v,"'","",-1)
	v=strings.Replace(v,"\"","",-1)
	v=strings.Replace(v,"_","",-1)
	v=strings.Replace(v,"%","",-1)
	v=strings.Replace(v,"#","",-1)
	v=strings.Replace(v,"/","",-1)
	return v
}