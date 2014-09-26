package node

import (
	"github.com/goraft/raft"

	"errors"
	"net/rpc"
	"strings"
)

type rpcWrapper struct {
	node *Node
}

type State struct {
	Name  string `json:"name"`
	State string `json:"state"`
	Path  string `json:"path"`
	Uri   string `json:"uri"`
}

type Nothing struct {
}

func (w *rpcWrapper) State(args Nothing, reply *State) error {
	*reply = w.node.State()
	return nil
}

func (w *rpcWrapper) Join(command raft.DefaultJoinCommand, reply *Nothing) (err error) {
	leaderName := w.node.raftServer.Leader()

	if w.node.raftServer.Name() == leaderName {
		_, err = w.node.raftServer.Do(&command)
	} else if leader, ok := w.node.raftServer.Peers()[leaderName]; ok {
		host := strings.Replace(leader.ConnectionString, "http://", "", 1)

		client, err := rpc.DialHTTP("tcp", host)
		if err != nil {
			return err
		}

		err = client.Call("Node.Join", command, &Nothing{})
		client.Close()
	} else {
		err = errors.New("No current leader?: " + leaderName)
	}

	return
}

func (w *rpcWrapper) Leave(command raft.DefaultLeaveCommand, reply *Nothing) (err error) {
	_, err = w.node.raftServer.Do(&command)
	return
}

func (w *rpcWrapper) Log(command LogCommand, reply *Nothing) (err error) {
	_, err = w.node.raftServer.Do(&command)
	return
}
