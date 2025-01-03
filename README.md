# GoP2PNetwork

This Project counts the Number of Nodes in a P2P Network

# Table Of Contents

- [GoP2PNetwork](#gop2pnetwork)
  - [Implementation Requirements](#implementation-requirements)
  - [How to build the project](#how-to-build-the-project)
  - [How to execute the project](#how-to-execute-the-project)
    - [Execution in physical nodes](#execution-in-physical-nodes)
    - [Docker](#docker)
      - [Clean up docker environment](#clean-up-docker-environment)
    - [Results intrepertation](#results-intrepertation)
      - [Server startup and Peer Registration](#server-startup-and-peer-registration)
      - [Event generator](#event-generator)
      - [Recive update message](#recive-update-message)
      - [Peers map with timestamps](#peers-map-with-timestamps)
      - [Node Count](#node-count)
      - [Node removel](#node-removel)


### Implementation Requirements
* Create a network of 6 peers (p1 to p6), running on different machines (m1 to m6), with
the topology shown in the figure above.
* Peers must have a thread waiting for others to connect and register themselves. This is
how peers get to know each other and the network is built.
* Each peer keeps a map data structure (e.g., Map or HashMap in Java) with the
IPs/names of the machines of the peers that have registered themselves with it. For
example, peer p2 (in machine m2) keeps the entries [m1, m3, m4] whereas peer p6 (in
machine m6) keeps only [m4].
* Now implement the Anti-Entropy Algorithm so peers can disseminate their maps. Each
peer updates its map twice per minute following a Poisson distribution, each time print-
ing the number of peers in the updated map.
* Once the dissemination algorithm starts, the number of items in the map at each peer
will approximate the total number of peers in the network.
* What happens to the number of peers if you remove some of them from the above
network, e.g., killing the processes? Devise a scheme that handles this situation (hint:
you can use timestamps for the map entries and remove entries older than a given
threshold).

### How to build the project
In the root directory of the project are present 2 binaries from **Linux (goP2PNetwork-linux)** and **MacOS (goP2PNetwork-macos)** arm architecture. However is possible to build the project manually by executing the following commands.

```bash
$ cd GoP2PNetwork
$ go build
```

### How to execute the project
To execute the project there is 2 possible ways.
1. Execute the code in phsical nodes
2. Build docker environment using docker compose.


#### Execution in physical nodes
Before executing the project, is required to pull the source code into the nodes that will be part of the network.
Furthermore, is also require to execute the command passing 2 input parametes:

* `--id <NodeID>`
* `--neighbourAddress <NeighbourAddress>`

Example:
```bash
$ cd GoP2PNetwork
$ ./goP2PNetwork-linux --id m1 --neighbourAddress 10.10.1.2
```

#### Docker

To execute the project using docker is required to have the docker and docker compose prior installed. Then execute the following commands.
Note: replace the neighbourAddress value to your values

```bash
$ docker compose up -d
```

```bash
$ docker exec -it gop2pnetwork-peer1-1 bash
$ ./goP2PNetwork-linux --id m1 --neighboursAddress 10.10.1.2
```
```bash
$ docker exec -it gop2pnetwork-peer2-1 bash
$ ./goP2PNetwork-linux  --id m2 --neighboursAddress 10.10.1.1,10.10.1.3,10.10.1.4
```
```bash
$ docker exec -it gop2pnetwork-peer3-1 bash
$ ./goP2PNetwork-linux  --id m3 --neighboursAddress 10.10.1.2
```
```bash
$ docker exec -it gop2pnetwork-peer4-1 bash
$ ./goP2PNetwork-linux  --id m4 --neighboursAddress 10.10.1.2,10.10.1.5,10.10.1.6
```
```bash
$ docker exec -it gop2pnetwork-peer5-1 bash
$ ./goP2PNetwork-linux --id m5 --neighboursAddress 10.10.1.4
```
```bash
$ docker exec -it gop2pnetwork-peer6-1 bash
$ ./goP2PNetwork-linux --id m6 --neighboursAddress 10.10.1.4
```

##### Clean up docker environment
```bash
$docker compose down && docker compose rm && docker rmi -f $(docker images -aq)
```


### Results intrepertation

##### Server startup and Peer Registration
Peers share its own iformation (id and address) with other Peers to form a list of neighbours.

```bash
2025/01/03 00:01:07 P2P server listening at 10.10.1.1:50003
2025/01/03 00:01:07 Regist new peer m2 to local map
```
##### Event generator
The events generated in this project aim to keep all the maps sync and updated. Therefore, the event generator creates 2 events per minute following a Poisson distribution,
```bash
2025/01/03 00:01:07 Minute:1 Nrequests:1
2025/01/03 00:01:07 Request 1 at 2.290438 seconds
2025/01/03 00:01:07 Sleep 2.29044 seconds...
2025/01/03 00:01:09 Requests for the minute 1 endend before finishing the 60s.
Waiting 57.709562 seconds to complete the cycle of 60s....
```
##### Recive update message
Update messages aims to update the timestamp of the neighbours
```bash
Recived Update m2
Old timestamp: 1735862667
New Timestamp: 1735862731
```

##### Peers map with timestamps
Map which the key is the node ID and the value is the timestamp that the nodes are valid
```bash
2025/01/03 00:02:39 Local map: {
    "m1": 1735862717,
    "m2": 1735862731,
    "m3": 1735862731,
    "m4": 1735862674,
    "m5": 1735862730,
    "m6": 1735862668
}
```
##### Node Count
Displayed the current node count with the list of Peer ID's
```bash
2025/01/03 00:02:39 Number nodes in the network: 6 -> [m3 m2 m1 m5 m6 m4]
```

##### Node removel
If a node does not recive any updates from peers. The timestamp will be expired, thus, resulting in node removal.
```bash
2025/01/03 00:06:24 TTL reached for Address: m5 with Maxtime: 2025-01-03 00:06:20 +0000 UTC
```
