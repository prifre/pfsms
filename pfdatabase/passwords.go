package pfdatabase

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
func MakeHash() (string,error) {
    var err error
	var sq string
	var rows *sql.Rows
    var h,hash string
	db := new(DBtype)
	db.Opendb()
	sq="SELECT hash FROM tblHashtable"
	rows, err = db.conn.Query(sq)
	if err==nil {
		for rows.Next() {
			err = rows.Scan(&hash)
			if err!=nil {
				log.Println("#1 MakeHash Error Scan)",err.Error())
			}
			if h=="" && hash>"" {
				h=hash
			}
		}
	}
	db.Closedatabase()
	if h>"" {
	    return h,nil
	}
	b := make([]rune, 32)
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := range b {
		b[i] = letters[mathrand.Intn(len(letters))]
	}
	h = string(b)
	sq="INSERT INTO tblHashtable (hash) "
	sq = fmt.Sprintf("%s VALUES ('%s')",sq,h)
	db.Opendb()
	db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
	if err != nil {
		log.Println("#2 MakeHash Prepare failed: ", sq, " ", err.Error())
		db.Closedatabase()
		return "",err
	}
	db.reply, err = db.statement.Exec() // Execute SQL Statements
	if err != nil {
		log.Println("#3 MakeHash Exec failed: ", sq, " ", err.Error())
		db.Closedatabase()
		return "",err
	}
	db.Closedatabase()
    return h,err
}
func EncryptPassword(text string,hash string) (string, error) {
	block, err := aes.NewCipher([]byte(hash))
	if err != nil {
	 return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}
func DecryptPassword(text string,hash string) (string, error) {
	block, err := aes.NewCipher([]byte(hash))
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