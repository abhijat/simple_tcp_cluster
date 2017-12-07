package simple_tcp_cluster

import (
	"net"
	"os"
	"bufio"
	"strings"
	"encoding/json"
	"time"
)

type ClusterState struct {
	ServerName string
	Peers      map[string]time.Time
}

func (c ClusterState) String() string {
	data, _ := json.MarshalIndent(c, " ", " ")
	return string(data)
}

var (
	clusterState = ClusterState{}
)

func ServerLoop(address string) {

	clusterState.ServerName = address
	clusterState.Peers = make(map[string]time.Time)

	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		logger.Println(err)
		os.Exit(1)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logger.Println(err)
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Println(err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	if err != nil {
		logger.Println(err)
		return
	}

	message := string(msg)
	message = strings.TrimSpace(message)

	if message == "request" {
		json.NewEncoder(conn).Encode(clusterState)
	} else {
		conn.Write([]byte("error"))
	}
}
