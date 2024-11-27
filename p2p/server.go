package p2p

import (
	"context"
	"flag"
	ppn "goP2PNetwork/p2p/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)


type P2PNetworkServer struct {
	ppn.UnimplementedP2PNetworkServer
}

func (s *P2PNetworkServer) NeighbourList(ctx context.Context, in *ppn.NeighbourListRequest) (*ppn.NeighbourListResponse, error) {

	// sending result to client
	return &ppn.NeighbourListResponse{Result: ""}, nil
}

func (s *P2PNetworkServer) PeerRegist(ctx context.Context, in *ppn.PeerRegistRequest) (*ppn.PeerRegistResponse, error) {

	// sending result to client
	return &ppn.PeerRegistResponse{Result: "Peer registred"}, nil
}

func Server() {
	flag.Parse()
	lis, err := net.Listen("tcp", "127.0.0.1:50003")
	if err != nil {
		log.Fatalf("failed to listen from calc server: %v", err)
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
