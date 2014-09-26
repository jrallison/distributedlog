package db

import (
	"bufio"
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

	f, err := os.Open(logpath)
	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			println(err.Error())
			log.Fatal(err)
		}
	} else {
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			n += 1
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}

	f, err = os.OpenFile(logpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return &DB{
		dir:     path,
		file:    f,
		written: n,
	}
}

func (db *DB) Log(id string) (err error) {
	db.count += 1

	if db.count > db.written {
		if _, err = db.file.WriteString(id + "\n"); err == nil {
			println(id)
		}
	}

	return
}
