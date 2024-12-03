package p2p

import (
	"context"
	"encoding/json"
	"fmt"
	"goP2PNetwork/config"
	ppn "goP2PNetwork/p2p/proto"
	"log"
	"os"
	"slices"
	"sync"
	"time"
)
type Node struct {
  Data *NodeData 
  Neighbours []string
  Mutex sync.Mutex 
}
type NodeData struct {
  Address string 
  UID string `json:"UID"`
  ValidTimestamp int64 `json:"ValidTimestampMap"`
}

type NetworkMapData struct {
  UID string `json:"UID"`
  NodeMap map[string]int64 `json:"NodeMap"`
}
type NetworkMap struct {
  Data *NetworkMapData
  Mutex sync.Mutex 
}

var (
	LocalNode Node
    LocalNetworkMap NetworkMap
	PeerAddr = os.Getenv("NGH_ADDR")+":50003"
	LocalAddr = os.Getenv("LOCAL_NODE")+":50003" 
	Uid = os.Getenv("UID")
)


type P2PNetworkServer struct {
	ppn.UnimplementedP2PNetworkServer
}

func (s *P2PNetworkServer) NeighbourList(ctx context.Context, in *ppn.NeighbourListRequest) (*ppn.NeighbourListResponse, error) {
	// sending result to client
	nbMap := NetworkMapDataJsonUnmarshal(in.GetPeer())


    fmt.Printf("Recived Update %s\n Old timestamp: %d\nNew Timestamp: %d\n", nbMap.UID, LocalNetworkMap.Data.NodeMap[nbMap.UID], nbMap.NodeMap[nbMap.UID])

    for node, timeStamp := range nbMap.NodeMap {
      LocalNetworkMap.UpdateNetworkMap(node, timeStamp)
    }
	return &ppn.NeighbourListResponse{Result: fmt.Sprintf("Neighbourlist %s updated", LocalAddr) }, nil
}

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

func (n *Node) AddPeer(peerData NodeData){
  if !slices.Contains(n.Neighbours, peerData.Address){
    n.Neighbours = append(n.Neighbours, peerData.Address)

    LocalNetworkMap.Mutex.Lock()
    LocalNetworkMap.Data.NodeMap[peerData.UID] = peerData.ValidTimestamp
    LocalNetworkMap.Mutex.Unlock()
  }
}

func (m *NetworkMap) UpdateNetworkMap(node string, timeStamp int64){
  m.Mutex.Lock()
  if m.Data.NodeMap[node] != timeStamp && m.Data.NodeMap[node] < timeStamp {
    m.Data.NodeMap[node] = timeStamp
  } 
  m.Mutex.Unlock()
}

func IncrementTimestamp() int64 {
  entryMaxTime := time.Now().Add(5 * time.Minute).Unix()
  return entryMaxTime
}

func (m *NetworkMap) CheckMapEntryExpiration(){
  for {
    if len(m.Data.NodeMap) > 0 {
      for uid, maxTime:= range m.Data.NodeMap {
        if maxTime <= time.Now().Unix() && uid != LocalNode.Data.UID && maxTime != 0{
            log.Printf(config.Yellow+"TTL reached for Address: %s with Maxtime: %s\n"+config.Reset, uid, time.Unix(maxTime, 0))

            m.Mutex.Lock()
            delete(m.Data.NodeMap, uid)
            m.Mutex.Unlock()

        }
      }
    }
    time.Sleep(time.Second*5)
  }
}

func (m *NetworkMap) CountNodes(){
  for {
    var nodesSlice []string

    jsonLocalMap, _ := json.MarshalIndent(m.Data.NodeMap,"","    ")
    log.Printf("Local map: %s\n",config.Cyan+string(jsonLocalMap)+config.Reset)
    for key := range m.Data.NodeMap {
      nodesSlice = append(nodesSlice, key)
    }

    log.Printf(config.Blue+"Number nodes in the network: %d -> %s\n"+config.Reset,len(nodesSlice), nodesSlice)
    time.Sleep(time.Second*30)
  }
}

func (n *Node) RemovePeer(addr string){
  n.Mutex.Lock()
  // Ensures that only attempts to remove is the map has somehting
  if len(n.Neighbours) == 0 {
    log.Print("No Neighbours in the Network")
  }else {
    if slices.Contains(n.Neighbours, addr){
      peerIndex := slices.Index(n.Neighbours, addr)
      n.Neighbours = append(n.Neighbours[:peerIndex], n.Neighbours[peerIndex+1:]...)
    }
  }
  n.Mutex.Unlock()
}
