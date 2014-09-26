package node

func (n *Node) State() State {
	return State{
		n.raftServer.Name(),
		n.raftServer.State(),
		n.path,
		n.uri(),
	}
}
