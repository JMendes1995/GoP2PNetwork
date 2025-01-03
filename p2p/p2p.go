package p2p

import (
	"context"
	"encoding/json"
	"fmt"
	"goP2PNetwork/config"
	ppn "goP2PNetwork/p2p/proto"
	"log"
	"sync"
	"time"
)

type Node struct {
	Data       *NodeData
	Neighbours []string
	Mutex      sync.Mutex
}
type NodeData struct {
	Address        string
	UID            string `json:"UID"`
	ValidTimestamp int64  `json:"ValidTimestampMap"`
}

type NetworkMapData struct {
	UID     string           `json:"UID"`
	NodeMap map[string]int64 `json:"NodeMap"`
}
type NetworkMap struct {
	Data  *NetworkMapData
	Mutex sync.Mutex
}

const (
	Port = ":50003"
)

var (
	LocalNode       Node
	LocalNetworkMap NetworkMap
	PeerAddr        string
	LocalAddr       string
	Uid             string
)

type P2PNetworkServer struct {
	ppn.UnimplementedP2PNetworkServer
}

// reciebve
func (s *P2PNetworkServer) NeighbourList(ctx context.Context, in *ppn.NeighbourListRequest) (*ppn.NeighbourListResponse, error) {
	// sending result to client
	nbMap := NetworkMapDataJsonUnmarshal(in.GetPeer())

	fmt.Printf("Recived Update %s\nOld timestamp: %d\nNew Timestamp: %d\n", nbMap.UID, LocalNetworkMap.Data.NodeMap[nbMap.UID], nbMap.NodeMap[nbMap.UID])

	for node, timeStamp := range nbMap.NodeMap {
		LocalNetworkMap.UpdateNetworkMap(node, timeStamp)
	}
	return &ppn.NeighbourListResponse{Result: fmt.Sprintf("Neighbourlist %s updated", LocalAddr)}, nil
}

// Recive a gRPC with the peer node data (ID, ADDRESS and timestamp=0) 
// After adding the new node to the list Sende as response the node info
func (s *P2PNetworkServer) PeerRegist(ctx context.Context, in *ppn.PeerRegistRequest) (*ppn.PeerRegistResponse, error) {
	peerNodeData := NodeDataJsonUnmarshal(in.GetPeer())
	LocalNode.Mutex.Lock()
	LocalNode.AddPeer(peerNodeData)

	entryMaxTime := IncrementTimestamp()
	LocalNode.Data.ValidTimestamp = entryMaxTime
	jsonNodeData, _ := json.Marshal(LocalNode.Data)
	LocalNode.Mutex.Unlock()

	log.Printf(config.Cyan+"Regist new peer %s to local map\n"+config.Reset, peerNodeData.UID)
	return &ppn.PeerRegistResponse{Result: string(jsonNodeData)}, nil
}



func (m *NetworkMap) UpdateNetworkMap(node string, timeStamp int64) {
	m.Mutex.Lock()
	if m.Data.NodeMap[node] != timeStamp && m.Data.NodeMap[node] < timeStamp {
		m.Data.NodeMap[node] = timeStamp
	}
	m.Mutex.Unlock()
}


func (m *NetworkMap) CheckMapEntryExpiration() {
	for {
		if len(m.Data.NodeMap) > 0 {
			for uid, maxTime := range m.Data.NodeMap {
				if maxTime <= time.Now().Unix() && uid != LocalNode.Data.UID && maxTime != 0 {
					log.Printf(config.Yellow+"TTL reached for Address: %s with Maxtime: %s\n"+config.Reset, uid, time.Unix(maxTime, 0))

					m.Mutex.Lock()
					delete(m.Data.NodeMap, uid)
					m.Mutex.Unlock()

				}
			}
		}
		time.Sleep(time.Second * 5)
	}
}

func (m *NetworkMap) CountNodes() {
	for {
		var nodesSlice []string

		jsonLocalMap, _ := json.MarshalIndent(m.Data.NodeMap, "", "    ")
		log.Printf("Local map: %s\n", config.Cyan+string(jsonLocalMap)+config.Reset)
		for key := range m.Data.NodeMap {
			nodesSlice = append(nodesSlice, key)
		}

		log.Printf(config.Blue+"Number nodes in the network: %d -> %s\n"+config.Reset, len(nodesSlice), nodesSlice)
		time.Sleep(time.Second * 30)
	}
}

