package pfdatabase

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	_ "github.com/mattn/go-sqlite3"
)

type DBtype struct {
	conn         *sql.DB
	statement    *sql.Stmt
	reply        sql.Result
	Databasepath string
}
func checkGUI() bool {
	b := os.Getenv("APP_MODE") == "GUI"
	return b
}
func (db *DBtype) Opendb() error {
	var err error
	// var temp fyne.URI
	if db.conn != nil {
		return nil // allready opened!
	}
	_,err=os.Stat(db.Databasepath)
	if errors.Is(err, os.ErrNotExist) {
		// err = errors.New("Database not found, creating new db: " + db.Databasepath)
		err = db.Setupdb()
		if err != nil {
			return errors.New("#1 Opendb SetupDB " + err.Error())
		}
	}
	db.conn, err = sql.Open("sqlite3", db.Databasepath) // Open the created SQLite File
	if err != nil {
		err =db.Setupdb()
		if err != nil {
			return errors.New("#2 Opendb SetupDB2 " + err.Error())
		}
		db.conn, err = sql.Open("sqlite3", db.Databasepath) // Open the created SQLite File
		if err != nil {
			return errors.New("#3 Opendb Open2 " + err.Error())
		}
	}
	return err
}
func (db *DBtype) Setupdb() error {
	var err error
	if db.Databasepath == "" {
		if !checkGUI() {
			db.Databasepath = fyne.CurrentApp().Preferences().String("pfsmsdb")
		} else {
			db.Databasepath = "pfsms.db"
		}
	}
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
	if db.conn != nil {
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

	s += "CREATE TABLE tblQueue (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, "
	s += "tstamp VARCHAR(20), groupname VARCHAR(100), phone VARCHAR(20), message TEXT);"
	sq := strings.Split(s, ";")
	for i := 0; i < len(sq); i++ {
		if len(sq[i]) < 10 {
			continue
		}
		db.statement, err = db.conn.Prepare(sq[i]) // Prepare SQL Statement
		if err != nil {
			if err.Error() == "table tblCustomers already exists" {
				err = nil
				return err
			}
			log.Println("#1 CreateTables: '", sq[i], "'", err.Error())
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
// CUSTOMERS FUNCTIONS
func (db *DBtype) ShowCustomers() [][]string {
	var phone, firstname, lastname string
	var data [][]string
	var err error
	db.Opendb()
	sq := "SELECT phone,firstname,lastname FROM tblCustomers ORDER BY phone ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ShowCustomers Query error:", err.Error())
		return nil
	}
	for rows.Next() {
		err = rows.Scan(&phone, &firstname, &lastname)
		if err != nil {
			log.Println("#2 ShowCustomers Scan error:", err.Error())
			return nil
		}
		data = append(data, []string{phone, firstname, lastname})
	}
	return data
}
func (db *DBtype) DeleteAllCustomers() {
	var sq string
	var err error
	db.Opendb()
	sq = "DELETE FROM tblCustomers;"
	db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
	if err != nil {
		log.Println("#1 DeleteCustomers Prepare DELETE", err.Error())
		return
	}
	_, err = db.statement.Exec() // Execute SQL Statements
	if err != nil {
		log.Println("#2 ImportCustomers Failed DELETE tblGroups")
		return
	}
}
func (db *DBtype) ImportCustomers(frfile string) {
	// customer textfile should be phone <<tab>> firstname <<tab>> lastname <<tab>> note <<cr>> <<lf>>
	var err error
	var b0 []byte
	var c int
	var sq, phone, firstname, lastname, note string
	b0, err = os.ReadFile(frfile) // SQL to make tables!
	if err != nil {
		log.Println("#1 ImportCustomers ReadFile", err.Error())
		return
	}
	b := string(b0)
	//fix bad characters
	b = Removebadsqlcharacters(b)
	b = strings.Replace(b, "\n", "", -1)
	// allcurrent := db.ShowCustomers()
	db.Opendb()
	for i := 0; i < len(strings.Split(b, "\r")); i++ {
		b1 := strings.Split(b, "\r")[i]
		b2 := strings.Split(b1, "\t")
		firstname = ""
		lastname = ""
		note = ""
		phone = ""
		if len(b2) > 0 {
			phone = Fixphonenumber(b2[0])
		} else {
			continue
		}
		if len(b2) > 1 {
			firstname = b2[1]
		}
		if len(b2) > 2 {
			lastname = b2[2]
		}
		if len(b2) > 3 {
			note = b2[3]
		}
		if len(phone) < 5 {
			continue
		}
		// check if phonenumber already exists
		var cnt string
		r, err := db.conn.Query("SELECT COUNT(*) AS cnt FROM tblCustomers WHERE phone = '" + phone + "'")
		if err != nil {
			log.Println("#2 ImportCustomers Query failed: ", sq, " ", err.Error())
		} else {
			for r.Next() {
				r.Scan(&cnt)
			}
		}
		if cnt > "0" {
			sq = "UPDATE tblCustomers SET phone='%s',firstname='%s',lastname='%s',note='%s' WHERE phone ='%s'"
			sq = fmt.Sprintf(sq, phone, firstname, lastname, note,phone)
		} else {
			sq = "INSERT INTO tblCustomers (phone,firstname,lastname,note)  VALUES ('%s','%s','%s','%s')"
			sq = fmt.Sprintf(sq, phone, firstname, lastname, note)
		}
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
	log.Printf("Imported %d customers.\r\n", c)
}
func (db *DBtype) ExportCustomers(tofile string) {
	var err error
	var sq, txt, phone, firstname, lastname, note string
	db.Opendb()
	sq = "SELECT phone,firstname,lastname,note FROM tblCustomers ORDER BY phone ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ExportCustomers Query ", err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&phone, &firstname, &lastname, &note)
		if err != nil {
			log.Println("#2 ExpotCustomers Scan ", err.Error())
		}
		if len(phone) > 0 {
			txt += fmt.Sprintf("%s\t%s\t%s\t%s\r\n", phone, firstname, lastname, note)
		}
	}
	if txt == "" {
		// export sample data!!!
		txt = "+46736290839\tPeter\tFreund\r\n"
		txt += "087543169\tLin\tZhang\r\n"
		txt += "0890510\tFröken\tUr\tClock\r\n"
		txt += "004690200\tTelia Support\tUtomlands\r\n"
		txt += "0046771735311\tTre Support\tCompany\r\n"
		txt += "+46708 222 222\tTelenor Support\tCompany\r\n"
		txt += "+12024561111\tWhite\tHouse\tin USA\r\n"
	}
	err = os.WriteFile(tofile, []byte(txt), 0644)
	if err != nil {
		log.Println("#3 ExportCustomers WriteFile ", err.Error())
	}
}
// GROUPS FUNCTIONS
func (db *DBtype) ShowGroups() []string {
	// should show all available groupnames
	var data []string
	var err error
	var groupname string
	db.Opendb()
	sq := "SELECT DISTINCT groupname FROM tblGroups ORDER BY groupname ASC"
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
		data = append(data, groupname)
	}
	return data
}
func (db *DBtype) ShowAllGroups() [][]string {
	// should show all available groupnames
	var data [][]string
	var err error
	var groupname, phone string
	db.Opendb()
	sq := "SELECT groupname,phone FROM tblGroups ORDER BY groupname ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ShowGropAllnames Query error:", err.Error())
		return nil
	}
	for rows.Next() {
		err = rows.Scan(&groupname, &phone)
		if err != nil {
			log.Println("#2 ShowAllGroups Scan ", err.Error())
		}
		data = append(data, []string{groupname, phone})
	}
	return data
}
func (db *DBtype) SaveGroup(groupname string, phones string) {
	var err error
	var sq string
	groupname = Removebadsqlcharacters(groupname)
	p := strings.Split(phones, ",")
	for i := 0; i < len(p); i++ {
		p[i] = Fixphonenumber(p[i])
	}
	db.Opendb()
	sq = fmt.Sprintf("DELETE FROM tblGroups WHERE groupname = '%s'", groupname)
	_, err = db.conn.Exec(sq)
	if err != nil {
		log.Println("#1 SaveGroup DELETE failed ", err.Error())
		return
	}
	for i := 0; i < len(p); i++ {
		sq = fmt.Sprintf("INSERT INTO tblGroups (groupname,phone) VALUES ('%s','%s')", groupname, p[i])
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
}
func (db *DBtype) DeleteGroup(g string) {
	var err error
	var sq string
	db.Opendb()
	sq = fmt.Sprintf("DELETE FROM tblGroups WHERE groupname = '%s'", g)
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
}
func (db *DBtype) DeleteAllGroups() {
	var sq string
	var err error
	db.Opendb()
	sq = "DELETE FROM tblGroups;"
	db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
	if err != nil {
		log.Println("#1 DeleteGroups Prepare DELETE", err.Error())
		return
	}
	_, err = db.statement.Exec() // Execute SQL Statements
	if err != nil {
		log.Println("#2 ImportGroups Failed DELETE tblGroups")
		return
	}
}
func (db *DBtype) ImportGroupsFromString(b string) {
	// b should be tab-separated \r separated string with -> grouname <tab> phone <CR>
	// since import is done directory from field in Messages... (not from file)
	var sq string
	var err error
	db.Opendb()
	b = Removebadsqlcharacters(b)
	for i := 0; i < len(strings.Split(b, "\r")); i++ {
		b1 := strings.Split(b, "\r")[i]
		b2 := strings.Split(b1, "\t")
		if len(b2) < 2 {
			continue
		}
		b2[1] = Fixphonenumber(b2[1])
		if len(b2[1]) < 5 {
			continue
		}
		sq = "INSERT INTO tblGroups (groupname,phone)"
		sq += fmt.Sprintf(" VALUES ('%s','%s')", b2[0], b2[1])
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
}
func (db *DBtype) ImportGroups(frfile string) {
	// b should be tab-separated \r separated string with -> grouname <tab> phone <CR>
	// since import is done directory from field in Messages... (not from file)
	var sq string
	var err error
	var b0 []byte
	var c int
	b0, err = os.ReadFile(frfile) // SQL to make tables!
	if err != nil {
		log.Println("#1 ImportGroups ReadFile", err.Error())
		return
	}
	b := string(b0)
	//fix bad characters
	b = Removebadsqlcharacters(b)
	b = strings.Replace(b, "\n", "", -1)
	// allcurrent := db.ShowGroups()
	db.Opendb()
	for i := 0; i < len(strings.Split(b, "\r")); i++ {
		b1 := strings.Split(b, "\r")[i]
		b2 := strings.Split(b1, "\t")
		if len(b2) < 2 {
			continue
		}
		b2[1] = Fixphonenumber(b2[1])
		if len(b2[1]) < 5 {
			continue
		}
		sq = "INSERT INTO tblGroups (groupname,phone)"
		sq += fmt.Sprintf(" VALUES ('%s','%s')", b2[0], b2[1])
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
		c++
	}
	log.Printf("Imported %d groups.\r\n", c)
}
func (db *DBtype) ExportGroups(tofile string) {
	var err error
	var sq, txt, groupname, phone string
	db.Opendb()
	sq = "SELECT groupname,phone FROM tblGroups ORDER BY groupname ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ExportGroups Query ", err.Error())
		return
	}
	for rows.Next() {
		err = rows.Scan(&groupname, &phone)
		if err != nil {
			log.Println("#2 ExpotGroups Scan ,", err.Error())
		}
		if len(phone) > 0 {
			txt += fmt.Sprintf("%s\t%s\r\n", groupname, phone)
		}
	}
	if txt == "" {
		// export sample data!!!
		txt = "Sample\t0046736290839\r\n"
		txt += "Sample\t004687543169\r\n"
		txt += "Sample\t04690510\r\n"
		txt += "Sample\t0012024561111\r\n"
	}
	err = os.WriteFile(tofile, []byte(txt), 0644)
	if err != nil {
		log.Println("#3 ExportGroups WriteFile,", err.Error())
	}
}
// HISTORY FUNCTIONS
func (db *DBtype) ShowHistory() [][]string {
	var r [][]string
	var err error
	var sq, tstamp, groupname, phone, message string
	db.Opendb()
	sq = "SELECT tstamp,groupname,phone,message FROM tblHistory ORDER BY tstamp ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ShowHistory Query ", err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&tstamp, &groupname, &phone, &message)
		if err != nil {
			log.Println("#2 ShowHistory Scan ", err.Error())
		}
		r = append(r, []string{tstamp, groupname, phone, message})
	}
	db.Closedatabase()
	return r
}
func (db *DBtype) ImportHistory(frfile string) {
	// s += "CREATE TABLE tblHistory (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, "
	// s += "tstamp VARCHAR(20), groupname VARCHAR(100), phone VARCHAR(20), message TEXT);"
	// b should be tab-separated \r separated string with -> tstamp <tab> groupname <tab> phone <tab> message <CR>
	// since import is done directory from field in Messages... (not from file)
	var sq string
	var err error
	var b0 []byte
	var c int
	b0, err = os.ReadFile(frfile) // SQL to make tables!
	if err != nil {
		log.Println("#1 ImportGroups ReadFile", err.Error())
		return
	}
	b := string(b0)
	//fix bad characters
	b = Removebadsqlcharacters(b)
	b = strings.Replace(b, "\n", "", -1)
	// allcurrent := db.ShowGroups()
	db.Opendb()
	b = Removebadsqlcharacters(b)
	for i := 0; i < len(strings.Split(b, "\r")); i++ {
		b1 := strings.Split(b, "\r")[i]
		b2 := strings.Split(b1, "\t")
		if len(b2) < 2 {
			continue
		}
		b2[1] = Fixphonenumber(b2[2])
		if len(b2[2]) < 5 {
			continue
		}
		sq = "INSERT INTO tblHistory (tstamp,groupname,phone,message)"
		sq += fmt.Sprintf(" VALUES ('%s','%s','%s','%s')", b2[0], b2[1], b2[2], b2[3])
		db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
		if err != nil {
			log.Println("#3 ImportHistory prepare", err.Error())
			return
		}
		db.reply, err = db.statement.Exec() // Execute SQL Statements
		if err != nil {
			log.Println("#4 ImportHistory Exec ", err.Error())
			return
		}
		c++
	}
	log.Println("Imported ",c," history.")
}
func (db *DBtype) ExportHistory(tofile string) {
	// s += "CREATE TABLE tblHistory (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, "
	// s += "tstamp VARCHAR(20), phone VARCHAR(20), groupname VARCHAR(100), message TEXT);"
	var err error
	var sq, txt, tstamp, groupname, phone, message string
	db.Opendb()
	sq = "SELECT tstamp,groupname,phone,message FROM tblHistory ORDER BY tstamp ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ExportHistory Query ", err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&tstamp, &groupname, &phone, &message)
		if err != nil {
			log.Println("#2 ExportHistory Scan ", err.Error())
		}
		message = Removebadsqlcharacters(message)
		message = Showdebugmsg(message)
		ts:=fmt.Sprintf("%s-%s-%s\t%s:%s:%s",tstamp[0:4],tstamp[4:6],tstamp[6:8],tstamp[8:10],tstamp[10:12],tstamp[12:14])
		txt += fmt.Sprintf("%s\t%s\t%s\t%s\r\n", ts, groupname, phone, message)
	}
	if txt == "" {
		// export sample History!!!
		txt = "2024-08-01\t13:50:55\ttest\t0046736290839\tThis is a test message 1\\r\\nwith 2 lines\r\n"
		txt += "2024-08-01\t13:50:55\ttest\t0046736290839\tThis is a test message 2\\r\\nwith 3 lines\\r\\nwith 2 lines\r\n"
		txt += "2024-08-01\t13:50:55\ttest\t0046736290839\tThis is a test message 3\r\n"
		txt += "2024-08-01\t13:50:55\ttest\t0046736290839\tThis is a test message 4\r\n"
	}
	err = os.WriteFile(tofile, []byte(txt), 0644)
	if err != nil {
		log.Println("#3 ExportHistory WriteFile ", err.Error())
	}
}
func (db *DBtype) SaveHistory(r []string) error {
	// resulting string with history from pfmobile = tstamp,phone,message
	// message in \"\"
	var sq string
	var err error
	db.Opendb()
	// tstamp := r[0]
	// phone := r[1]
	// groupname := r[2]
	// message := r[3]
	sq = "INSERT INTO tblHistory (tstamp,groupname,phone,message)"
	sq += fmt.Sprintf(" VALUES ('%s','%s','%s','%s')", r[0], r[1], r[2], r[3])
	db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
	if err != nil {
		return errors.New("#1 SaveHistory Prepare: "+ err.Error())
	}
	db.reply, err = db.statement.Exec() // Execute SQL Statements
	if err != nil {
		return errors.New("#2 SaveHistory Exec: "+ err.Error())
	}
	return err
}
// QUE FUNCTIONS
func (db *DBtype) AddToQueue(r []string,groupname string,sendtext string) error {
	var tstamp string
	var sq string
	var err error
	db.Opendb()
	for i:=0;i<len(r);i++ {
		r[i]=Removebadsqlcharacters(r[i])
		b:=strings.Split(r[i],"\t")
		tstamp = time.Now().Format("20060102150405")
		phone:=Fixphonenumber(b[0])
		if len(phone)<5 {
			return errors.New("Invalid phone number: " + phone)
		}
		if strings.Contains(sendtext, "<<Fname>>") || strings.Contains(sendtext, "<<Lname>>") {
			fname:=Removebadsqlcharacters(b[1])
			lname:=Removebadsqlcharacters(b[2])
			sendtext = strings.ReplaceAll(sendtext, "<<Fname>>", fname)
			sendtext = strings.ReplaceAll(sendtext, "<<Lname>>", lname)
		}
		msg:=Removebadsqlcharacters(sendtext)
		sq = "INSERT INTO tblQueue (tstamp,phone,groupname,message)"
		sq += fmt.Sprintf(" VALUES ('%s','%s','%s','%s')", tstamp, phone, groupname, msg)
		db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
		if err != nil {
			return errors.New("AddToQueue Prepare error: " + err.Error())
		}
		db.reply, err = db.statement.Exec() // Execute SQL Statements
		if err != nil {
			return errors.New("AddToQueue Exec error: " + err.Error())
		}
	}
	return err
}
func (db *DBtype) DeletefromQueue(id string) error {
	// resulting string with history from pfmobile = tstamp,phone,message
	// message in \"\"
	var sq string
	var err error
	db.Opendb()
	sq = fmt.Sprintf("DELETE FROM tblQueue WHERE ID = %s", id)
	db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
	if err != nil {
		return errors.New("DeletefromQueue Prepare error: " + err.Error())
	}
	db.reply, err = db.statement.Exec() // Execute SQL Statements
	if err != nil {
		return errors.New("DeletefromQueue Exec error: " + err.Error())
	}
	return err
}
func (db *DBtype) ShowQueue() [][]string {
	var r [][]string
	var fixtstamp string
	var err error
	var sq, tstamp,  phone string
	db.Opendb()
	sq = "SELECT tstamp,phone FROM tblQueue ORDER BY tstamp ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ShowQueue Query ", err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&tstamp, &phone)
		if err != nil {
			log.Println("#2 ShowQueue Scan ", err.Error())
		}
		if len(tstamp)>13 {
			fixtstamp=fmt.Sprintf("%s-%s-%s %s:%s:%s",tstamp[0:4],tstamp[4:6],tstamp[6:8],tstamp[8:10],tstamp[10:12],tstamp[12:14])
		}
		r = append(r, []string{fixtstamp, phone})
	}
	db.Closedatabase()
	return r
}
func (db *DBtype) DeleteAllQueue() error {
	var sq string
	var err error
	db.Opendb()
	sq = "DELETE FROM tblQueue;"
	db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
	if err != nil {
		return errors.New("#1 DeleteAllQueue Prepare DELETE"+ err.Error())
	}
	db.reply, err = db.statement.Exec() // Execute SQL Statements
	if err != nil {
		return errors.New("#2 DeleteAllQueue Exec "+ err.Error())
	}
	return err
}
func (db *DBtype) ImportQueue(b string) {
	// b should be tab-separated \r separated string with -> tstamp <tab> groupname <tab> phone <tab> message <CR>
	// since import is done directory from field in Messages... (not from file)
	var sq string
	var err error
	db.Opendb()
	b = Removebadsqlcharacters(b)
	for i := 0; i < len(strings.Split(b, "\r")); i++ {
		b1 := strings.Split(b, "\r")[i]
		b2 := strings.Split(b1, "\t")
		if len(b2) < 2 {
			continue
		}
		b2[2] = Fixphonenumber(b2[2])
		if len(b2[2]) < 5 {
			continue
		}
		sq = "INSERT INTO tblQueue (tstamp,groupname,phone,message)"
		sq += fmt.Sprintf(" VALUES ('%s','%s','%s','%s')", b2[0], b2[1], b2[2], b2[3])
		db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
		if err != nil {
			log.Println("#3 ImportQueue prepare", err.Error())
			return
		}	
		db.reply, err = db.statement.Exec() // Execute SQL Statements
		if err != nil {
			log.Println("#4 ImportQueue Exec ", err.Error())
			return
		}
	}
}
func (db *DBtype) ExportQueue(tofile string) error {
	var err error
	var sq, txt, tstamp, groupname, phone, message string
	db.Opendb()
	sq = "SELECT tstamp,groupname,phone,message FROM tblQueue ORDER BY tstamp ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		return errors.New("#1 ExportQueue Query " + err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&tstamp, &groupname, &phone, &message)
		if err != nil {
			return errors.New("#2 ExportQueue Scan " + err.Error())
		}
		message = Removebadsqlcharacters(message)
		message = Showdebugmsg(message)
		ts:=fmt.Sprintf("%s-%s-%s\t%s:%s:%s",tstamp[0:4],tstamp[4:6],tstamp[6:8],tstamp[8:10],tstamp[10:12],tstamp[12:14])
		txt += fmt.Sprintf("%s\t%s\t%s\t%s\r\n", ts, groupname, phone, message)
	}
	if txt == "" {
		// export sample Que!!!
		txt = "2024-08-01\t13:50:55\ttest\t0046736290839\tThis is a test message 1\\r\\nwith 2 lines\r\n"
		txt += "2024-08-01\t13:50:55\ttest\t0046736290839\tThis is a test message 2\\r\\nwith 3 lines\\r\\nwith 2 lines\r\n"
		txt += "2024-08-01\t13:50:55\ttest\t0046736290839\tThis is a test message 3\r\n"
		txt += "2024-08-01\t13:50:55\ttest\t0046736290839\tThis is a test message 4\r\n"
	}
	err = os.WriteFile(tofile, []byte(txt), 0644)
	if err != nil {
		return errors.New("#3 ExportQueue WriteFile " + err.Error())
	}
	return err
}
func (db *DBtype) CountinQueue() int {
	var err error
	var cnt int
	db.Opendb()
	r, err := db.conn.Query("SELECT COUNT(*) AS cnt FROM tblQueue")
	if err != nil {
		log.Println("#1 CountQueue Query failed: ", err.Error())
		return 0
	}
	for r.Next() {
		r.Scan(&cnt)
	}
	db.Closedatabase()
	return cnt
}
func (db *DBtype) GetNextinQueue() ([]string,error) {
	var err error
	var id, tstamp, groupname, phone, message string
	db.Opendb()
	rows, err := db.conn.Query("SELECT id,tstamp,groupname,phone,message FROM tblQueue ORDER BY tstamp ASC LIMIT 1")
	if err != nil {
		log.Println("#1 GetNextinQueue Query failed: ", err.Error())
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&id, &tstamp, &groupname, &phone, &message)
		if err != nil {
			log.Println("#2 GetNextinQueue Scan failed: ", err.Error())
			return nil, err
		}
	}
	db.Closedatabase()
	return []string{id, tstamp, groupname, phone, message}, err
}
// OTHER FUNCTIONS
func (db *DBtype) GetFname(phone string) string {
	var firstname string
	db.Opendb()
	rows, err := db.conn.Query("SELECT firstname FROM tblCustomers WHERE phone = '" + phone + "'")
	if err != nil {
		fmt.Println("#2 ShowCustomers Query error:", err.Error())
		return ""
	}
	for rows.Next() {
		err = rows.Scan(&firstname)
		if err != nil {
			log.Println("rows.Scan failed in GetFname")
		}
	}
	return firstname
}
func (db *DBtype) GetLname(phone string) string {
	var lastname string
	db.Opendb()
	rows, err := db.conn.Query("SELECT lastname FROM tblCustomers WHERE phone = '" + phone + "'")
	if err != nil {
		return ""
	}
	for rows.Next() {
		if err != nil {
			fmt.Println("#2 ShowCustomers Query error:", err.Error())
			return ""
		}
		err = rows.Scan(&lastname)
		if err != nil {
			log.Println("rows.Scan failed in GetFname")
		}
	}
	return lastname
}
func Fixphonenumber(pn string) string {
	var cc string
	var cci string = "00"
	// pn phonenumber  cc coutrycode
	// Sweden (+46) converts to 0046
	if checkGUI() == false {
		cc = "Sweden(+46)"
	} else {
		cc = fyne.CurrentApp().Preferences().StringWithFallback("mobilecountry", "Sweden(+46)")
	}
	for i := 0; i < len(cc); i++ {
		if strings.Index("0123456789", string(cc[i])) > 0 {
			cci += string(cc[i])
		}
	}
	var cc2, pn2 string
	for i := 0; i < len(cc); i++ {
		if strings.Index("0123456789", string(cc[i])) > 0 {
			cc2 += string(cc[i])
		}
	}
	for i := 0; i < len(pn); i++ {
		if strings.Index("+0123456789", string(pn[i])) > 0 {
			pn2 += string(pn[i])
		}
	}
	if len(pn2) < 6 {
		return ""
	}
	if pn2[0:2] == string("00") {
		return pn2
	}
	if string(pn2[0]) == "0" {
		return "00" + cc2 + pn2[1:]
	}
	if string(pn2[0]) == "+" {
		return "00" + pn2[1:]
	}
	return "00" + cc2 + pn2
}
func Removebadsqlcharacters(v string) string {
	v = strings.Replace(v, "´", "", -1)
	v = strings.Replace(v, "`", "", -1)
	v = strings.Replace(v, "@", "", -1)
	v = strings.Replace(v, "'", "", -1)
	v = strings.Replace(v, "\"", "", -1)
	v = strings.Replace(v, "_", "", -1)
	v = strings.Replace(v, "%", "", -1)
	v = strings.Replace(v, "#", "", -1)
	v = strings.Replace(v, "/", "", -1)
	return v
}
func Showdebugmsg(s string) string {
	r2 := s
	r2 = strings.Replace(r2, string(rune(13)), "\\r", -1)
	r2 = strings.Replace(r2, string(rune(10)), "\\n", -1)
	r2 = strings.Replace(r2, string(rune(0)), "\\0", -1)
	r2 = strings.Replace(r2, string(rune(9)), "\\t", -1)
	r2 = strings.Replace(r2, string(rune(26)), "\\z", -1)
	return r2
}
