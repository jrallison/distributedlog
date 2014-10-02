package db

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type DB struct {
	dir     string
	count   int
	file    *os.File
	written int
}

func New(path string) *DB {
	logpath := filepath.Join(path, "log")

	var n int

	err := os.Remove(logpath)
	if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
		log.Fatal(err)
	}

	f, err := os.OpenFile(logpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	println("read", n, "events from log")

	return &DB{
		dir:     path,
		file:    f,
		written: n,
	}
}

func (db *DB) Log(id string) (err error) {
	db.count += 1

	println(db.count, db.written)
	if db.count > db.written {
		println("acked", id)
		_, err = db.file.WriteString(id + "\n")
		if err != nil {
			panic(err)
		}
	}

	return
}
