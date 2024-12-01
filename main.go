package main

import (
	"goP2PNetwork/p2p"
	"sync"
)

func main(){
    p2p.LocalNode = p2p.Node{
		UID: p2p.Uid,
		Address: p2p.LocalAddr,
		Neighbours: &p2p.NeigbboursData{
			Data: make([]string, 0),
			Version: 1,
		},
	}
    p2p.LocalNeighboursMap = p2p.NeighboursMap{
		Data: &p2p.NeighboursMapData{
			Address: p2p.LocalAddr,
			MapPeer: make(map[string]p2p.NeigbboursData),
		},
		ValidTimestampMap: struct{Data map[string]int64; Mutex sync.Mutex} {
			Data: make(map[string]int64),
		},
    }

	//peer.Neighbours = append(peer.Neighbours, "10.10.1.1")
    //mapnode.MapPeer[peer.UID] = peer.Neighbours

	go p2p.Server()
	go p2p.LocalNode.NodeRegisterClient()
	go p2p.LocalNeighboursMap.CheckMapEntryExpiration()
    //go p2p.LocalNeighboursMap.EventGenerator()

	select {}
}
