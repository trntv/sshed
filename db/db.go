package db

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/boltdb/bolt"
	"time"
)

var servers = []byte("servers")
var meta = []byte("meta")

type Server struct {
	Host     string
	Port     string
	User     string
	Password string
	KeyFile  string
}

type DB struct {
	Path         string
	Password     string
	Bootstrapped bool

	boltdb *bolt.DB
}

type ErrNotFound struct {
	Key string
}

func (err ErrNotFound) Error() string {
	return fmt.Sprintf("Server with key %s not found", err.Key)
}

func NewDB(path string) (*DB, error) {
	database := &DB{
		Path:         path,
		Bootstrapped: false,
	}

	err := database.Open()
	if err != nil {
		return nil, err
	}
	defer database.Close()

	err = database.boltdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(servers)
		if err != nil {
			if err == bolt.ErrBucketExists {
				database.Bootstrapped = true
			} else {
				return err
			}
		}

		if database.Bootstrapped == false {
			_, err = tx.CreateBucketIfNotExists(meta)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return database, err
}

func (DB *DB) Open() error {
	var err error
	DB.boltdb, err = bolt.Open(DB.Path, 0600, &bolt.Options{Timeout: 1 * time.Second})

	return err
}

func (DB *DB) MustOpen() {
	if DB.boltdb != nil {
		return
	}

	err := DB.Open()
	if err != nil {
		panic(err)
	}
}

func (DB *DB) Close() error {
	if DB.boltdb == nil {
		return nil
	}

	err := DB.boltdb.Close()
	if err != nil {
		return err
	}

	DB.boltdb = nil

	return nil
}

func (DB *DB) Get(key string) (srv *Server, err error) {
	DB.MustOpen()

	srv = &Server{}
	err = DB.boltdb.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(servers)
		value := b.Get([]byte(key))
		if value == nil {
			return ErrNotFound{Key: key}
		}

		err = DB.unmarshal(value, srv)

		return err
	})

	return srv, err
}

func (DB *DB) Put(key string, server *Server) error {
	DB.MustOpen()

	err := DB.boltdb.Update(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(servers)
		data, err := DB.marshal(server)
		if err != nil {
			return err
		}
		err = b.Put([]byte(key), data)

		return err
	})

	return err
}

func (DB *DB) Remove(key string) error {
	DB.MustOpen()

	err := DB.boltdb.Update(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(servers)
		err := b.Delete([]byte(key))

		return err
	})

	return err
}

func (DB *DB) GetAll() (srvs map[string]*Server, err error) {
	DB.MustOpen()

	srvs = make(map[string]*Server)
	err = DB.boltdb.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(servers)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var srv = &Server{}
			err = DB.unmarshal(v, srv)
			if err != nil {
				return err
			}
			srvs[string(k)] = srv
		}

		return nil
	})

	return srvs, err
}

func (DB *DB) marshal(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	data := b.Bytes()

	isEncrypted, err := DB.IsEncrypted()
	if err != nil {
		return nil, err
	}
	if isEncrypted == true {
		data, err = DB.encrypt(data, DB.Password)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (DB *DB) unmarshal(b []byte, v interface{}) error {
	isEncrypted, err := DB.IsEncrypted()
	if err != nil {
		return err
	}
	if isEncrypted == true {
		b, err = DB.decrypt(b, DB.Password)
		if err != nil {
			return err
		}
	}

	r := bytes.NewReader(b)
	dec := gob.NewDecoder(r)
	return dec.Decode(v)
}
