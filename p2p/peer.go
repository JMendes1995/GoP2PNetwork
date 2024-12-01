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
  UID string 
  Address string 
  Neighbours *NeigbboursData 
  Mutex sync.Mutex 
}
type NeigbboursData struct {
  Data []string `json:"Neighbours"`
  Version int `json:"Version"`
}
// The creation of NeighboursMapData serves the purpose of not sending to the client the mutex field.
type NeighboursMapData struct {
	Address string `json:"Address"`
	MapPeer map[string]NeigbboursData `json:"Map"`
}

// Creation of substruct to lock only the ValidTimestampMap withlocking the entire stuckt
// This is ment to increase scalability and efficiency.
type NeighboursMap struct {
  Data *NeighboursMapData
  Mutex sync.Mutex 
  ValidTimestampMap  struct{
      Data map[string]int64
      Mutex sync.Mutex 
  }
}

var (
	LocalNode Node
	LocalNeighboursMap NeighboursMap
	PeerAddr = os.Getenv("NGH_ADDR")+":50003"
	LocalAddr = os.Getenv("LOCAL_NODE")+":50003" 
	Uid = os.Getenv("UID")
)


type P2PNetworkServer struct {
	ppn.UnimplementedP2PNetworkServer
}

func (s *P2PNetworkServer) NeighbourList(ctx context.Context, in *ppn.NeighbourListRequest) (*ppn.NeighbourListResponse, error) {
	// sending result to client
	nbMap := NeighboursMapJsonUnmarshal(in.GetPeer())
    //jsonNbMap, _ := json.Marshal(nbMap.MapPeer)
    //log.Printf("\nNew update from peer %s: %s\n", nbMap.Address, jsonNbMap)
    
    LocalNeighboursMap.UpdateNeighboursMap(nbMap.MapPeer, nbMap.Address)
	return &ppn.NeighbourListResponse{Result: fmt.Sprintf("Neighbourlist %s updated", LocalAddr) }, nil
}

func (s *P2PNetworkServer) PeerRegist(ctx context.Context, in *ppn.PeerRegistRequest) (*ppn.PeerRegistResponse, error) {
    
  //peerInfo := NodeJsonUnmarshal(in.GetPeer())
  LocalNode.Mutex.Lock()
  if !slices.Contains(LocalNode.Neighbours.Data, in.GetPeer()){
      LocalNode.Neighbours.Data = append(LocalNode.Neighbours.Data, in.GetPeer())
      LocalNode.Neighbours.Version++

      LocalNeighboursMap.Mutex.Lock()
      LocalNeighboursMap.Data.MapPeer[LocalNode.Address] = *LocalNode.Neighbours
      LocalNeighboursMap.Mutex.Unlock()
  }
  LocalNode.Mutex.Unlock()
  jsonLocalMap, _ := json.Marshal(*LocalNeighboursMap.Data)

  //log.Printf("\nRegist new peer %s to local map %s\n", LocalNeighboursMap.Data.Address, jsonLocalMap )    
  log.Printf(config.Cyan+"Regist new peer %s to local map\n"+config.Reset, in.GetPeer())    
  //fmt.Printf("\njsonmap: %s\n", localmap)    
  return &ppn.PeerRegistResponse{Result: string(jsonLocalMap)}, nil
}


func (m *NeighboursMap) AddValidationTime(addr string) {
  entryMaxTime := time.Now().Add(time.Second * 100).Unix()
  
  m.ValidTimestampMap.Mutex.Lock()
  m.ValidTimestampMap.Data[addr] = entryMaxTime
  m.ValidTimestampMap.Mutex.Unlock()

  log.Printf(config.Yellow+"Added validation to entry from Peer %s with timestamp %d\n"+config.Reset, addr, m.ValidTimestampMap.Data[addr])  
}

func (m *NeighboursMap) CheckMapEntryExpiration(){
  for {
    if len(m.ValidTimestampMap.Data) > 0 {
      for addr, maxTime:= range m.ValidTimestampMap.Data {
        if maxTime <= time.Now().Unix(){
            log.Printf(config.Yellow+"TTL reached for Address: %s with Maxtime: %s\n"+config.Reset, addr, time.Unix(maxTime, 0))
            //remove peer.
            m.RemovePeer(addr)
            m.ValidTimestampMap.Mutex.Lock()
            delete(m.ValidTimestampMap.Data, addr)
            m.ValidTimestampMap.Mutex.Unlock()  
        }
      }
    }
    jsonLocalMap, _ := json.MarshalIndent(m.Data.MapPeer,"","    ")
    log.Printf("Local map: %s\n",config.Cyan+string(jsonLocalMap)+config.Reset)
    m.CountNodes()
    time.Sleep(time.Second*5)
  }
}


// Function that updates Neighbours map
func (m *NeighboursMap) UpdateNeighboursMap(peerMap map[string]NeigbboursData, peerAddr string){
  m.Mutex.Lock()
  for key, value := range peerMap {
    // if the update recived is not the local node and if the version of the map recived is higher than the one existing in the map 
    if key != LocalNode.Address && m.Data.MapPeer[key].Version <= peerMap[key].Version{
      m.Data.MapPeer[key] = value
      m.AddValidationTime(key)
    }
  }
  //jsonLocalMap, _ := json.Marshal(m.Data.MapPeer)
  //log.Printf("\nUpdated Local map: %s\n",jsonLocalMap)
  m.Mutex.Unlock()
}


func (m *NeighboursMap) RemovePeer(addr string){
  LocalNode.Mutex.Lock()
  // Ensures that only attempts to remove is the map has somehting
  if len(LocalNode.Neighbours.Data) == 0 {
    log.Print("No Neighbours in the Network")
  }else {
    if slices.Contains(LocalNode.Neighbours.Data, addr){
      peerIndex := slices.Index(LocalNode.Neighbours.Data, addr)
      LocalNode.Neighbours.Data = append(LocalNode.Neighbours.Data[:peerIndex], LocalNode.Neighbours.Data[peerIndex+1:]...)
      LocalNode.Neighbours.Version++
    }
    LocalNode.Mutex.Unlock()
  }
  m.Mutex.Lock()
  //remove peer from map
  delete(m.Data.MapPeer, addr)
  // update the local map with the new neighbours list.
  m.Data.MapPeer[LocalAddr] = *LocalNode.Neighbours
  m.Mutex.Unlock()
}


func (m *NeighboursMap) CountNodes(){
  var nodesSlice []string
  for key := range m.Data.MapPeer {
    nodesSlice = append(nodesSlice, key)
  }

  log.Printf(config.Blue+"Number nodes in the network: %d -> %s\n"+config.Reset,len(nodesSlice), nodesSlice)
}
