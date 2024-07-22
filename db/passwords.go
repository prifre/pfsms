package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
)
func (db *DBtype) SetHash(hash string) error {
    var err error
	var sq string
	var r *sql.Rows
    var h string
	db.Opendb()
	// fmt.Println("nanostamp=",nanostamp,"tstamp=",tstamp)	
	sq="SELECT hash FROM tblHashtable"
	r, err = db.conn.Query(sq)
	if r != nil && err==nil {
		r.Next()
		r.Scan(&h)
        if h>"" {
		    return errors.New("#1 SetHash hash already exists!!! cannot overwrite")
        }
	}
	sq="INSERT INTO tblHashtable (hash) "
	sq = fmt.Sprintf("%s VALUES ('%s')",sq,hash)
	db.statement, err = db.conn.Prepare(sq) // Prepare SQL Statement
	if err != nil {
		log.Println("#2 prepare failed: ", sq, " ", err.Error())
		return err
	}
	db.reply, err = db.statement.Exec() // Execute SQL Statements
    return err
}
func (db *DBtype) GetHash() (string,error) {
    var err error
	var sq string
	var r *sql.Rows
    var h string
	db.Opendb()
	sq="SELECT hash FROM tblHashtable"
	r, err = db.conn.Query(sq)
	if r != nil && err==nil {
		r.Next()
		r.Scan(&h)
        if h>"" {
            db.hash=h
    		return h,err
        }
	}
    db.hash=""
    return "",err
}
func (db *DBtype) EncryptPassword(password string) (string,error) {
	var err error
	var ciphertext,key,plaintext []byte
	if len(db.hash)<32 {
		db.hash=fmt.Sprintf("%sabcdefghijklmnopqrstuvwxyz1234567890",db.hash)[0:32]
	}
    key = []byte(string(db.hash)) // 32 bytes
    plaintext = []byte(password)
    ciphertext, err = encrypt(key, plaintext)
	return string(ciphertext),err
}
func (db *DBtype) DecryptPassword(password string) (string,error) {
	var err error
    var ciphertext,key,result []byte
	ciphertext = []byte(password)
	fmt.Printf("%0x\n", ciphertext)
    key=[]byte(db.hash)
	result, err = decrypt(key, ciphertext)
    return string( result),err
}

// See alternate IV creation from ciphertext below
//var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func encrypt(key, text []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    b := base64.StdEncoding.EncodeToString(text)
    ciphertext := make([]byte, aes.BlockSize+len(b))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }
    cfb := cipher.NewCFBEncrypter(block, iv)
    cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
    return ciphertext, nil
}

func decrypt(key, text []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    if len(text) < aes.BlockSize {
        return nil, errors.New("ciphertext too short")
    }
    iv := text[:aes.BlockSize]
    text = text[aes.BlockSize:]
    cfb := cipher.NewCFBDecrypter(block, iv)
    cfb.XORKeyStream(text, text)
    data, err := base64.StdEncoding.DecodeString(string(text))
    if err != nil {
        return nil, err
    }
    return data, nil
}