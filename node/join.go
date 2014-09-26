package node

import (
	"github.com/goraft/raft"

	"net/rpc"
)

func (n *Node) Join(node string) error {
	client, err := rpc.DialHTTP("tcp", node)
	if err != nil {
		return err
	}

	command := raft.DefaultJoinCommand{
		Name:             n.raftServer.Name(),
		ConnectionString: n.uri(),
	}

	return client.Call("Node.Join", command, &Nothing{})
}
