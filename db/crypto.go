package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"io"
)

var metaKey = []byte("encrypted")

func (DB *DB) EncryptDatabase(key string) error {
	DB.MustOpen()

	records, err := DB.GetAll()
	if err != nil {
		return err
	}

	DB.Password = key

	err = DB.boltdb.Update(func(tx *bolt.Tx) error {
		ciphertext, err := DB.encrypt([]byte("plaintext"), key)
		if err != nil {
			return err
		}

		b := tx.Bucket(meta)
		err = b.Put(metaKey, ciphertext)

		return err
	})

	DB.Close()

	for key, value := range records {
		err := DB.Put(key, value)
		if err != nil {
			return err
		}
	}

	return err
}

func (DB *DB) IsEncrypted() (bool, error) {
	var isEncrypted bool
	DB.MustOpen()

	err := DB.boltdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(meta)
		v := b.Get(metaKey)
		if v != nil && string(v) != "plaintext" {
			isEncrypted = true
		}
		return nil
	})

	return isEncrypted, err
}

func (DB *DB) TestKey(key string) (bool, error) {
	var isValid bool
	DB.MustOpen()

	err := DB.boltdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(meta)
		v := b.Get(metaKey)
		if v == nil {
			return nil
		}

		p, err := DB.decrypt(v, key)
		if err != nil {
			return err
		}

		if string(p) != "plaintext" {
			isValid = true
		}

		return nil
	})

	return isValid, err
}

func (DB *DB) encrypt(plaintext []byte, password string) ([]byte, error) {
	cKey, err := DB.makeCipherKey([]byte(password))
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(cKey)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func (DB *DB) decrypt(ciphertext []byte, password string) ([]byte, error) {
	cKey, err := DB.makeCipherKey([]byte(password))
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(cKey)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New(fmt.Sprintf("ciphertext too short: %d", len(ciphertext)))
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)

	plaintext := make([]byte, len(ciphertext))
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

func (DB *DB) makeCipherKey(data []byte) ([]byte, error) {
	hash := md5.New()
	_, err := hash.Write(data)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}
