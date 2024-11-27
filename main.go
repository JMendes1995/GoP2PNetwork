package main

import (
	"fmt"
	"goP2PNetwork/p2p"
	"log"
	"os"
)

func main(){
	peer := p2p.Peer{
		UID: "m1",
		Address: "127.0.0.1:50003",
		NeighbourList: make([]string, 0),
	}
	peer.NeighbourList = append(peer.NeighbourList, "10.10.1.1")
	fmt.Print(peer)
	fmt.Println()


	if os.Args[1] == "server" {
		go p2p.Server()
	} else if os.Args[1] == "peer" {
		go p2p.Client(peer)

	} else {
		log.Fatalf("node type requires to be peer or server")
	}
	select {}
}
