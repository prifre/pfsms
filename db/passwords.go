package db

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	mathrand "math/rand"
)
var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}
func (db *DBtype) MakeHash() error {
    var err error
	var sq string
	var r *sql.Rows
    var h string
	// fmt.Println("nanostamp=",nanostamp,"tstamp=",tstamp)	
	sq="SELECT hash FROM tblHashtable"
	r, err = db.conn.Query(sq)
	if r != nil && err==nil {
		r.Next()
		r.Scan(&h)
        if h>"" {
			db.hash=h
		    return nil
        }
	}
	b := make([]rune, 32)
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := range b {
		b[i] = letters[mathrand.Intn(len(letters))]
	}
	h = string(b)
	sq="INSERT INTO tblHashtable (hash) "
	sq = fmt.Sprintf("%s VALUES ('%s')",sq,h)
	db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
	if err != nil {
		log.Println("#2 prepare failed: ", sq, " ", err.Error())
		return err
	}
	db.reply, err = db.statement.Exec() // Execute SQL Statements
	db.hash=h
    return err
}
func (db *DBtype) EncryptPassword(text string) (string, error) {
	block, err := aes.NewCipher([]byte(db.hash))
	if err != nil {
	 return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}
func (db *DBtype) DecryptPassword(text string) (string, error) {
	block, err := aes.NewCipher([]byte(db.hash))
	if err != nil {
		return "", err
	}
	b,_:=base64.StdEncoding.DecodeString(text)
	cipherText :=b
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
