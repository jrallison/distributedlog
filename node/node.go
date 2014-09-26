package node

import (
	"github.com/goraft/raft"
	"github.com/jrallison/distributedlog/db"

	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"path/filepath"
)

type Node struct {
	name       string
	host       string
	port       int
	path       string
	raftServer raft.Server
	db         *db.DB
}

func New(path, host string, port int) (n *Node) {
	n = &Node{
		host: host,
		port: port,
		path: path,
	}

	// Read existing name or generate a new one.
	if b, err := ioutil.ReadFile(filepath.Join(path, "name")); err == nil {
		n.name = string(b)
	} else {
		n.name = fmt.Sprintf("%07x", rand.Int())[0:7]
		if err = ioutil.WriteFile(filepath.Join(path, "name"), []byte(n.name), 0644); err != nil {
			panic(err)
		}
	}

	return
}

func (n *Node) uri() string {
	return fmt.Sprintf("http://%s:%d", n.host, n.port)
}

func (n *Node) Start(join string) (err error) {
	log.Printf("Initializing Raft Server: %s", n.path)

	if err = n.ConnectCluster(join); err != nil {
		log.Fatal(err)
	}

	log.Println("Initializing HTTP server")

	return n.Serve()
}
