package p2p

import (
	"encoding/json"
	"flag"
	ppn "goP2PNetwork/p2p/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

func NodeDataJsonUnmarshal(data string) NodeData {
	var node NodeData
	err := json.Unmarshal([]byte(data), &node)
	if err != nil {
		log.Fatalf("Error in Unmarshalling Node: %s", err)
	}
	return node
}
func NetworkMapDataJsonUnmarshal(data string) NetworkMapData {
	var nbData NetworkMapData
	err := json.Unmarshal([]byte(data), &nbData)
	if err != nil {
		log.Fatalf("Error in Unmarshalling Node: %s", err)
	}
	return nbData
}

func Server() {
	flag.Parse()
	lis, err := net.Listen("tcp", LocalAddr+Port)
	if err != nil {
		log.Fatalf("failed to listen from server: %v", err)
	}
	//start grpc server
	server := grpc.NewServer()
	//register server
	ppn.RegisterP2PNetworkServer(server, &P2PNetworkServer{})
	log.Printf("P2P server listening at %v", lis.Addr())

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
