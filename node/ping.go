package node

import (
	"github.com/goraft/raft"

	"net/rpc"
	"strings"
)

func Ping(peer *raft.Peer) (*rpc.Call, error) {
	host := strings.Replace(peer.ConnectionString, "http://", "", 1)

	client, err := rpc.DialHTTP("tcp", host)
	if err != nil {
		return nil, err
	}

	state := new(State)

	call := client.Go("Node.State", Nothing{}, state, nil)

	return call, nil
}
