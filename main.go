package main

import (
	"flag"
	"fmt"
	"goP2PNetwork/p2p"
	"os"
	"strings"
)

func cli() (string, string) {
	// Define flags
	id := flag.String("id", "", "Specify the Local Node ID")
	neighboursAddress := flag.String("neighboursAddress", "", "Specify the list of NeighbourS addresses separated by comma")

	help := flag.Bool("help", false, "help")

	// Parse flags
	flag.Parse()

	if *help {
		fmt.Println("Usage: cli [options]")
		flag.PrintDefaults()
		os.Exit(0)
	}

	return *id, *neighboursAddress
}

func main() {
	id, neighboursAddress := cli()

	localAddress, _ := p2p.GetLocalAddr()
	p2p.LocalAddr = localAddress
	p2p.ListPeerAddr = strings.Split(neighboursAddress, ",")
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
	for i := range p2p.ListPeerAddr {
		go p2p.LocalNode.NodeRegisterClient(p2p.ListPeerAddr[i])
	}

	go p2p.LocalNetworkMap.CheckMapEntryExpiration()
	go p2p.LocalNode.EventGenerator()
	go p2p.LocalNetworkMap.CountNodes()
	//
	select {}
}
