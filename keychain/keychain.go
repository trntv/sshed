package keychain

import (
	"bytes"
	"encoding/gob"
	"github.com/boltdb/bolt"
	"time"
)

var Bootstrapped bool

var keys = []byte("keys")
var meta = []byte("meta")

var path string
var Password string
var Encrypted bool

var boltdb *bolt.DB

type Record struct {
	Password   string
	PrivateKey string
}

func Open(p string) error {
	path = p

	err := open()
	if err != nil {
		return err
	}
	defer close()

	err = boltdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(keys)
		if err != nil {
			if err == bolt.ErrBucketExists {
				Bootstrapped = true
			} else {
				return err
			}
		}

		if Bootstrapped == false {
			_, err = tx.CreateBucketIfNotExists(meta)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = boltdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(meta)
		v := b.Get(metaKey)
		if v != nil && string(v) != "plaintext" {
			Encrypted = true
		}
		return nil
	})

	return err
}

func Get(host string) (k *Record, err error) {
	mustOpen()
	defer close()

	k = &Record{}
	err = boltdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(keys)
		value := b.Get([]byte(host))
		if value == nil {
			return ErrNotFound{Key: host}
		}

		err = unmarshal(value, k)

		return err
	})

	return k, err
}

func Put(key string, record *Record) error {
	mustOpen()
	defer close()

	err := boltdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(keys)
		data, err := marshal(record)
		if err != nil {
			return err
		}
		err = b.Put([]byte(key), data)

		return err
	})

	return err
}

func Remove(key string) error {
	mustOpen()
	defer close()

	err := boltdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(keys)
		err := b.Delete([]byte(key))

		return err
	})

	return err
}

func list() (records map[string]*Record, err error) {
	mustOpen()
	defer boltdb.Close()

	records = make(map[string]*Record)
	err = boltdb.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(keys)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var srv = &Record{}
			err = unmarshal(v, srv)
			if err != nil {
				return err
			}
			records[string(k)] = srv
		}

		return nil
	})

	return records, err
}

func close() error {
	if boltdb == nil {
		return nil
	}

	err := boltdb.Close()
	if err != nil {
		return err
	}

	boltdb = nil

	return nil
}

func open() error {
	var err error
	boltdb, err = bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})

	return err
}

func mustOpen() {
	if boltdb == nil {
		err := open()
		if err != nil {
			panic(err)
		}
	}
}

func marshal(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	data := b.Bytes()

	if Encrypted == true {
		data, err = encrypt(data, Password)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func unmarshal(b []byte, v interface{}) error {
	var err error
	if Encrypted == true {
		b, err = decrypt(b, Password)
		if err != nil {
			return err
		}
	}

	r := bytes.NewReader(b)
	dec := gob.NewDecoder(r)
	return dec.Decode(v)
}
