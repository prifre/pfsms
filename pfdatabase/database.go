package pfdatabase

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"pfsms/general"

	"fyne.io/fyne/v2"
	_ "github.com/mattn/go-sqlite3"
)

type DBtype struct {
	conn         *sql.DB
	statement    *sql.Stmt
	reply        sql.Result
	Databasepath string
}

func (db *DBtype) Opendb() error {
	if db.conn != nil {
		return nil // Redan öppen
	}

	// Sätt standard-sökväg om den inte är angiven
	if db.Databasepath == "" {
		if fyne.CurrentApp() != nil {
			db.Databasepath = fyne.CurrentApp().Preferences().String("pfsmsdb")
		}
		if db.Databasepath == "" {
			db.Databasepath = "pfsms.db"
		}
	}

	_, err := os.Stat(db.Databasepath)
	if errors.Is(err, os.ErrNotExist) {
		err = db.Setupdb()
		if err != nil {
			return fmt.Errorf("#1 Opendb SetupDB: %w", err)
		}
	}

	db.conn, err = sql.Open("sqlite3", db.Databasepath)
	if err != nil {
		return fmt.Errorf("#2 Opendb Open: %w", err)
	}

	return nil
}

func (db *DBtype) Setupdb() error {
	if db.Databasepath == "" {
		if fyne.CurrentApp() != nil {
			db.Databasepath = fyne.CurrentApp().Preferences().String("pfsmsdb")
		}
		if db.Databasepath == "" {
			db.Databasepath = "pfsms.db"
		}
	}

	if _, err := os.Stat(db.Databasepath); os.IsNotExist(err) {
		log.Println("#1 Setupdb database not found, creating new db: " + db.Databasepath)
		file, err := os.Create(db.Databasepath)
		if err != nil {
			log.Println("#2 Setupdb Failed to create db:", err)
			return err
		}
		file.Close()

		err = db.Createtables()
		if err != nil {
			log.Println("#3 Setupdb Could not create tables!:", err)
			return err
		}
		log.Println("Database tables created successfully")
	}
	return nil
}

func (db *DBtype) Closedatabase() error {
	if db.conn != nil {
		err := db.conn.Close()
		db.conn = nil
		return err
	}
	return nil
}

func (db *DBtype) Createtables() error {
	var err error
	if db.conn == nil {
		db.conn, err = sql.Open("sqlite3", db.Databasepath)
		if err != nil {
			return fmt.Errorf("#1 Createtables sql.Open: %w", err)
		}
	}

	// Kontrollera om tabellen finns
	var tableName string
	err = db.conn.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='tblCustomers';").Scan(&tableName)
	if err == nil && tableName == "tblCustomers" {
		return nil // Tabellerna finns redan
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS tblCustomers (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			phone VARCHAR(20), 
			firstname VARCHAR(100), 
			lastname VARCHAR(100), 
			note TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS tblGroups (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			groupname VARCHAR(100), 
			phone VARCHAR(100)
		);`,
		`CREATE TABLE IF NOT EXISTS tblHistory (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			tstamp VARCHAR(20), 
			groupname VARCHAR(100), 
			phone VARCHAR(20), 
			message TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS tblHashtable (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			hash VARCHAR(100)
		);`,
		`CREATE TABLE IF NOT EXISTS tblQueue (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			tstamp VARCHAR(20), 
			groupname VARCHAR(100), 
			phone VARCHAR(20), 
			message TEXT
		);`,
	}

	for _, q := range queries {
		_, err := db.conn.Exec(q)
		if err != nil {
			log.Printf("#2 CreateTables failed for query [%s]: %v\n", q, err)
			return err
		}
	}

	return nil
}

// GetFname hämtar förnamn säkert via parametriserad fråga
func (db *DBtype) GetFname(phone string) string {
	if err := db.Opendb(); err != nil {
		log.Println("GetFname Opendb error:", err)
		return ""
	}

	var firstname string
	err := db.conn.QueryRow("SELECT firstname FROM tblCustomers WHERE phone = ?", phone).Scan(&firstname)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println("GetFname Query error:", err)
		}
		return ""
	}
	return firstname
}

// GetLname hämtar efternamn säkert via parametriserad fråga
func (db *DBtype) GetLname(phone string) string {
	if err := db.Opendb(); err != nil {
		log.Println("GetLname Opendb error:", err)
		return ""
	}

	var lastname string
	err := db.conn.QueryRow("SELECT lastname FROM tblCustomers WHERE phone = ?", phone).Scan(&lastname)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println("GetLname Query error:", err)
		}
		return ""
	}
	return lastname
}

// Behövs inte längre när parametriserade frågor (SQL placeholders ?) används överallt.
func Removebadsqlcharacters(v string) string {
	return general.Showdebugmsg(v)
}
