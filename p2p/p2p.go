package p2p

import (
	"fmt"
	"goP2PNetwork/config"
	"goP2PNetwork/poisson"
	"log"
	"math/rand"
	"time"
)


func (n *Node) EventGenerator(){
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))

	lambda := 2.0 // rate of 2 updates per minute
	poissonProcess := poisson.PoissonProcess{Lambda: lambda, Rng: rng}
	totalRequests := 0
	t := 1 // current minute

	for {	
		currentTime := 0.0
		previousTime := 0.0
		
		nRequests := poissonProcess.PoissonRandom()
		for i := 1; i <= nRequests; i++ {
			totalRequests ++

			log.Printf(config.Green+"Minute:%d Nrequests:%d"+config.Reset, t, nRequests)


			// get the time for the next request to be executed
			interArrivalTime := poissonProcess.ExponentialRandom()
			previousTime = currentTime
			currentTime = (interArrivalTime * 60) + currentTime


			log.Printf(config.Green+"Request %d at %f seconds\n"+config.Reset, i, currentTime)
			log.Printf(config.Green+"Sleep %.5f seconds...\n"+config.Reset, float64(currentTime-previousTime))
			
			delta := time.Duration(currentTime-previousTime) * time.Second
			time.Sleep(delta)
			n.PushUpdates()
			if i == nRequests && currentTime < 60 {
				log.Printf(config.Green+"Requests for the minute %d endend before finishing the 60s.\nWaiting %f seconds to complete the cycle of 60s....\n"+config.Reset, t, float64(60-currentTime))
				time.Sleep((time.Duration(60-currentTime) * time.Second))
			}


		}
		fmt.Println()
		log.Printf(config.Green+"Statistics: Total requests: %d Minutes spent: %d rate:%f\n"+config.Reset, totalRequests, t, float64(totalRequests)/float64(t))
		t++
		time.Sleep(time.Second * 5)
	}
}

func (n *Node) PushUpdates(){
	for i := range n.Neighbours {
		if n.Neighbours[i] != n.Data.Address {
			LocalNetworkMap.NodeNeighbour(n.Neighbours[i])
		}	
	}
}
