package pfdatabase

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"pfsms/general"
)

// QUE FUNCTIONS
func (db *DBtype) AddToQueue(r []string, groupname string, sendtext string) error {
	var err error

	if err = db.Opendb(); err != nil {
		return fmt.Errorf("AddToQueue Opendb error: %w", err)
	}

	sq := "INSERT INTO tblQueue (tstamp, phone, groupname, message) VALUES (?, ?, ?, ?)"
	db.statement, err = db.conn.Prepare(sq)
	if err != nil {
		return fmt.Errorf("AddToQueue Prepare error: %w", err)
	}
	defer db.statement.Close()

	for i := 0; i < len(r); i++ {
		b := strings.Split(r[i], "\t")
		tstamp := time.Now().Format("20060102150405")
		phone := general.Fixphonenumber(b[0])

		if len(phone) < 5 {
			return fmt.Errorf("invalid phone number: %s", phone)
		}

		currentMsg := sendtext
		if len(b) >= 3 && (strings.Contains(currentMsg, "<<Fname>>") || strings.Contains(currentMsg, "<<Lname>>")) {
			fname := b[1]
			lname := b[2]
			currentMsg = strings.ReplaceAll(currentMsg, "<<Fname>>", fname)
			currentMsg = strings.ReplaceAll(currentMsg, "<<Lname>>", lname)
		}

		db.reply, err = db.statement.Exec(tstamp, phone, groupname, currentMsg)
		if err != nil {
			return fmt.Errorf("AddToQueue Exec error: %w", err)
		}
	}
	return nil
}

func (db *DBtype) DeletefromQueue(id string) error {
	var err error

	if err = db.Opendb(); err != nil {
		return fmt.Errorf("DeletefromQueue Opendb error: %w", err)
	}

	sq := "DELETE FROM tblQueue WHERE id = ?"
	db.statement, err = db.conn.Prepare(sq)
	if err != nil {
		return fmt.Errorf("DeletefromQueue Prepare error: %w", err)
	}
	defer db.statement.Close()

	db.reply, err = db.statement.Exec(id)
	if err != nil {
		return fmt.Errorf("DeletefromQueue Exec error: %w", err)
	}
	return nil
}

func (db *DBtype) ShowQueue() [][]string {
	var r [][]string
	var err error
	var phone, fname, lname string

	if err = db.Opendb(); err != nil {
		log.Println("#0 ShowQueue Opendb error:", err)
		return nil
	}

	// Använd LEFT JOIN så att meddelanden visas även om kunden inte finns i tblCustomers
	sq := "SELECT q.phone, COALESCE(c.firstname, ''), COALESCE(c.lastname, '') FROM tblQueue AS q LEFT JOIN tblCustomers AS c ON q.phone = c.phone ORDER BY q.tstamp ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		log.Println("#1 ShowQueue Query error:", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&phone, &fname, &lname)
		if err != nil {
			log.Println("#2 ShowQueue Scan error:", err)
			continue
		}
		r = append(r, []string{phone, fname, lname})
	}
	return r
}

func (db *DBtype) DeleteAllQueue() error {
	var err error

	if err = db.Opendb(); err != nil {
		return fmt.Errorf("#0 DeleteAllQueue Opendb error: %w", err)
	}

	sq := "DELETE FROM tblQueue;"
	db.statement, err = db.conn.Prepare(sq)
	if err != nil {
		return fmt.Errorf("#1 DeleteAllQueue Prepare DELETE: %w", err)
	}
	defer db.statement.Close()

	db.reply, err = db.statement.Exec()
	if err != nil {
		return fmt.Errorf("#2 DeleteAllQueue Exec: %w", err)
	}
	return nil
}

func (db *DBtype) ImportQueue(b string) {
	var err error

	if err = db.Opendb(); err != nil {
		log.Println("#0 ImportQueue Opendb error:", err)
		return
	}

	sq := "INSERT INTO tblQueue (tstamp, groupname, phone, message) VALUES (?, ?, ?, ?)"
	db.statement, err = db.conn.Prepare(sq)
	if err != nil {
		log.Println("#1 ImportQueue prepare error:", err)
		return
	}
	defer db.statement.Close()

	lines := strings.Split(b, "\r")
	for _, line := range lines {
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
			log.Println("#2 ImportQueue Exec error:", err)
			return
		}
	}
}

func (db *DBtype) ExportQueue(tofile string) error {
	var err error
	var txt, tstamp, groupname, phone, message string

	if err = db.Opendb(); err != nil {
		return fmt.Errorf("#0 ExportQueue Opendb error: %w", err)
	}

	sq := "SELECT tstamp, groupname, phone, message FROM tblQueue ORDER BY tstamp ASC"
	rows, err := db.conn.Query(sq)
	if err != nil {
		return fmt.Errorf("#1 ExportQueue Query error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&tstamp, &groupname, &phone, &message)
		if err != nil {
			return fmt.Errorf("#2 ExportQueue Scan error: %w", err)
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
		return fmt.Errorf("#3 ExportQueue WriteFile error: %w", err)
	}
	return nil
}

func (db *DBtype) CountinQueue() int {
	var err error
	var cnt int

	if err = db.Opendb(); err != nil {
		log.Println("#0 CountinQueue Opendb error:", err)
		return 0
	}

	err = db.conn.QueryRow("SELECT COUNT(*) FROM tblQueue").Scan(&cnt)
	if err != nil {
		log.Println("#1 CountQueue QueryRow failed:", err)
		return 0
	}

	return cnt
}

func (db *DBtype) GetNextinQueue() ([]string, error) {
	var err error
	var id, tstamp, groupname, phone, message string

	if err = db.Opendb(); err != nil {
		return nil, fmt.Errorf("#0 GetNextinQueue Opendb error: %w", err)
	}

	err = db.conn.QueryRow("SELECT id, tstamp, groupname, phone, message FROM tblQueue ORDER BY tstamp ASC LIMIT 1").Scan(&id, &tstamp, &groupname, &phone, &message)
	if err != nil {
		return nil, err
	}

	return []string{id, tstamp, groupname, phone, message}, nil
}
