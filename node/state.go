package node

func (n *Node) State() State {
	return State{
		n.raftServer.Name(),
		n.raftServer.State(),
		n.raftServer.CommitIndex(),
		n.path,
		n.uri(),
	}
}
