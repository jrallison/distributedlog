package node

import (
	"github.com/goraft/raft"

	"errors"
	"net/http"
	"net/rpc"
	"strings"
)

func (n *Node) Leave() error {
	leaderName := n.raftServer.Leader()

	command := raft.DefaultLeaveCommand{
		Name: n.raftServer.Name(),
	}

	if n.raftServer.Name() == leaderName {
		_, err := n.raftServer.Do(&command)
		return err
	} else if leader, ok := n.raftServer.Peers()[leaderName]; ok {
		host := strings.Replace(leader.ConnectionString, "http://", "", 1)

		client, err := rpc.DialHTTP("tcp", host)
		if err != nil {
			return err
		}

		return client.Call("Node.Leave", command, &Nothing{})
	} else {
		return errors.New("No current leader?: " + leaderName)
	}
}

func (n *Node) leaveHandler(w http.ResponseWriter, req *http.Request) {
	if err := n.Leave(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
