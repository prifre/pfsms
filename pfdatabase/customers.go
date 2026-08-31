package pfdatabase

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"pfsms/general"
)

// CUSTOMERS FUNCTIONS
func (db *DBtype) ShowCustomers() [][]string {
	var phone, firstname, lastname string
	var data [][]string

	if err := db.Opendb(); err != nil {
		return nil
	}

	sq := "SELECT phone, firstname, lastname FROM tblCustomers ORDER BY phone ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ShowCustomers Query error:", err.Error())
		return nil
	}
	defer rows.Close()

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
	if err := db.Opendb(); err != nil {
		return
	}

	var err error
	db.reply, err = db.conn.Exec("DELETE FROM tblCustomers;")
	if err != nil {
		log.Println("#1 DeleteCustomers Failed DELETE:", err.Error())
	}
}

func (db *DBtype) ImportCustomers(frfile string) {
	b0, err := os.ReadFile(frfile)
	if err != nil {
		log.Println("#1 ImportCustomers ReadFile:", err.Error())
		return
	}

	b := string(b0)
	b = strings.ReplaceAll(b, "\n", "")

	if err := db.Opendb(); err != nil {
		return
	}

	// Starta transaktion för snabbare import
	tx, err := db.conn.Begin()
	if err != nil {
		log.Println("#2 ImportCustomers Begin transaction failed:", err.Error())
		return
	}
	// Om något kraschar rullas ändringarna tillbaka automatiskt
	defer tx.Rollback()

	// Förbered statements inom transaktionen
	stmtCheck, err := tx.Prepare("SELECT COUNT(*) FROM tblCustomers WHERE phone = ?")
	if err != nil {
		log.Println("#3 ImportCustomers Prepare Check failed:", err.Error())
		return
	}
	defer stmtCheck.Close()

	stmtUpdate, err := tx.Prepare("UPDATE tblCustomers SET firstname=?, lastname=?, note=? WHERE phone=?")
	if err != nil {
		log.Println("#4 ImportCustomers Prepare Update failed:", err.Error())
		return
	}
	defer stmtUpdate.Close()

	stmtInsert, err := tx.Prepare("INSERT INTO tblCustomers (phone, firstname, lastname, note) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println("#5 ImportCustomers Prepare Insert failed:", err.Error())
		return
	}
	defer stmtInsert.Close()

	lines := strings.Split(b, "\r")
	var count int

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		phone := ""
		firstname := ""
		lastname := ""
		note := ""

		if len(parts) > 0 {
			phone = general.Fixphonenumber(parts[0])
		}
		if len(parts) > 1 {
			firstname = parts[1]
		}
		if len(parts) > 2 {
			lastname = parts[2]
		}
		if len(parts) > 3 {
			note = parts[3]
		}

		if len(phone) < 5 {
			continue
		}

		var cnt int
		err := stmtCheck.QueryRow(phone).Scan(&cnt)
		if err != nil && err != sql.ErrNoRows {
			log.Println("#6 ImportCustomers QueryRow failed:", err.Error())
			continue
		}

		if cnt > 0 {
			// UPDATE
			_, err = stmtUpdate.Exec(firstname, lastname, note, phone)
			if err != nil {
				log.Println("#7 ImportCustomers UPDATE failed:", err.Error())
				continue
			}
		} else {
			// INSERT
			_, err = stmtInsert.Exec(phone, firstname, lastname, note)
			if err != nil {
				log.Println("#8 ImportCustomers INSERT failed:", err.Error())
				continue
			}
		}
		count++
	}

	// Verkställ alla ändringar samtidigt till disken
	if err := tx.Commit(); err != nil {
		log.Println("#9 ImportCustomers Commit failed:", err.Error())
		return
	}

	log.Printf("Imported %d customers.\r\n", count)
}

func (db *DBtype) ExportCustomers(tofile string) {
	if err := db.Opendb(); err != nil {
		return
	}

	sq := "SELECT phone, firstname, lastname, note FROM tblCustomers ORDER BY phone ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ExportCustomers Query:", err.Error())
		return
	}
	defer rows.Close()

	var txt string
	var c int
	var phone, firstname, lastname, note string

	for rows.Next() {
		c++
		err = rows.Scan(&phone, &firstname, &lastname, &note)
		if err != nil {
			log.Println("#2 ExportCustomers Scan:", err.Error())
			continue
		}
		if len(phone) > 0 {
			txt += fmt.Sprintf("%s\t%s\t%s\t%s\r\n", phone, firstname, lastname, note)
		}
	}

	if txt == "" {
		// Sample data
		txt = "+46736290839\tPeter\tFreund\r\n" +
			"087543169\tLin\tZhang\r\n"
	}

	err = os.WriteFile(tofile, []byte(txt), 0644)
	if err != nil {
		log.Println("#3 ExportCustomers WriteFile:", err.Error())
	}
	log.Printf("Exported %d customers.\r\n", c)
}
