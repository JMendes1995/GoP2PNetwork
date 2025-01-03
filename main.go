package main

import (
	"flag"
	"fmt"
	"goP2PNetwork/p2p"
	"os"
)

func cli() (string, string) {
	// Define flags
	id := flag.String("id", "", "Specify the Local Node ID")
	neighbourAddress := flag.String("neighbourAddress", "", "Specify the Neighbour address")

	help := flag.Bool("help", false, "help")

	// Parse flags
	flag.Parse()

	if *help {
		fmt.Println("Usage: cli [options]")
		flag.PrintDefaults()
		os.Exit(0)
	}

	return *id, *neighbourAddress
}

func main() {
	id, neighbourAddress := cli()

	localAddress, _ := p2p.GetLocalAddr()
	p2p.LocalAddr = localAddress
	p2p.PeerAddr = neighbourAddress
	p2p.Uid = id
	p2p.LocalNode = p2p.Node{
		Data: &p2p.NodeData{
			UID:     id,
			Address: localAddress,
		},
		Neighbours: make([]string, 0),
	}
	p2p.LocalNetworkMap = p2p.NetworkMap{
		Data: &p2p.NetworkMapData{
			UID:     p2p.Uid,
			NodeMap: make(map[string]int64),
		},
	}

	go p2p.Server()
	go p2p.LocalNode.NodeRegisterClient()
	go p2p.LocalNetworkMap.CheckMapEntryExpiration()
	go p2p.LocalNode.EventGenerator()
	go p2p.LocalNetworkMap.CountNodes()
	//
	select {}
}
