Usage:

    go run runner.go

It'll walk you through the rest of the setup.

## A quick overview:

The `runner.go` script will help you set up a cluster of three nodes. Then it will begin bombarding the cluster
with 500 threads, each attempting to write log entries to a random node in the cluster.  The cluster
will accept writes as long as a majority are running and can communicate with each other.

The cluster is responsible for making sure each node has the *exact same* list of events in it's log.
You can kill nodes, prevent them from communicating with each other, wreck havok, do you worst.  The end
result should be that two things:

1. All nodes have the exact same event history (consistency)
2. All events that a node in the cluster returned with a 200 response code, should be present in the log (no lost data)

## Your mission, should you choose to accept it:

Through any means necessary\*, cause the cluster to break one of the two promises above. `runner.go` will verify
both promises at the end of the run.

\* - you can trivally break it by deleting the data behind the nodes... let's limit this to networking
and liveness tricks (e.g. force killing nodes, pausing processes, etc, etc).

## Example output of `runner.go`:

    > go run runner.go
    start node in another console with the command "go run node.go -p 4000 tmp/0". Then press enter.
    start node in another console with the command "go run node.go -p 4001 -join=localhost:4000 tmp/1". Then press enter.
    start node in another console with the command "go run node.go -p 4002 -join=localhost:4000 tmp/2". Then press enter.

    slamming events into the cluster until you stop me by hitting enter.
      - feel free to kill nodes and see what happens.
      - to restart a node, run the same command you ran to start it, minus the '-join' parameter.
    1000 events
    2000 events
    3000 events
    4000 events
    5000 events
    6000 events
    7000 events
    waiting for events to finish...
    8000 events
    9000 events
    10000 events
    11000 events
    12000 events
    13000 events
    14000 events
    15000 events
    16000 events
    whew... that was rough, let's rest a bit and then verify results...
    COOL: found 16384 / 16384 acknowledged events in log
    SUCCESS: all nodes are consistent and all acknowledged events are present


