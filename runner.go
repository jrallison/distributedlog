package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var CONCURRENCY = 500

var NODES = []string{
	"localhost:4000",
	"localhost:4001",
	"localhost:4002",
}

func main() {
	checkNodesArentRunning()
	os.RemoveAll("tmp")
	startNodes()

	requests := make(map[string]bool)
	stop := make(chan bool)
	exited := make(chan bool)

	go slam(requests, stop, exited)

	fmt.Println("slamming events into the cluster until you stop me by hitting enter.")
	fmt.Println(" - feel free to kill nodes and see what happens.")
	fmt.Println(" - to restart a node, run the same command you ran to start it, minus the '-join' parameter.")
	fmt.Scanln()

	stop <- true
	fmt.Println("waiting for requests to finish...")
	<-exited

	checkNodesAreRunning()

	fmt.Println("whew... that was rough, let's rest a bit and then verify results...")
	time.Sleep(10 * time.Second)

	checkNodesAreRunning()

	consistent := verifyConsistentLogs()
	requests0 := verifyRequests(requests, []string{"tmp/0/db/log"})
	requests1 := verifyRequests(requests, []string{"tmp/1/db/log"})
	requests2 := verifyRequests(requests, []string{"tmp/2/db/log"})
	allrequests := verifyRequests(requests, []string{"tmp/0/db/log", "tmp/1/db/log", "tmp/2/db/log"})

	if consistent && requests0 && requests1 && requests2 && allrequests {
		fmt.Println("SUCCESS: all nodes are consistent and all acknowledged events are present")
	} else {
		fmt.Println("WHOOPSIE: send the output of all consoles to John :(")
	}
}

func verifyConsistentLogs() bool {
	var log string
	var fails int

	for i, node := range NODES {
		data, err := ioutil.ReadFile("tmp/" + strconv.Itoa(i) + "/db/log")
		if err != nil {
			fmt.Println("FAIL: error reading log:", err)
			fails += 1
		}

		if i == 0 {
			log = string(data)
		} else {
			if string(data) != log {
				fmt.Println("FAIL:", NODES[0], "and", node, "have inconsistent logs!!!!")
				fails += 1
			}
		}
	}

	return fails == 0
}

func verifyRequests(requests map[string]bool, paths []string) bool {
	var fails int

	found := make(map[string]bool)

	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			fmt.Println("FAIL: error reading log:", err)
			fails += 1
		}

		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			found[scanner.Text()] = true
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("FAIL: error reading log:", err)
			fails += 1
		}
	}

	var ackedpresent int
	var acknowledged int
	var unackedpresent int
	var unacknowledged int

	for id, acked := range requests {
		if _, ok := found[id]; ok {
			if acked {
				ackedpresent += 1
				acknowledged += 1
			} else {
				unackedpresent += 1
				unacknowledged += 1
			}
		} else {
			if acked {
				acknowledged += 1
			} else {
				unacknowledged += 1
			}
		}

		delete(found, id)
	}

	if len(found) > 0 {
		fmt.Println("FAIL: non existant requests in log:", found)
		fails += 1
	}

	if ackedpresent != acknowledged {
		fmt.Println("FAIL: only found", ackedpresent, "/", acknowledged, "acknowledged events in logs", paths)
		fails += 1
	} else {
		fmt.Println("COOL: found", ackedpresent, "/", acknowledged, "acknowledged events in logs", paths)
	}

	if unackedpresent > 0 {
		fmt.Println("COOL: found", unackedpresent, "/", unacknowledged, "unacknowledged events in logs", paths)
	}

	return fails == 0
}

func slam(requests map[string]bool, stop chan bool, exited chan bool) {
	var mutex sync.Mutex
	var total int

	quit := make(chan bool)
	exit := make(chan bool)

	for i := 0; i < CONCURRENCY; i++ {
		go (func(index int) {
			var count int

			for {
				count += 1

				mutex.Lock()
				total += 1
				mutex.Unlock()

				if total%1000 == 0 {
					fmt.Println(total, "events")
				}

				select {
				case <-quit:
					exit <- true
					return
				default:
					id := fmt.Sprint(index, "-", count)

					node := NODES[rand.Intn(len(NODES))]

					resp, err := http.Get("http://" + node + "/log/" + id)
					if err == nil {
						io.Copy(ioutil.Discard, resp.Body)
						resp.Body.Close()
					}

					mutex.Lock()
					if err == nil && resp.StatusCode == 200 {
						requests[id] = true
					} else {
						requests[id] = false
					}
					mutex.Unlock()
				}
			}
		})(i)
	}

	<-stop

	for i := 0; i < CONCURRENCY; i++ {
		quit <- true
	}

	for i := 0; i < CONCURRENCY; i++ {
		<-exit
	}

	exited <- true
}

func checkNodesAreRunning() {
	for _, node := range NODES {
		for !ping(node) {
			fmt.Println("For verification, let's bring back up the node listening on http://" + node + " then press enter.")
			fmt.Scanln()
		}
	}
}

func checkNodesArentRunning() {
	var notified int

	for _, node := range NODES {
		for ping(node) {
			if notified == 0 {
				fmt.Println("Let's start from scratch!")
				notified += 1
			}

			fmt.Println("shut down the process running on http://" + node + " then press enter.")
			fmt.Scanln()
		}
	}
}

func startNodes() {
	join := " -join=" + NODES[0]

	for i, node := range NODES {
		port := strings.SplitN(node, ":", 2)[1]

		options := join

		if i == 0 {
			options = ""
		}

		for !ping(node) {
			fmt.Println("start node in another console with the command \"go run node.go -p " + port + options + " tmp/" + strconv.Itoa(i) + "\". Then press enter.")
			fmt.Scanln()
		}
	}
}

func ping(node string) bool {
	_, err := http.Get("http://" + node)
	return err == nil
}
