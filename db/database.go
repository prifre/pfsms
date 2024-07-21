package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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
	db.conn.SetMaxOpenConns(1)
	db.conn.SetMaxIdleConns(0)
	db.conn.SetConnMaxIdleTime(time.Hour * 2)
	db.conn.SetConnMaxLifetime(time.Hour * 2)
}
func (db *DBtype) Setupdb() error {
	var err error
	var path string
	path, err = os.Getwd()
	if err!=nil {
		panic("path")
	}
	if path[len(path)-2:]=="db" {
		path = path[:len(path)-3]
	}
	if path[len(path)-4:]!="data" {
		path = path + string(os.PathSeparator) + "data"
	}
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
	db.Opendb()
	_, table_check := db.conn.Query("SELECT * FROM tblCustomers;")

	if table_check == nil {
		return nil
		//table tblMain exists, so probably all is well...
	}
	//create tables...
    b, err := os.ReadFile("createtables.sql") // SQL to make tables!
    if err != nil {
        fmt.Print(err)
    }
    s := string(b) // convert content to a 'string'	for _, s := range sq {
	s=strings.Replace(s,"\r"," ",-1)
	s=strings.Replace(s,"\n"," ",-1)
	s=strings.TrimSpace(s)
	for i:=0;i<len(strings.Split(s,";"));i++ {
		sq:=strings.Split(s,";")[i]
		if len(sq)<10 {
			continue
		}
		db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
		if err != nil {
			if err.Error() == "table tblCustomers already exists" {
				err = nil
				return err
			}
			log.Println("#1 CreateTables: ", err.Error())
		}
		db.reply, err = db.statement.Exec() // Execute SQL Statements
		if err != nil {
			log.Println("#2 CreateTables failed: ", sq, " ", err.Error(), db.reply)
			return err
		}
	}
	return err
}
func (db *DBtype) Getsql(sq string) ([]string, error) {
// get one value from database quickly...
var err error
	var k []string
	var s sql.NullString
	var s2 string
	db.Opendb()
	rows, err := db.conn.Query(sq)
	if err != nil {
		fmt.Println("#2 Getsql Query error:", err.Error())
		return nil, err
	}
	col, err := rows.Columns()
	if err != nil {
		fmt.Println("#3 Getsql Col error", err.Error())
		return nil, err
	}
	if len(col) > 1 {
		log.Println("#4 Getsql too many columns in query! Do your own query!")
		return nil, fmt.Errorf("too many columns!%v", "")
	}
	var ct []*sql.ColumnType
	ct, err = rows.ColumnTypes()
	if err != nil {
		fmt.Println("#5 Getsql CT error", err.Error())
	}
	for rows.Next() {
		switch strings.ToUpper(ct[0].DatabaseTypeName()) {
		case "INTEGER":
			var x int64
			err = rows.Scan(&x)
			if err != nil {
				fmt.Println("#6 Getsql Scan error", err.Error())
			}
			s2 = fmt.Sprintf("%v", x)
		case "TEXT":
			err = rows.Scan(&s)
			s2 = ""
			if s.Valid {
				s2 = fmt.Sprintf("%v", s.String)
			}
			if err != nil {
				fmt.Println("#6 Getsql Scan error", err.Error())
			}
		default:
			err = rows.Scan(&s)
			// COUNT(*)...
			if s.Valid {
				s2 = fmt.Sprintf("%v", s.String)
			}
		}
		k = append(k, s2)
	}
	return k, err
}
func (db *DBtype) Deleteall(n string) error {
// delete the database file?

	var err error
	var sq []string
	// remove from database
	db.Opendb()
	sq = append(sq, "DELETE FROM tblDustTrak WHERE nanostamp="+n)
	sq = append(sq, "DELETE FROM tblPTrak WHERE nanostamp="+n)
	sq = append(sq, "DELETE FROM tblAeroTrak WHERE nanostamp="+n)
	sq = append(sq, "DELETE FROM tblMain WHERE nanostamp="+n)
	for i := 0; i < len(sq); i++ {
		db.statement, err = db.conn.Prepare(sq[i]) // Prepare SQL Statement
		if err != nil {
			log.Println("#2 deleteall prepare failed: ", sq[i], " ", err.Error())
			return err
		}
		db.reply, err = db.statement.Exec() // Execute SQL Statements
		if err != nil {
			log.Println("#3 deleteall exec failed: ", sq[i], " ", err.Error(), db.reply)
			return err
		}
	}
	return err
}
func (db *DBtype) ImportCustomers(frfile string) error {
	var err error
	var b0 []byte
	var sq string
    b0, err = os.ReadFile(frfile) // SQL to make tables!
    if err != nil {
        fmt.Print(err)
    }
	b:=string(b0)
	db.Opendb()
	for i:=0;i<len(strings.Split(b,"\r\n"));i++ {
		b1:=strings.Split(b,"\r\n")[i]
		b2:=strings.Split(b1,"\t")
		if len(b2)<9 || b2[1]=="Exp.Nota" {
			continue
		}
		if i % 100==0 {
			fmt.Println(i,b2)
		}
		sq ="INSERT INTO tblCustomers (expnote,phone,firstname,lastname,indate,outdate)"
		sq = fmt.Sprintf("%s VALUES (\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\")",sq,b2[1],b2[2],b2[3],b2[4],b2[9],b2[10])
		db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
		if err != nil {
			log.Println("#2 prepare failed: ", sq, " ", err.Error())
			return err
		}
		db.reply, err = db.statement.Exec() // Execute SQL Statements

		
		// KUNDER_ALLA innehåller:
		// Nr	Exp.Nota	Mobilnr	Förnamn	Efternamn	In	Ut	Sign	Regdatum	Indatum	Utdatum	År	Pärm	Väntetid
	}
	return err
}
func (db *DBtype) AddMessage(reference string, message string) error {
	var err error
	db.Opendb()
	nanostamp := time.Now().UnixNano()
	tstamp := time.Now().Format(time.RFC3339)
	fmt.Println("nanostamp=",nanostamp,"tstamp=",tstamp)	
	sq:="INSERT INTO tblMessages (nanostamp,tstamp,reference,message) "
	sq = fmt.Sprintf("%s VALUES (%d,\"%s\",\"%s\",\"%s\")",sq,nanostamp,tstamp,reference,message)
	db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
	if err != nil {
		log.Println("#1 prepare failed: ", sq, " ", err.Error())
		return err
	}
	db.reply, err = db.statement.Exec() // Execute SQL Statements
return err
}
func (db *DBtype) ShowCustomers(from int,to int) ([][]string,error) {
	var data [][]string
	var err error
	var id int
	var phone,firstname,lastname string
	db.Opendb()
	sq:=fmt.Sprintf("SELECT id,phone,firstname,lastname from tblCustomers WHERE id>%d AND id<=%d",from,to)
	rows, err := db.conn.Query(sq)
	if err != nil {
		fmt.Println("#2 ShowCustomers Query error:", err.Error())
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&id,&phone,&firstname,&lastname)
		data=append(data,[]string{fmt.Sprintf("%d",id),phone,firstname,lastname})
	}
	return data, err
}
func (db *DBtype) ShowGroupnames() ([][]string,error) {
	// should show all available groupnames
	var data [][]string
	var err error
	var id int
	var phone,firstname,lastname string
	db.Opendb()
	sq:="SELECT * FROM tblGroupnames ORDER BY groupname ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		fmt.Println("#2 ShowGropupnames Query error:", err.Error())
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&id,&phone,&firstname,&lastname)
		data=append(data,[]string{fmt.Sprintf("%d",id),phone,firstname,lastname})
	}
	return data, err
}

func (db *DBtype) ShowGroups(from int,to int) ([][]string,error) {
	// the ShowGroups should show groups based on CustomerID selected []id
	var data [][]string
	var err error
	data=append(data,[]string{"test","test2","test3"})
	
	// var id int
	// var phone,firstname,lastname string
	// db.Opendb()
	// sq:=fmt.Sprintf("SELECT id FROM tblGroups WHERE id>%d AND id<=%d",from,to)
	// rows, err := db.conn.Query(sq)
	// if err != nil {
	// 	fmt.Println("#2 ShowGroups Query error:", err.Error())
	// 	return nil, err
	// }
	// for rows.Next() {
	// 	err = rows.Scan(&id,&phone,&firstname,&lastname)
	// 	data=append(data,[]string{fmt.Sprintf("%d",id),phone,firstname,lastname})
	// }
	return data, err
}

