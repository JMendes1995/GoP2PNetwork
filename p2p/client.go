package p2p

import (
	"context"
	"encoding/json"
	"fmt"
	"goP2PNetwork/config"
	ppn "goP2PNetwork/p2p/proto"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (n *Node) NodeRegisterClient(peerAddr string) {
	conn, err := grpc.NewClient(peerAddr+Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error connecting to server %s", err)
	}

	c := ppn.NewP2PNetworkClient(conn)
	defer conn.Close()

	// Prevent exiting due to error of client not reaching the peer
	done := false
	for !done {
		// setup connection with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		jsonNodeData, _ := json.Marshal(n.Data)
		peerNodeDataJson, nodeRegisterError := c.PeerRegist(ctx, &ppn.PeerRegistRequest{Peer: string(jsonNodeData)})

		// if the connection is successfull
		// save the node info recived in the response message into the Neighbours slice
		if nodeRegisterError == nil {
			peerNodeData := NodeDataJsonUnmarshal(peerNodeDataJson.Result)
			fmt.Print(peerNodeData.UID)
			n.AddPeer(peerNodeData)

			done = true
		} else {
			log.Printf(config.Yellow+"Peer not ready %s"+config.Reset, peerAddr+Port)
		}
		time.Sleep(time.Second * 5)
	}
}

func (m *NetworkMap) NodeNeighbour(peerAddr string) {

	conn, err := grpc.NewClient(peerAddr+Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Error connecting to server %s", err)
		// try to remove from linked nodes
		LocalNode.RemovePeer(peerAddr + Port)
	}

	c := ppn.NewP2PNetworkClient(conn)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	m.Mutex.Lock()
	entryMaxTime := IncrementTimestamp()
	m.Data.NodeMap[LocalNode.Data.UID] = entryMaxTime
	nbmapJson, _ := json.Marshal(m.Data)
	m.Mutex.Unlock()

	c.NeighbourList(ctx, &ppn.NeighbourListRequest{Peer: string(nbmapJson)})

}
