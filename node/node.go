package node

import (
	"github.com/goraft/raft"
	"github.com/jrallison/distributedlog/db"

	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/rpc"
	"path/filepath"
	"strings"
	"sync"
)

var clientMutex sync.Mutex

type Node struct {
	name       string
	host       string
	port       int
	path       string
	raftServer raft.Server
	clients    map[string]*rpc.Client
	db         *db.DB
}

func New(path, host string, port int) (n *Node) {
	n = &Node{
		host:    host,
		port:    port,
		path:    path,
		clients: make(map[string]*rpc.Client),
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

func (n *Node) client(node string) (*rpc.Client, error) {
	var err error

	if _, ok := n.clients[node]; !ok {
		clientMutex.Lock()

		if c, ok := n.clients[node]; !ok {
			println("creating client for", node)

			host := strings.Replace(node, "http://", "", 1)

			if c, err = rpc.DialHTTP("tcp", host); err == nil {
				n.clients[node] = c
			}
		}

		clientMutex.Unlock()
	}

	return n.clients[node], err
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
