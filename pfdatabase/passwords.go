package pfdatabase

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
)

// Fasta initialiseringsvektorn (IV) för CFB
var ivBytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func MakeHash() (string, error) {
	var err error
	var hash string

	db := new(DBtype)
	if err = db.Opendb(); err != nil {
		return "", fmt.Errorf("#0 MakeHash Opendb: %w", err)
	}

	sq := "SELECT hash FROM tblHashtable LIMIT 1"
	rows, err := db.conn.Query(sq)
	if err == nil {
		if rows.Next() {
			err = rows.Scan(&hash)
			if err != nil {
				log.Println("#1 MakeHash Scan error:", err.Error())
			}
		}
		rows.Close() // 👈 Viktigt! Stäng raderna direkt
	}

	if hash != "" {
		return hash, nil
	}

	// Skapa en kryptografiskt säker slumpmässig nyckel (32 tecken/bytes)
	randomBytes := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
		return "", fmt.Errorf("#2 MakeHash crypto/rand failed: %w", err)
	}
	// Omvandla till 32 hex-tecken för att matcha 256-bitars AES-nyckel
	hash = fmt.Sprintf("%x", randomBytes)

	// Spara i databasen med parametriserad fråga
	sqInsert := "INSERT INTO tblHashtable (hash) VALUES (?)"
	db.statement, err = db.conn.Prepare(sqInsert)
	if err != nil {
		log.Println("#3 MakeHash Prepare failed:", err.Error())
		return "", err
	}
	defer db.statement.Close()

	db.reply, err = db.statement.Exec(hash)
	if err != nil {
		log.Println("#4 MakeHash Exec failed:", err.Error())
		return "", err
	}

	return hash, nil
}

func EncryptPassword(text string, hash string) (string, error) {
	if len(hash) == 0 {
		return "", fmt.Errorf("hash key cannot be empty")
	}

	block, err := aes.NewCipher([]byte(hash))
	if err != nil {
		return "", err
	}

	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, ivBytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptPassword(text string, hash string) (string, error) {
	if len(hash) == 0 {
		return "", fmt.Errorf("hash key cannot be empty")
	}

	block, err := aes.NewCipher([]byte(hash))
	if err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %w", err)
	}

	cfb := cipher.NewCFBDecrypter(block, ivBytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}
