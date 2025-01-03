package p2p

import (
	"log"
	"net"
	"slices"
	"time"
)

func GetLocalAddr() (string, string) {
	// Get a list of all network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	// Loop through the interfaces and get their addresses
	for _, iface := range interfaces {
		// Skip loopback interfaces (127.0.0.1)
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		// Get the addresses for the interface
		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		// Loop through addresses and print the first non-loopback IP address
		for _, addr := range addrs {
			// We check if the address is an IP address (not a MAC address)
			ipnet, ok := addr.(*net.IPNet)
			if ok && !ipnet.IP.IsLoopback() {
				// Print the first non-loopback IP address
				return ipnet.IP.String(), ""
			}
		}
	}
	return "", "No addresses found"
}

// Increment the timestamp for 3 minutes
func IncrementTimestamp() int64 {
	entryMaxTime := time.Now().Add(3 * time.Minute).Unix()
	return entryMaxTime
}

// Add the new node to the Neighbours slice
func (n *Node) AddPeer(peerData NodeData) {
	if !slices.Contains(n.Neighbours, peerData.Address) {
		n.Neighbours = append(n.Neighbours, peerData.Address)

		LocalNetworkMap.Mutex.Lock()
		LocalNetworkMap.Data.NodeMap[peerData.UID] = peerData.ValidTimestamp
		LocalNetworkMap.Mutex.Unlock()
	}
}

// when a node stops responding, remove from the Neighbours slice
func (n *Node) RemovePeer(addr string) {
	n.Mutex.Lock()
	// Ensures that only attempts to remove is the map has somehting
	if len(n.Neighbours) == 0 {
		log.Print("No Neighbours in the Network")
	} else {
		if slices.Contains(n.Neighbours, addr) {
			peerIndex := slices.Index(n.Neighbours, addr)
			n.Neighbours = append(n.Neighbours[:peerIndex], n.Neighbours[peerIndex+1:]...)
		}
	}
	n.Mutex.Unlock()
}
