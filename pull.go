package main

import (
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	check("scp")
	check("sshpass")

	numConns := 8
	numBeacons := 1
	numShards := 2
	nodesPerShard := 4

	p := &Pool{
		Jobs: make(chan string, numConns),
		Done: make(chan struct{}, numConns),
		Fin:  make(chan struct{}),
	}

	for i := 0; i < numConns; i++ {
		go waitAndPull(p, i)
	}

	for i := 0; i < numBeacons; i++ {
		go pullBeacon(i, p)
	}

	for i := 0; i < numShards; i++ {
		for j := 0; j < nodesPerShard; j++ {
			go pullShard(i, j, p)
		}
	}

	for i := 0; i < numBeacons+numShards*nodesPerShard; i++ {
		<-p.Done
	}
	close(p.Fin)
	for i := 0; i < numConns; i++ {
		<-p.Done
	}
}

type Pool struct {
	Jobs chan string
	Done chan struct{}
	Fin  chan struct{}
}

func waitAndPull(p *Pool, i int) {
	for {
		select {
		case node := <-p.Jobs:
			log.Printf("Worker %d received node %s\n", i, node)
			pull(node)
			p.Done <- struct{}{}

		case _, ok := <-p.Fin:
			if !ok {
				log.Println("End worker", i)
				p.Done <- struct{}{}
				return
			}
		}
	}
}

func pullBeacon(i int, p *Pool) {
	node := strings.Join([]string{"beacon", strconv.Itoa(i)}, "")
	p.Jobs <- node
}

func pullShard(i, j int, p *Pool) {
	node := strings.Join([]string{"shard", strconv.Itoa(i), "-", strconv.Itoa(j)}, "")
	p.Jobs <- node
}

func pull(node string) {
	remoteFile, localFile := getPaths(node)

	password := ""
	cmd := exec.Command("sshpass", "-p", password, "scp", remoteFile, localFile)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(node, err)
		return
	}
	if err := cmd.Start(); err != nil {
		log.Println(node, err)
		return
	}
	slurp, _ := ioutil.ReadAll(stderr)
	if len(slurp) > 0 {
		log.Printf("%s stderr: %s", node, string(slurp))
	}
	log.Println(node, "done")
}

func getPaths(node string) (string, string) {
	ip := ""
	host := strings.Join([]string{"root@", ip, ":/data/"}, "")
	file := "log.txt"
	remoteFile := filepath.Join(host, node, file)

	saveDir := "/tmp"
	localFile := filepath.Join(saveDir, strings.Join([]string{node, "_", file}, ""))
	log.Println("Pulling", remoteFile)
	return remoteFile, localFile
}

func check(cmd string) {
	_, err := exec.LookPath(cmd)
	if err != nil {
		log.Fatal(err)
	}
}
