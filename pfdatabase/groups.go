package pfdatabase

import (
	"fmt"
	"log"
	"os"
	"strings"

	"pfsms/general"
)

// GROUPS FUNCTIONS
func (db *DBtype) ShowGroups() []string {
	var data []string
	var err error
	var groupname string

	if err = db.Opendb(); err != nil {
		log.Println("#0 ShowGroups Opendb error:", err)
		return nil
	}

	sq := "SELECT DISTINCT groupname FROM tblGroups ORDER BY groupname ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ShowGroups Query error:", err)
		return nil
	}
	defer rows.Close() // Stänger resurserna när funktionen är klar

	for rows.Next() {
		err = rows.Scan(&groupname)
		if err != nil {
			log.Println("#2 ShowGroups Scan error:", err)
			return nil
		}
		data = append(data, groupname)
	}
	return data
}

func (db *DBtype) ShowAllGroups() [][]string {
	var data [][]string
	var err error
	var groupname, phone string

	if err = db.Opendb(); err != nil {
		log.Println("#0 ShowAllGroups Opendb error:", err)
		return nil
	}

	sq := "SELECT groupname, phone FROM tblGroups ORDER BY groupname ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ShowAllGroups Query error:", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&groupname, &phone)
		if err != nil {
			log.Println("#2 ShowAllGroups Scan error:", err)
			continue
		}
		data = append(data, []string{groupname, phone})
	}
	return data
}

func (db *DBtype) SaveGroup(groupname string, phones string) {
	var err error

	if err = db.Opendb(); err != nil {
		log.Println("#0 SaveGroup Opendb error:", err)
		return
	}

	p := strings.Split(phones, ",")
	for i := 0; i < len(p); i++ {
		p[i] = general.Fixphonenumber(p[i])
	}

	// Säker DELETE med parametrering
	sq := "DELETE FROM tblGroups WHERE groupname = ?"
	db.statement, err = db.conn.Prepare(sq)
	if err != nil {
		log.Println("#1 SaveGroup DELETE prepare failed:", err)
		return
	}
	db.reply, err = db.statement.Exec(groupname)
	if err != nil {
		log.Println("#2 SaveGroup DELETE exec failed:", err)
		return
	}

	// Säker INSERT
	sqInsert := "INSERT INTO tblGroups (groupname, phone) VALUES (?, ?)"
	db.statement, err = db.conn.Prepare(sqInsert)
	if err != nil {
		log.Println("#3 SaveGroup INSERT prepare failed:", err)
		return
	}

	for i := 0; i < len(p); i++ {
		if strings.TrimSpace(p[i]) == "" {
			continue
		}
		db.reply, err = db.statement.Exec(groupname, p[i])
		if err != nil {
			log.Println("#4 SaveGroup INSERT exec failed:", err)
			return
		}
	}
}

func (db *DBtype) DeleteGroup(g string) {
	var err error

	if err = db.Opendb(); err != nil {
		log.Println("#0 DeleteGroup Opendb error:", err)
		return
	}

	sq := "DELETE FROM tblGroups WHERE groupname = ?"
	db.statement, err = db.conn.Prepare(sq)
	if err != nil {
		log.Println("#1 DeleteGroup Prepare failed:", err)
		return
	}

	db.reply, err = db.statement.Exec(g)
	if err != nil {
		log.Println("#2 DeleteGroup Exec failed:", err)
		return
	}
}

func (db *DBtype) DeleteAllGroups() {
	var err error

	if err = db.Opendb(); err != nil {
		log.Println("#0 DeleteAllGroups Opendb error:", err)
		return
	}

	sq := "DELETE FROM tblGroups;"
	db.statement, err = db.conn.Prepare(sq)
	if err != nil {
		log.Println("#1 DeleteGroups Prepare DELETE failed:", err)
		return
	}

	db.reply, err = db.statement.Exec()
	if err != nil {
		log.Println("#2 DeleteGroups Exec DELETE failed:", err)
		return
	}
}

func (db *DBtype) ImportGroupsFromString(b string) {
	var err error

	if err = db.Opendb(); err != nil {
		log.Println("#0 ImportGroupsFromString Opendb error:", err)
		return
	}

	sq := "INSERT INTO tblGroups (groupname, phone) VALUES (?, ?)"
	db.statement, err = db.conn.Prepare(sq)
	if err != nil {
		log.Println("#1 ImportGroupsFromString prepare failed:", err)
		return
	}

	lines := strings.Split(b, "\r")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}

		phone := general.Fixphonenumber(parts[1])
		if len(phone) < 5 {
			continue
		}

		db.reply, err = db.statement.Exec(parts[0], phone)
		if err != nil {
			log.Println("#2 ImportGroupsFromString Exec failed:", err)
			return
		}
	}
}

func (db *DBtype) ImportGroups(frfile string) {
	b0, err := os.ReadFile(frfile)
	if err != nil {
		log.Println("#1 ImportGroups ReadFile error:", err)
		return
	}

	b := string(b0)
	b = strings.ReplaceAll(b, "\n", "")

	if err = db.Opendb(); err != nil {
		log.Println("#2 ImportGroups Opendb error:", err)
		return
	}

	sq := "INSERT INTO tblGroups (groupname, phone) VALUES (?, ?)"
	db.statement, err = db.conn.Prepare(sq)
	if err != nil {
		log.Println("#3 ImportGroups prepare failed:", err)
		return
	}

	var c int
	lines := strings.Split(b, "\r")
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}

		phone := general.Fixphonenumber(parts[1])
		if len(phone) < 5 {
			continue
		}

		db.reply, err = db.statement.Exec(parts[0], phone)
		if err != nil {
			log.Println("#4 ImportGroups Exec failed:", err)
			return
		}
		c++
	}
	log.Printf("Imported %d groups.\r\n", c)
}

func (db *DBtype) ExportGroups(tofile string) {
	var err error
	var txt, groupname, phone string
	var c int

	if err = db.Opendb(); err != nil {
		log.Println("#0 ExportGroups Opendb error:", err)
		return
	}

	sq := "SELECT groupname, phone FROM tblGroups ORDER BY groupname ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ExportGroups Query error:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		c++
		err = rows.Scan(&groupname, &phone)
		if err != nil {
			log.Println("#2 ExportGroups Scan error:", err)
			continue
		}
		if len(phone) > 0 {
			txt += fmt.Sprintf("%s\t%s\r\n", groupname, phone)
		}
	}

	if txt == "" {
		txt = "Sample\t0046736290839\r\n"
		txt += "Sample\t004687543169\r\n"
		txt += "Sample\t04690510\r\n"
		txt += "Sample\t0012024561111\r\n"
	}

	err = os.WriteFile(tofile, []byte(txt), 0644)
	if err != nil {
		log.Println("#3 ExportGroups WriteFile error:", err)
	}
	log.Printf("Exported %d groups.\r\n", c)
}
