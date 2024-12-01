package p2p

import (
	"context"
	"encoding/json"
	"fmt"
	"goP2PNetwork/config"
	ppn "goP2PNetwork/p2p/proto"
	"log"
	"slices"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (n *Node) NodeRegisterClient()  {
  conn, err := grpc.NewClient(PeerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
  if err != nil {
      log.Fatalf("Error connecting to server %s", err)
  }

  c := ppn.NewP2PNetworkClient(conn)
  defer conn.Close()

  // Prevent exiting due to error of client not reaching the peer
  done := false
  for !done{
      // setup connection with timeout
      ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) 
      defer cancel()

      peerPullMapJson, nodeRegisterError := c.PeerRegist(ctx, &ppn.PeerRegistRequest{Peer: string(LocalAddr)})

      if nodeRegisterError == nil {
          //log.Printf("\nRecived Map from Node %s: %s\n", PeerAddr, peerPullMapJson )  
          peerMapData := NeighboursMapJsonUnmarshal(peerPullMapJson.Result)
          fmt.Print(peerMapData.Address)
          if !slices.Contains(n.Neighbours.Data, peerMapData.Address){
            n.Mutex.Lock()
            n.Neighbours.Data = append(n.Neighbours.Data, peerMapData.Address)

            LocalNeighboursMap.Mutex.Lock()
            LocalNeighboursMap.Data.MapPeer[n.Address] = *n.Neighbours
            LocalNeighboursMap.Mutex.Unlock()
            n.Mutex.Unlock()
          }
          LocalNeighboursMap.UpdateNeighboursMap(peerMapData.MapPeer, peerMapData.Address)

          done = true
      }else{
          log.Printf(config.Yellow+"Peer not ready %s"+config.Reset, PeerAddr)	
      }
      time.Sleep(time.Second * 10)
  }
  n.EventGenerator()
}

func (m *NeighboursMap) NodeNeighbour(peerAddr string) {

	conn, err := grpc.NewClient(peerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error connecting to server %s", err)
	}

	c := ppn.NewP2PNetworkClient(conn)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) 
	defer cancel()

	filteredMap := m.FilterMap()
	nbmapJson, _ := json.Marshal(filteredMap)
	//log.Printf("\nSending Map update to node %s: %s\n",PeerAddr, nbmapJson)
	response, err := c.NeighbourList(ctx, &ppn.NeighbourListRequest{Peer: string(nbmapJson)})
	if err == nil {
		log.Printf(config.Cyan+"%s\n"+config.Reset, response.Result)
	}	
}

//Filter The map to be sent to neighbours.
//This filter consists in sending only maps by direct link nodes.
func (m *NeighboursMap) FilterMap() NeighboursMapData {
  // create new temporary NeighboursMapData struct to store filtered maps
  filteredMap := NeighboursMapData{
    Address: LocalNode.Address,
    MapPeer: make(map[string]NeigbboursData),
  }
  m.Mutex.Lock()
  filteredMap.MapPeer[LocalAddr] = m.Data.MapPeer[LocalAddr]
  for i := range LocalNode.Neighbours.Data {
    filteredMap.MapPeer[LocalNode.Neighbours.Data[i]] = m.Data.MapPeer[LocalNode.Neighbours.Data[i]]
  }
  m.Mutex.Unlock()
  return filteredMap
}
