package node

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"
)

func (n *Node) Serve() error {
	rpc.RegisterName("Node", &rpcWrapper{n})
	rpc.HandleHTTP()

	n.HandleFunc("/cluster", n.clusterHandler)
	n.HandleFunc("/leave", n.leaveHandler)
	n.HandleFunc("/log/", n.logHandler)

	log.Println("Listening at:", n.uri())

	return http.ListenAndServe(fmt.Sprintf(":%d", n.port), nil)
}

func (n *Node) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, handler)
}
