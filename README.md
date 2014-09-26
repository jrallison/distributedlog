Usage:

    go run runner.go

It'll walk you through the rest of the setup.

## A quick overview:

The `runner.go` script will help you set up a cluster of three nodes. Then it will begin bombarding the cluster
with 500 threads, each attempting to write log entries to a random node in the cluster.

The cluster is responsible for making sure each node has the *exact same* list of requests in it's log.
You can kill nodes, prevent them from communicating with each other, wreck havok, do you worst.  The end
result should be that two things:

1) All nodes have the exact same request log history (consistency)
2) All requests that a node in the cluster returned with a 200 response code, should be present in the log (no lost data)

## Your mission, should you choose to accept it:

Through any means necessary, cause the cluster to break one of the two promises above. `runner.go` will verify
both promises at the end of the run.

## Example output of `runner.go`:

    > go run runner.go
    start node in another console with the command "go run node.go -p 4000 tmp/0". Then press enter.
    start node in another console with the command "go run node.go -p 4001 -join=localhost:4000 tmp/1". Then press enter.
    start node in another console with the command "go run node.go -p 4002 -join=localhost:4000 tmp/2". Then press enter.

    slamming requests into the cluster until you stop me by hitting enter.
      - feel free to kill nodes and see what happens.
      - to restart a node, run the same command you ran to start it, minus the '-join' parameter.
    1000 requests
    2000 requests
    3000 requests
    4000 requests
    5000 requests
    6000 requests
    7000 requests
    waiting for requests to finish...
    8000 requests
    9000 requests
    10000 requests
    11000 requests
    12000 requests
    13000 requests
    14000 requests
    15000 requests
    16000 requests
    whew... that was rough, let's rest a bit and then verify results...
    COOL: found 16384 / 16384 acknowledged requests in log
    SUCCESS: all nodes are consistent and all acknowledged requests are present


