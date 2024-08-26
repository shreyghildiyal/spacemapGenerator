package mapGen

import (
	"fmt"

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
	InClusterConnectionRatio float32 // for each star in a cluster the number of connections to other cluster stars it has divided by the number of such possible connections
	ClusterConnectionRatio   float32 // for each cluster, number of connections to stars in other clusters/ number of such connections possible
}

func AddStarNeighbours(stars []Star) {
	clusterStarIds := getClusterStarIds(stars)
	for _, clusterStarIds := range clusterStarIds {
		connectClusterStars(stars, clusterStarIds)
	}

	connectClusters(stars)
}

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

func connectClusterStars(stars []Star, clusterStarIds []int) {
	// create a minimum spanning tree

	createMinSpanningTree(clusterStarIds, stars)

	// use a probability system to use the remaining borders
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

func connectClusters(stars []Star) {

}

func getDistSquared(loc1, loc2 cartesian.Vector2) float64 {

	return (loc1.X-loc2.X)*(loc1.X-loc2.X) + (loc1.Y-loc2.Y)*(loc1.Y-loc2.Y)

}
