package db

import (
	"database/sql"
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
func (db *DBtype) Opendb() error {
	var err error
	// var temp fyne.URI
	if db.conn != nil {
		return nil // allready opened!
	}
	db.conn, err = sql.Open("sqlite3", db.Databasepath) // Open the created SQLite File
	if err != nil {
		log.Fatal("setupdatabase storage.Child error", err.Error())
		db.Setupdb()
	}
	db.conn.SetMaxOpenConns(1)
	db.conn.SetMaxIdleConns(0)
	db.conn.SetConnMaxIdleTime(time.Hour * 2)
	db.conn.SetConnMaxLifetime(time.Hour * 2)
	return err
}
func (db *DBtype) Setupdb() error {
	var err error
	db.Databasepath = "pfsms.db"
	if _, err = os.Stat(db.Databasepath); err == nil {
		err = db.Opendb()
		if err != nil {
			log.Println("#1 setupdb Failed to open db '"+db.Databasepath+"'", db.conn)
			return err
		}
	} else {
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
		err = db.Opendb()
		if err != nil {
			log.Println("#5 setupdb Failed to open db", db.conn)
			return err
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
	err = db.Opendb()
	if err != nil {
		log.Println("#1 CreateTables failed opendb: ", err.Error())
		return err
	}
	// check if table exists
	_, table_check := db.conn.Query("select * from tblMain;")

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
	s=strings.Replace(s,"\r\n\r\n","¤",-1)
	s=strings.Replace(s,"\r\n"," ",-1)
	for i:=0;i<len(strings.Split(s,"¤"));i++ {
		sq:=strings.Split(s,"¤")[i]
		db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
		if err != nil {
			if err.Error() == "table tblMain already exists" {
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
	err = db.Opendb()
	if err != nil {
		log.Println("#1 Getsql opendb error: ", err.Error())
		return nil, err
	}
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
	var err error
	var sq []string
	// remove from database
	err = db.Opendb()
	if err != nil {
		log.Println("#1 deleteall open Failed", err.Error())
	}
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
    b0, err = os.ReadFile(frfile) // SQL to make tables!
    if err != nil {
        fmt.Print(err)
    }
	b:=string(b0)
	for i:=0;i<len(strings.Split(b,"\r\n"));i++ {
		b1:=strings.Split(b,"\r\n")[i]
		b2:=strings.Split(b1,"\t")
	//	fmt.Printf("%s\t%s\t%s\r\n",b2[2],b2[3],b2[4]) // mobilnrm förnamn, efternamn
		sq:="INSERT INTO tblCustomers (expnote,phone,firstname,lastname,indate, outdate)"
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
func (db *DBtype) AddMessage(messagetitle string,message string) error {
	var err error
	nanostamp := time.Now().UnixNano()
	tstamp := time.Now().Format(time.RFC3339)
	fmt.Println("nanostamp=",nanostamp,"tstamp=",tstamp)	
	sq:="INSERT INTO tblMessages (nanostamp,tstamp,messagetitle,message) "
	sq = fmt.Sprintf("%s VALUES (%d,\"%s\",\"%s\",\"%s\")",sq,nanostamp,tstamp,messagetitle,message)
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
	err = db.Opendb()
	if err != nil {
		log.Println("#1 ShowCustomers opendb error: ", err.Error())
		return nil, err
	}
	sq:=fmt.Sprintf("SELECT id,phone,firstname,lastname from tblCustomers WHERE id>%d AND id<=%d",from,to)
	rows, err := db.conn.Query(sq)
	if err != nil {
		fmt.Println("#2 ShoweCustomers Query error:", err.Error())
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&id,&phone,&firstname,&lastname)
		data=append(data,[]string{fmt.Sprintf("%d",id),phone,firstname,lastname})
	}
	return data, err
}

