package main

import (
	cluster "github.com/abhijat/simple_tcp_cluster"
	"time"
	"os"
	"fmt"
)

var (
	peers = []string{
		":8080",
		":8081",
		":8082",
	}
)

func main() {

	args := os.Args
	if len(args) != 2 {
		fmt.Println("need a port")
		os.Exit(1)
	}

	listenAddress := args[1]
	go cluster.ServerLoop(listenAddress)

	for {
		for _, peer := range peers {
			if peer != listenAddress {
				cluster.RetrieveInfo(peer)
			}
		}
		time.Sleep(time.Second * 5)
	}
}
