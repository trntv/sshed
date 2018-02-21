package keychain

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

func EncryptDatabase(password string) error {
	mustOpen()
	defer close()

	records, err := list()
	if err != nil {
		return err
	}

	Password = password

	err = boltdb.Update(func(tx *bolt.Tx) error {
		ciphertext, err := encrypt([]byte("plaintext"), password)
		if err != nil {
			return err
		}

		b := tx.Bucket(meta)
		err = b.Put(metaKey, ciphertext)

		return err
	})

	for key, value := range records {
		err := Put(key, value)
		if err != nil {
			return err
		}
	}

	return err
}

func encrypt(plaintext []byte, password string) ([]byte, error) {
	cKey, err := makeCipherKey([]byte(password))
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

func decrypt(ciphertext []byte, password string) ([]byte, error) {
	cKey, err := makeCipherKey([]byte(password))
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

func makeCipherKey(data []byte) ([]byte, error) {
	hash := md5.New()
	_, err := hash.Write(data)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}
