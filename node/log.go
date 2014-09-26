package node

import (
	"github.com/goraft/raft"
	"github.com/jrallison/distributedlog/db"

	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"
	"strings"
)

type LogCommand struct {
	Id string `json:"id"`
}

func NewLogCommand(id string) *LogCommand {
	return &LogCommand{
		Id: id,
	}
}

func (c *LogCommand) CommandName() string {
	return "log"
}

func (c *LogCommand) Apply(server raft.Server) (interface{}, error) {
	db := server.Context().(*db.DB)
	return "", db.Log(c.Id)
}

func (n *Node) Log(id string) error {
	leaderName := n.raftServer.Leader()

	command := NewLogCommand(id)

	var err error

	if n.raftServer.Name() == leaderName {
		_, err = n.raftServer.Do(command)
	} else if leader, ok := n.raftServer.Peers()[leaderName]; ok {
		host := strings.Replace(leader.ConnectionString, "http://", "", 1)

		client, e := rpc.DialHTTP("tcp", host)
		if e != nil {
			err = e
		}

		if err == nil {
			result := new(String)
			err = client.Call("Node.Log", command, &result)
		}
	} else {
		err = errors.New("No current leader?: " + leaderName)
	}

	return err
}

func (n *Node) logHandler(w http.ResponseWriter, req *http.Request) {
	id := strings.Replace(req.URL.Path, "/log/", "", 1)

	err := n.Log(id)

	var body []byte

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		body, _ = json.MarshalIndent(map[string]string{
			"info": err.Error(),
		}, "", "  ")
	} else {
		body, _ = json.MarshalIndent(map[string]string{}, "", "  ")
	}

	w.Write(append(body, []byte("\n")...))
}
