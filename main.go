package main

import (
	"goP2PNetwork/p2p"
)

func main(){
    p2p.LocalNode = p2p.Node{
		Data: &p2p.NodeData{
			UID: p2p.LocalAddr,
			Address: p2p.LocalAddr,
		},
		Neighbours: make([]string, 0),
	}
    p2p.LocalNetworkMap = p2p.NetworkMap{
		Data: &p2p.NetworkMapData{
			UID: p2p.LocalAddr,
			NodeMap:  make(map[string]int64),
		},
    }

	go p2p.Server()
	go p2p.LocalNode.NodeRegisterClient()
	go p2p.LocalNetworkMap.CheckMapEntryExpiration()
	go p2p.LocalNode.EventGenerator()
	go p2p.LocalNetworkMap.CountNodes()
	
	select {}
}
