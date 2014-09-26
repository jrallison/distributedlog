package node

import (
	"github.com/goraft/raft"
	"github.com/jrallison/distributedlog/db"

	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func (n *Node) ConnectCluster(connect string) (err error) {
	raft.RegisterCommand(&LogCommand{})

	transporter := raft.NewHTTPTransporter("/raft", 200*time.Millisecond)

	if err := os.MkdirAll(n.path+"/db", 0744); err != nil {
		log.Fatalf("Unable to create db directory: %v", err)
	}

	n.db = db.New(n.path + "/db")

	n.raftServer, err = raft.NewServer(n.name, n.path, transporter, nil, n.db, "")
	if err != nil {
		return
	}

	transporter.Install(n.raftServer, n)
	err = n.raftServer.Start()
	if err != nil {
		return
	}

	if connect != "" {
		log.Println("Attempting to join cluster:", connect)

		if !n.raftServer.IsLogEmpty() {
			log.Fatal("Cannot join with an existing log")
		}

		if err = n.Join(connect); err != nil {
			return
		}
	} else if n.raftServer.IsLogEmpty() {
		log.Println("Initializing new cluster")

		_, err = n.raftServer.Do(&raft.DefaultJoinCommand{
			Name:             n.raftServer.Name(),
			ConnectionString: n.uri(),
		})

		if err != nil {
			return
		}
	} else {
		log.Println("Recovered from log")
	}

	return
}

func (n *Node) clusterHandler(w http.ResponseWriter, req *http.Request) {
	body := make(map[string]interface{})

	if n.raftServer.Running() {
		reachablePeers := 1 // local node

		statuses := make(map[string]interface{})

		for name, peer := range n.raftServer.Peers() {
			call, err := Ping(peer)

			if err == nil {
				reachablePeers += 1

				select {
				case <-call.Done:
					if call.Error == nil {
						statuses[name] = call.Reply
					} else {
						err = call.Error
					}
				case <-time.After(100 * time.Millisecond):
					err = errors.New("timeout")
				}
			}

			if err != nil {
				statuses[name] = "error: " + err.Error()
			}
		}

		statuses[n.raftServer.Name()] = n.State()

		body["_self"] = n.raftServer.Name()

		status := "available"

		if reachablePeers < n.raftServer.QuorumSize() {
			status = "unavailable"
		}

		body["cluster"] = map[string]interface{}{
			"term":   n.raftServer.Term(),
			"status": fmt.Sprint(status, " (", reachablePeers, "/", n.raftServer.MemberCount(), " nodes reachable)"),
			"nodes":  statuses,
		}
	} else {
		body["_self"] = "not connected"
	}

	js, _ := json.MarshalIndent(body, "", "  ")
	w.Write(js)
}
