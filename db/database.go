package db

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
			panic(err.Error())
		}
	}
}
func (db *DBtype) Setupdb() error {
	var err error
	var path string
	path, err = os.UserHomeDir()
	if err!=nil {
		panic("path")
	}
	path = fmt.Sprintf("%s%c%s",path ,os.PathSeparator, "pfsms")
	if _, err = os.Stat(path); err != nil {
		log.Println("#1 Adding folder data: " + path)
		if os.IsNotExist(err) {
			err = os.Mkdir(path, 0755)
			if err!=nil {
				panic(err.Error())
			}
			// file does not exist
		} else {
			panic(err.Error())
			// other error
		}
	}
	db.Databasepath = path + string(os.PathSeparator) + "pfsms.db"
	if _, err = os.Stat(db.Databasepath); err != nil {
		log.Println("#2 database not found, creating new db: " + db.Databasepath)
		var file *os.File
		file, err = os.Create(db.Databasepath) // Create SQLite file
		if err != nil {
			log.Println("#3 setupdb Failed to create db", err.Error())
			return err
		}
		file.Close()
		err = db.Createtables() // Create Database Tables
		if err != nil {
			log.Println("#4 Could not create tables!", err.Error())
			return err
		} else {
			log.Println("Database tables created")
		}
	}
	return err
}
func (db *DBtype) Closedatabase() error {
	var err error
	db.conn.Close()
	db.conn = nil
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
	var s =	 "CREATE TABLE tblMessages (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, "
		s += "tstamp VARCHAR(25), reference VARCHAR(100), message TEXT);"

		s += "CREATE TABLE tblCustomers (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, "
   		s += "phone VARCHAR(20), firstname VARCHAR(100), lastname VARCHAR(100), note TEXT);"

		s += "CREATE TABLE tblGroups (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, "
   		s += "groupname VARCHAR(100), phone VARCHAR(100));"

		s += "CREATE TABLE tblHistory (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, "
   		s += "tstamp VARCHAR(20), phone VARCHAR(20), reference VARCHAR(100), message TEXT);"

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
func (db *DBtype) ShowMessages() ([][]string,error) {
	var err error
	var ids,reference,message string
	var data [][]string
	var rows *sql.Rows
	db.Opendb()
	rows, err = db.conn.Query("SELECT id,reference,message FROM tblMessages ORDER BY ID ASC")
	if err != nil {
		log.Println("#2 ShowCustomers Query error:", err.Error())
		return nil,err
	}
	for rows.Next() {
		err = rows.Scan(&ids,&reference,&message)
		data=append(data,[]string{ids,reference,message})
	}
	db.Closedatabase()
	return data,err
}
func (db *DBtype) AddMessage(reference string, message string) (int,error) {
	var err error
	var sq string
	var r *sql.Rows
	var reply int
	db.Opendb()
	tstamp := time.Now().Format("2006-01-02 15:04:05")
	// fmt.Println("nanostamp=",nanostamp,"tstamp=",tstamp)	
	sq=fmt.Sprintf("SELECT * FROM tblMessages WHERE reference='%s'",reference)
	r, err = db.conn.Query(sq)
	if r != nil && err==nil {
		r.Next()
		r.Scan(&reply)
		log.Println("#1 cannot insert, reference exists")
		return reply,errors.New("#1 cannot insert, reference exists")
		//reference exists in db, so just Update
	}
	sq="INSERT INTO tblMessages (nanostamp,tstamp,reference,message) "
	sq = fmt.Sprintf("%s VALUES (\"%s\",\"%s\",\"%s\")",sq,tstamp,reference,message)
	db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
	if err != nil {
		log.Println("#2 prepare failed: ", sq, " ", err.Error())
		return -1,err
	}
	db.reply, err = db.statement.Exec() // Execute SQL Statements
	if err != nil {
		log.Println("#3 Exec error failed: ", sq, " ", err.Error())
		return -1,err
	}
	var r0 int64
	r0,err =db.reply.LastInsertId()
	reply=int(r0)
	db.Closedatabase()
	return reply,err
}
func (db *DBtype) UpdateMessage(id int,reference string, message string) error {
	var err error
	db.Opendb()
	tstamp := time.Now().Format("2006-01-02 15:04:05")
	sq:="UPDATE tblMessages (tstamp,reference,message) "
	sq = fmt.Sprintf("%s VALUES (\"%s\",\"%s\",\"%s\")",sq,tstamp,reference,message)
	sq = fmt.Sprintf("%s WHERE id = %d",sq,id)
	db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
	if err != nil {
		log.Println("#1 UpdateMessage prepare failed: ", sq, " ", err.Error())
		return err
	}
	db.reply, err = db.statement.Exec() // Execute SQL Statements
	db.Closedatabase()
	return err
}
func (db *DBtype) DeleteMessage(id int) error {
	var err error
	var sq string
	db.Opendb()
	sq = fmt.Sprintf("DELETE FROM tblMessages WHERE id = %d",id)
	db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
	if err != nil {
		log.Println("#1 prepare failed: ", sq, " ", err.Error())
		return err
	}
	db.reply, err = db.statement.Exec() // Execute SQL Statements
	db.Closedatabase()
	return err
}
func (db *DBtype) ShowCustomers(from int,to int) ([][]string,error) {
	var data [][]string
	var err error
	var id int
	var phone,firstname,lastname string
	db.Opendb()
	// r,_:=db.conn.Exec("SELECT COUNT(*) FROM tblCustomers")
	// fmt.Println(r)
	sq:=fmt.Sprintf("SELECT id,phone,firstname,lastname FROM tblCustomers WHERE id>%d AND id<=%d",from,to)
	rows, err := db.conn.Query(sq)
	if err != nil {
		fmt.Println("#2 ShowCustomers Query error:", err.Error())
		return nil, err
	}
	if !rows.Next() {
		return nil,err
	}
	for rows.Next() {
		err = rows.Scan(&id,&phone,&firstname,&lastname)
		data=append(data,[]string{fmt.Sprintf("%d: %s   %s %s",id,phone,firstname,lastname)})
	}
	db.Closedatabase()
	return data, err
}
func (db *DBtype) ShowGroupnames() ([][]string,error) {
	// should show all available groupnames
	var data [][]string
	var err error
	var id int
	var groupname string
	db.Opendb()
	sq:="SELECT DISTINCT groupname FROM tblGroups ORDER BY groupname ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		fmt.Println("#2 ShowGropupnames Query error:", err.Error())
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&id,&groupname)
		data=append(data,[]string{groupname})
	}
	db.Closedatabase()
	return data, err
}
func (db *DBtype) ImportCustomers(frfile string) error {
	// customer textfile should be phone <<tab>> firstname <<tab>> lastname <<tab>> note <<cr>> <<lf>>
	var err error
	var b0 []byte
	var sq,phone,firstname,lastname,note string
    b0, err = os.ReadFile(frfile) // SQL to make tables!
    if err != nil {
        log.Println("#1 ImportCustomers", err.Error())
    }
	b:=string(b0)
	b= strings.Replace(b,"\n","",-1)
	cc := fyne.CurrentApp().Preferences().StringWithFallback("mobileCountry", "Sweden(+46)")
	var cci string ="00"
	for i:=0;i<len(cc);i++ {
		if strings.Index("0123456789",string(cc[i]))>0 {
			cci +=string(cc[i])
		}
	}

	db.Opendb()
	for i:=0;i<len(strings.Split(b,"\r"));i++ {
		b1:=strings.Split(b,"\r")[i]
		b2:=strings.Split(b1,"\t")
		if i % 100==0 {
			fmt.Println(i,b2)
		}
		// check if phonenumber already exists
		firstname=""
		lastname=""
		note=""
		phone=""
		if len(b2)>0 {
			phone = Fixphonenumber(b2[0],cci)
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
		var cnt string
		r, err := db.conn.Query("SELECT COUNT(*) AS cnt FROM tblCustomers WHERE phone = '"+phone+"'")
		if err==nil {
			for r.Next() {
				err = r.Scan(&cnt)
				if err!=nil {
					log.Println("#1 ImportCustomers Error Scan)",err.Error())
				}
				fmt.Println(cnt)
			}
		}
			if r.Next() {
			continue
		}
		sq ="INSERT INTO tblCustomers (phone,firstname,lastname,note)"
		sq = fmt.Sprintf("%s VALUES ('%s','%s','%s','%s')",sq,phone,firstname,lastname,note)
		db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
		if err != nil {
			log.Println("#2 ImportCustomers prepare failed: ", sq, " ", err.Error())
			return err
		}
		db.reply, err = db.statement.Exec() // Execute SQL Statements
		if err != nil {
			log.Println("#3 ImportCustomers Exec failed: ", sq, " ", err.Error())
			return err
		}
	}
	db.Closedatabase()
	return err
}
func (db *DBtype) ExportCustomers(tofile string) error {
	var err error
	var sq,txt,phone,firstname,lastname,note string
	db.Opendb()
	sq ="SELECT phone,firstname,lastname,note FROM tblCustomers ORDER BY phone ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ExportCustomers ", err.Error())
	}
	if !rows.Next() {
		// export sample data!!!
		txt ="+46736290839\tPeter\tFreund\r\n"
		txt +="087543169\tLin\tZhang\r\n"
		txt +="04690510\tFröken\tUr\tClock\r\n"
		txt +="+12024561111\tWhite\tHouse\tin USA\r\n"
	} else {
		for rows.Next() {
			err = rows.Scan(&phone,&firstname,&lastname,&note)
			if err!=nil {
				log.Println("#2 ExpotCustomers error,",err.Error())
			}
			txt +=fmt.Sprintf("%s\t%s\t%s\t%s\r\n",phone,firstname,lastname,note)
		}
	}
	f, err := os.OpenFile(tofile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("#3 ExpotCustomers error,",err.Error())
	}
	defer f.Close()
	if _, err := f.WriteString(txt); err != nil {
		log.Println("#4 ExpotCustomers error,",err.Error())
	}
	db.Closedatabase()
	return err
}
func (db *DBtype) ImportGroups(frfile string) error {
	// customer textfile should be phone <<tab>> firstname <<tab>> lastname <<tab>> note <<cr>> <<lf>>
	var err error
	var b0 []byte
	var sq string
    b0, err = os.ReadFile(frfile) // SQL to make tables!
    if err != nil {
        log.Println("#1 ImportGroups", err.Error())
    }
	b:=string(b0)
	b= strings.Replace(b,"\n","",-1)
	db.Opendb()
	for i:=0;i<len(strings.Split(b,"\r"));i++ {
		b1:=strings.Split(b,"\r")[i]
		b2:=strings.Split(b1,"\t")
		sq ="INSERT INTO tblGroups (groupname,phone)"
		sq = fmt.Sprintf("%s VALUES (\"%s\",\"%s\")",sq,b2[0],b2[1])
		db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
		if err != nil {
			log.Println("#2 ImportGroups prepare", err.Error())
			return err
		}
		db.reply, err = db.statement.Exec() // Execute SQL Statements
		if err != nil {
			log.Println("#2 ImportGroups Exec", err.Error())
			return err
		}
	}
	db.Closedatabase()
	return err
}
func (db *DBtype) ExportGroups(tofile string) error {
	var err error
	var sq,txt,groupname,phone string
	db.Opendb()
	sq ="SELECT groupname,phone FROM tblGroups ORDER BY groupname ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ExportGroups ", err.Error())
	}
	if !rows.Next() {
		// export sample data!!!
		txt ="Sample\t0046736290839\r\n"
		txt +="Sample\t004687543169\r\n"
		txt +="Sample\t04690510\r\n"
		txt +="Sample\t0012024561111\r\n"
	} else {
		for rows.Next() {
			err = rows.Scan(&groupname,&phone)
			if err!=nil {
				log.Println("#2 ExpotGroups error,",err.Error())
			}
			txt +=fmt.Sprintf("%s\t%s\r\n",groupname,phone)
		}
	}
	f, err := os.OpenFile(tofile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("#4 ExpotGroups error,",err.Error())
	}
	defer f.Close()
	if _, err := f.WriteString(txt); err != nil {
		log.Println("#4 ExpotGroups error,",err.Error())
	}
	db.Closedatabase()
	return err
}
func Fixphonenumber(pn string,cc string) string {
	// pn phonenumber  cc coutrycode
	// Sweden (+46) converts to 0046
	var cci string ="00"
	if len(pn)<5 {
		return ""
	}
	for i:=0;i<len(cc);i++ {
		if strings.Index("0123456789",string(cc[i]))>0 {
			cci +=string(cc[i])
		}
	}
	if pn[0:2]==string("00") {
		return pn
	}
	if string(pn[0])=="0" {
		return  cci+pn[1:]
	}
	if string(pn[0])=="+" {
		return "00"+pn[1:]
	}
	return cci+pn
}
