package pfdatabase

import (
	"fmt"
	"log"
	"os"
	"strings"

	"pfsms/general"
)

// Hjälpfunktion för att säkert formatera tidsstämplar (YYYYMMDDHHMMSS -> YYYY-MM-DD HH:MM:SS)
func formatTimestamp(tstamp string) string {
	if len(tstamp) < 14 {
		return tstamp // Returnera som den är om den inte stämmer med formatet
	}
	return fmt.Sprintf("%s-%s-%s %s:%s:%s",
		tstamp[0:4], tstamp[4:6], tstamp[6:8],
		tstamp[8:10], tstamp[10:12], tstamp[12:14])
}

// HISTORY FUNCTIONS
func (db *DBtype) ShowHistory() [][]string {
	var r [][]string
	var err error
	var tstamp, groupname, phone, message string

	if err = db.Opendb(); err != nil {
		log.Println("#0 ShowHistory Opendb error:", err)
		return nil
	}

	sq := "SELECT tstamp, groupname, phone, message FROM (SELECT tstamp, groupname, phone, message FROM tblHistory ORDER BY tstamp DESC LIMIT 30) AS subquery ORDER BY tstamp ASC;"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ShowHistory Query error:", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&tstamp, &groupname, &phone, &message)
		if err != nil {
			log.Println("#2 ShowHistory Scan error:", err)
			continue
		}

		fixtstamp := formatTimestamp(tstamp)
		if len(message) > 40 {
			message = message[:40] + "..."
		}
		r = append(r, []string{fixtstamp, groupname, phone, message})
	}
	return r
}

func (db *DBtype) ImportHistory(frfile string) {
	b0, err := os.ReadFile(frfile)
	if err != nil {
		log.Println("#1 ImportHistory ReadFile error:", err)
		return
	}

	b := string(b0)
	b = strings.ReplaceAll(b, "\n", "")

	if err = db.Opendb(); err != nil {
		log.Println("#2 ImportHistory Opendb error:", err)
		return
	}

	sq := "INSERT INTO tblHistory (tstamp, groupname, phone, message) VALUES (?, ?, ?, ?)"
	db.statement, err = db.conn.Prepare(sq)
	if err != nil {
		log.Println("#3 ImportHistory Prepare error:", err)
		return
	}
	defer db.statement.Close()

	var c int
	lines := strings.Split(b, "\r")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "\t")
		if len(parts) < 4 {
			continue
		}

		tstamp := parts[0]
		groupname := parts[1]
		phone := general.Fixphonenumber(parts[2])
		message := parts[3]

		if len(phone) < 5 {
			continue
		}

		db.reply, err = db.statement.Exec(tstamp, groupname, phone, message)
		if err != nil {
			log.Println("#4 ImportHistory Exec error:", err)
			continue
		}
		c++
	}
	log.Printf("Imported %d history records.\r\n", c)
}

func (db *DBtype) ExportHistory(tofile string) {
	var err error
	var txt, tstamp, groupname, phone, message string

	if err = db.Opendb(); err != nil {
		log.Println("#0 ExportHistory Opendb error:", err)
		return
	}

	sq := "SELECT tstamp, groupname, phone, message FROM tblHistory ORDER BY tstamp ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ExportHistory Query error:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&tstamp, &groupname, &phone, &message)
		if err != nil {
			log.Println("#2 ExportHistory Scan error:", err)
			continue
		}

		ts := formatTimestamp(tstamp)
		txt += fmt.Sprintf("%s\t%s\t%s\t%s\r\n", ts, groupname, phone, message)
	}

	if txt == "" {
		txt = "2024-08-01 13:50:55\ttest\t0046736290839\tThis is a test message 1\r\n" +
			"2024-08-01 13:50:55\ttest\t0046736290839\tThis is a test message 2\r\n"
	}

	err = os.WriteFile(tofile, []byte(txt), 0644)
	if err != nil {
		log.Println("#3 ExportHistory WriteFile error:", err)
	}
}

func (db *DBtype) SaveHistory(r []string) error {
	var err error
	// r förväntas innehålla: [0]=tstamp, [1]=groupname, [2]=phone, [3]=message
	if len(r) < 4 {
		return fmt.Errorf("SaveHistory: expected at least 4 elements, got %d", len(r))
	}

	if err := db.Opendb(); err != nil {
		return fmt.Errorf("#0 SaveHistory Opendb error: %w", err)
	}

	tstamp := r[0]
	groupname := r[1]
	phone := general.Fixphonenumber(r[2])
	message := r[3]

	sq := "INSERT INTO tblHistory (tstamp, groupname, phone, message) VALUES (?, ?, ?, ?)"
	db.statement, err = db.conn.Prepare(sq)
	if err != nil {
		return fmt.Errorf("#1 SaveHistory Prepare error: %w", err)
	}
	defer db.statement.Close()

	db.reply, err = db.statement.Exec(tstamp, groupname, phone, message)
	if err != nil {
		return fmt.Errorf("#2 SaveHistory Exec error: %w", err)
	}

	return nil
}
