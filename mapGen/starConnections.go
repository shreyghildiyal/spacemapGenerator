package mapGen

import (
	"fmt"
	"math/rand"

	"github.com/shreyghildiyal/spacemapGenerator/cartesian"
)

func AddDummyNeighbours(stars []Star) {
	addedStars := []int{}     // index in stars of the star
	remainingStars := []int{} // index in stars of the star

	distanceGrid := make([][]float64, len(stars))
	for i := range stars {
		distanceGrid[i] = make([]float64, len(stars))
	}

	for i := 0; i < len(stars); i++ {
		for j := 0; j < len(stars); j++ {
			distanceGrid[i][j] = getDist(stars[i], stars[j]) // distance using the index within stars slice
		}
	}

	addedStars = append(addedStars, 0) // the first star in stars is where we start
	for i := 1; i < len(stars); i++ {
		remainingStars = append(remainingStars, i) // all other indexes are 'remaining'
	}

	for len(remainingStars) > 0 {
		bestRemIndex := 0
		bestStarIndex := remainingStars[bestRemIndex]
		bestDistance, bestFromStarIndex := getMinDist(bestStarIndex, addedStars, distanceGrid)

		for remIndex := 1; remIndex < len(remainingStars); remIndex++ {
			starIndex := remainingStars[remIndex]
			dist, fromStarIndex := getMinDist(starIndex, addedStars, distanceGrid)
			if dist < bestDistance {
				bestDistance = dist
				bestRemIndex = remIndex
				bestStarIndex = starIndex
				bestFromStarIndex = fromStarIndex
			}
		}

		remainingStars = append(remainingStars[:bestRemIndex], remainingStars[bestRemIndex+1:]...)
		addedStars = append(addedStars, bestStarIndex)
		stars[bestFromStarIndex].Neighbours = append(stars[bestFromStarIndex].Neighbours, bestStarIndex)
		stars[bestStarIndex].Neighbours = append(stars[bestStarIndex].Neighbours, bestFromStarIndex)
	}
	fmt.Println("Neighbours have been added. There should now be a path from every star to every other star")
}

type NeighbourConfigs struct {
	MinInClusterConnectionRatio float64 // for each star in a cluster the number of connections to other cluster stars it has divided by the number of such possible connections
	MinClusterConnectionRatio   float64 // for each cluster, number of connections to stars in other clusters/ number of such connections possible
	MaxClusterConnectionRatio   float64 // for each cluster, number of connections to stars in other clusters/ number of such connections possible
	MaxInClusterConnectionRatio float64
}

func AddStarNeighbours(stars []Star, configs NeighbourConfigs) {
	clusterStarsIds := getClusterStarIds(stars)
	for _, clusterStarIds := range clusterStarsIds {
		connectClusterStars(stars, clusterStarIds, configs.MinInClusterConnectionRatio, configs.MaxInClusterConnectionRatio)
	}

	connectClusters(stars, configs.MinClusterConnectionRatio, configs.MaxClusterConnectionRatio)
}

/*
Gives a map of clusterId:list(Id of star in cluster)
*/
func getClusterStarIds(stars []Star) map[int][]int {
	clusterStars := map[int][]int{}
	for _, star := range stars {
		clusterId := star.ClusterId
		if _, ok := clusterStars[clusterId]; !ok {
			clusterStars[clusterId] = []int{}
		}
		clusterStars[clusterId] = append(clusterStars[clusterId], star.Id)
	}
	return clusterStars
}

func connectClusterStars(stars []Star, clusterStarIds []int, minConnectionRatio, maxConnectionRatio float64) {
	// create a minimum spanning tree

	createMinSpanningTree(clusterStarIds, stars)

	// use a probability system to use the remaining borders
	logger.Printf("Created min tree for cluster %d", stars[clusterStarIds[0]].ClusterId)
	for _, starId := range clusterStarIds {
		createNeighbourConnections(starId, stars, minConnectionRatio, maxConnectionRatio)
	}
}

func createNeighbourConnections(starId int, stars []Star, minConnectionRatio, maxConnectionRatio float64) {

	actualIntendedRatio := minConnectionRatio + rand.Float64()*(maxConnectionRatio-minConnectionRatio)

	// neighbourCount := 0

	candidates := []int{}
	candidateValid := map[int]bool{}
	for cId := range stars[starId].Borders {
		if stars[cId].ClusterId == stars[starId].ClusterId {
			candidates = append(candidates, cId)
			candidateValid[cId] = true
		}
	}

	// logger.Println("Created candidates arrays")
	for _, nId := range stars[starId].Neighbours {
		if _, ok := candidateValid[nId]; !ok {
			candidateValid[nId] = false
		}
	}
	// logger.Println("Updated existing neighbours")

	intendedBorderCount := actualIntendedRatio * float64(len(candidates))

	for len(stars[starId].Neighbours) < int(intendedBorderCount) {
		logger.Printf("Neighbour Count: %d, intendedCount: %f, candidatesCount: %d", len(stars[starId].Neighbours), intendedBorderCount, len(candidates))
		selectedIndex := rand.Int() % (len(stars[starId].Borders) - len(stars[starId].Neighbours))
		logger.Printf("Selected Index: %d", selectedIndex)
		foundvalidCount := 0

		for i, cId := range candidates {
			if valid, ok := candidateValid[cId]; ok && valid {

				if foundvalidCount >= selectedIndex {

					stars[starId].Neighbours = append(stars[starId].Neighbours, cId)
					stars[cId].Neighbours = append(stars[cId].Neighbours, starId)
					candidateValid[i] = false

					// do the connection
				} else {
					foundvalidCount++
				}
			}
		}
	}

}

