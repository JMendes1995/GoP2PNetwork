package p2p

import (
	"encoding/json"
	"flag"
	ppn "goP2PNetwork/p2p/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

//func NodeJsonUnmarshal(data string) Node {
//	var node Node
//	err := json.Unmarshal([]byte(data), &node)
//	if err != nil {
//		log.Fatalf("Error in Unmarshalling Node: %s", err)
//	}
//	return node
//}


func NeighboursJsonUnmarshal(data string) NeigbboursData {
	var nbData NeigbboursData
	err := json.Unmarshal([]byte(data), &nbData)
	if err != nil {
		log.Fatalf("Error in Unmarshalling Node: %s", err)
	}
	return nbData
}

func NeighboursMapJsonUnmarshal(data string) NeighboursMapData {
	var nbMap NeighboursMapData
	err := json.Unmarshal([]byte(data), &nbMap)
	if err != nil {
		log.Fatalf("Error in Unmarshalling NeighboursMap: %s", err)
	}
	return nbMap
}

func Server() {
	flag.Parse()
	lis, err := net.Listen("tcp", LocalAddr)
	if err != nil {
		log.Fatalf("failed to listen from calc server: %v", err)
	}
	//start grpc server
	server := grpc.NewServer()
	//register server
	ppn.RegisterP2PNetworkServer(server, &P2PNetworkServer{})
	//logger.Info("P2P server listening at %v", lis.Addr())
	log.Printf("P2P server listening at %v", lis.Addr())

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
