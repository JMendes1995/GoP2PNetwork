package config

import (
	"context"
	"net"

	"google.golang.org/grpc/peer"
)

var(
	Reset = "\033[0m"
    bold = "\033[1m"
    underline = "\033[4m"
    strike = "\033[9m"
    italic = "\033[3m"

    Red = "\033[31m"
    Green = "\033[32m"
    Yellow = "\033[33m"
    Blue = "\033[34m"
    cPurple = "\033[35m"
    Cyan = "\033[36m"
    cWhite = "\033[37m"
)
func GetPeerMetadata(ctx context.Context) (string, string, error) {
	p, _ := peer.FromContext(ctx)

	addr := p.Addr
	ip, port, err := net.SplitHostPort(addr.String())
	if err != nil {
		return "", "", err
	}
	return ip, port, nil

}