func createMinSpanningTree(clusterStarIds []int, stars []Star) {
	addedStars := map[int]bool{}
	remainingStars := map[int]bool{}

	addedStars[clusterStarIds[0]] = true
	for i := 1; i < len(clusterStarIds); i++ {
		remainingStars[clusterStarIds[i]] = true
	}

	for len(remainingStars) > 0 {
		bestAddedStarId, bestRemainingStarId := getNextConnection(addedStars, remainingStars, stars)
		stars[bestAddedStarId].Neighbours = append(stars[bestAddedStarId].Neighbours, bestRemainingStarId)
		stars[bestRemainingStarId].Neighbours = append(stars[bestRemainingStarId].Neighbours, bestAddedStarId)

		delete(remainingStars, bestRemainingStarId)
		addedStars[bestRemainingStarId] = true
	}
}

func getNextConnection(addedStars map[int]bool, remainingStars map[int]bool, stars []Star) (int, int) {

	bestAddedStarIndex := -1
	bestRemainingStarIndex := -1
	minDist := -1.0

	for addedStarId := range addedStars {
		addedStar := stars[addedStarId]
		for remainingStarId := range addedStar.Borders {
			if _, ok := remainingStars[remainingStarId]; ok {
				dist := getDistSquared(addedStar.Vector2, stars[remainingStarId].Vector2)
				if minDist < 0 {
					bestAddedStarIndex = addedStarId
					bestRemainingStarIndex = remainingStarId
					minDist = dist
				} else {
					if dist < minDist {
						bestAddedStarIndex = addedStarId
						bestRemainingStarIndex = remainingStarId
						minDist = dist
					}
				}
			}
		}
	}

	return bestAddedStarIndex, bestRemainingStarIndex

}

func connectClusters(stars []Star, minConnectionRatio, maxConnectionRatio float64) {
	// cluster potential borders

	potentialConnections := map[int]map[int][]StarConnection{}
	for _, star := range stars {
		if _, ok := potentialConnections[star.ClusterId]; !ok {
			potentialConnections[star.ClusterId] = map[int][]StarConnection{}
		}
		for otherStarId, _ := range star.Borders {
			if stars[otherStarId].ClusterId > star.ClusterId {
				if _, ok := potentialConnections[star.ClusterId][stars[otherStarId].ClusterId]; !ok {
					potentialConnections[star.ClusterId][stars[otherStarId].ClusterId] = []StarConnection{}
				}
				conn := StarConnection{
					FromStar: star.Id,
					ToStar:   otherStarId,
				}
				potentialConnections[star.ClusterId][stars[otherStarId].ClusterId] = append(potentialConnections[star.ClusterId][stars[otherStarId].ClusterId], conn)
			}
		}
	}

	fmt.Println("Generated map for potential connections")

	// now we have all the possible connections that can be made from one cluster to any other

	for clusterId := range potentialConnections {
		fmt.Println("creating inter cluster connections from ", clusterId)
		for toClusterId, potConns := range potentialConnections[clusterId] {
			fmt.Println("generating connections to cluster", toClusterId)
			actualRatio := minConnectionRatio + (maxConnectionRatio-minConnectionRatio)*rand.Float64()
			desiredCount := float64(len(potConns)) * actualRatio
			if desiredCount < 1 {
				desiredCount = 1
			}
			fmt.Println("desired count", desiredCount)
			createdConnections := 0
			validitySlice := make([]bool, len(potConns))
			for i := range validitySlice {
				validitySlice[i] = true
			}

			for createdConnections < int(desiredCount) {
				index := rand.Intn(len(potConns) - createdConnections)
				fmt.Println("selected index", index, len(potConns)-createdConnections)
				conn := getValidItem(potConns, index, validitySlice)
				validitySlice[index] = false
				stars[conn.FromStar].Neighbours = append(stars[conn.FromStar].Neighbours, conn.ToStar)
				stars[conn.ToStar].Neighbours = append(stars[conn.ToStar].Neighbours, conn.FromStar)
				createdConnections++
			}
		}
	}
}

func getValidItem(connections []StarConnection, index int, validitySlice []bool) StarConnection {
	fmt.Println("trying to get desired item")
	foundValid := 0

	for i, conn := range connections {
		if validitySlice[i] {
			if foundValid < index {
				foundValid++
			} else {
				return conn
			}
		}
	}
	return connections[0]
}

func getDistSquared(loc1, loc2 cartesian.Vector2) float64 {

	return (loc1.X-loc2.X)*(loc1.X-loc2.X) + (loc1.Y-loc2.Y)*(loc1.Y-loc2.Y)

}
