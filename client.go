package simple_tcp_cluster

import (
	"net"
	"io/ioutil"
	"log"
	"os"
	"encoding/json"
	"time"
)

var (
	logger *log.Logger
)

func init() {
	logger = log.New(os.Stdout, "cluster: ", log.LstdFlags)
}

func RetrieveInfo(address string) {

	conn, err := net.Dial("tcp", address)
	if err != nil {
		logger.Println(err)
		return
	}

	_, err = conn.Write([]byte("request\n"))
	if err != nil {
		logger.Println(err)
		return
	}

	result, err := ioutil.ReadAll(conn)
	if err != nil {
		logger.Println(err)
		return
	}

	v := ClusterState{}
	err = json.Unmarshal(result, &v)
	if err != nil {
		logger.Println(err)
		return
	}

	clusterState.Peers[v.ServerName] = time.Now()

	// add peers which others have, but we do not
	for peer, timeStamp := range v.Peers {
		if _, ok := clusterState.Peers[peer]; !ok && peer != clusterState.ServerName {
			log.Println("adding entry from peer list", peer)
			clusterState.Peers[peer] = timeStamp
		}
	}

	logger.Println(clusterState)
}
