package p2p

import (
	"context"
	"encoding/json"
	ppn "goP2PNetwork/p2p/proto"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Client(peer Peer) {
	conn, err := grpc.NewClient("127.0.0.1:50003", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error connecting to server %s", err)
	}

	c := ppn.NewP2PNetworkClient(conn)
	defer conn.Close()
	// setupt connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) 
	defer cancel()

	peerJson, _ := json.Marshal(peer)
	response, err := c.PeerRegist(ctx, &ppn.PeerRegistRequest{Peer: string(peerJson)})
	if err != nil {
		log.Fatalf("Error requesting message from server %v", err)
	}
	log.Printf("Response: %s\n", response.Result)
	
}
