package mapGen

import (
	"errors"
	"fmt"
	"math"

	"math/rand"

	"github.com/shreyghildiyal/spacemapGenerator/cartesian"
)

type Segment struct {
}

// type Location struct {
// 	X float64
// 	Y float64
// }

type MapGenConfigs struct {
	MaxX                    float64
	MaxY                    float64
	StarCount               int
	ClusterCount            int
	MaxStarCountRatio       float64
	StarRepulsionFactor     float64
	ClusterAttractionFactor float64
	Iterations              int
	Seed                    int64
	InnerIterations         int
}

var randGen *rand.Rand

func InitMap(configs MapGenConfigs) ([]Star, error) {

	randGen = rand.New(rand.NewSource(configs.Seed))

	if configs.ClusterCount > configs.StarCount {
		return nil, errors.New("the cluster count cant be greater than star count")
	}

	bestStars := []Star{}
	bestStandardDeviation := -1.0

	j := 0

	for j < configs.Iterations {

		stars := initStarsAtRandomLocations(configs.StarCount, configs.MaxX, configs.MaxY)

		// we will assume that the first N stars are cluster cores.
		for i := 0; i < configs.ClusterCount; i++ {
			stars[i].IsClusterCore = true
		}

		distanceGrid := getDistanceGrid(configs, stars)
		// movestars away from each other
		moveStarsAwayFromEachOther(stars, distanceGrid, configs.StarRepulsionFactor)
		moveStarsAwayFromEdges(stars, configs.StarRepulsionFactor, configs.MaxX, configs.MaxY)
		for i := 0; i < configs.InnerIterations; i++ {

			// assign stars to clusters
			updateStarClusters(stars, distanceGrid, configs.ClusterCount)

			// move clusterCore
			moveClusterCore(stars, configs.ClusterCount)
		}

		updateStarClusters(stars, distanceGrid, configs.ClusterCount)

		clusterStarCounts := getClusterStarCounts(stars)

		deviation := getDeviation(clusterStarCounts)

		// fmt.Println("Standard Deviation:", deviation)
		minClusterStarCount := getMinClusterStarCount(clusterStarCounts)
		maxClusterStarCount := getMaxClusterStarCount(clusterStarCounts)
		ratio := float64(maxClusterStarCount) / float64(minClusterStarCount)

		if ratio >= configs.MaxStarCountRatio {

			if len(bestStars) < configs.StarCount || deviation < float64(bestStandardDeviation) {
				bestStars = stars
				bestStandardDeviation = deviation
				fmt.Println("Best standard deviation in cluster star count: ", bestStandardDeviation, j, configs.Iterations)

			}
			j++

		}

	}
	return bestStars, nil
}

func getMaxClusterStarCount(clusterStarCounts []int) int {
	max := clusterStarCounts[0]
	for i := 1; i < len(clusterStarCounts); i++ {
		if clusterStarCounts[i] > max {
			max = clusterStarCounts[i]
		}
	}
	return max
}

func getMinClusterStarCount(clusterStarCounts []int) int {

	min := clusterStarCounts[0]
	for i := 1; i < len(clusterStarCounts); i++ {
		if clusterStarCounts[i] < min {
			min = clusterStarCounts[i]
		}
	}
	return min
}

func getDeviation(clusterStarCounts []int) float64 {
	mean := getMean(clusterStarCounts)

	variance := getVariance(clusterStarCounts, mean)

	return math.Sqrt(variance)
}

func getVariance(clusterStarCounts []int, mean float64) float64 {
	variance := 0.0
	for _, count := range clusterStarCounts {
		diff := mean - float64(count)
		variance += diff * diff
	}
	variance = variance / float64(len(clusterStarCounts))
	return variance
}

func getMean(clusterStarCounts []int) float64 {
	mean := 0.0
	for _, count := range clusterStarCounts {
		mean += float64(count)
	}
	mean = mean / float64(len(clusterStarCounts))
	return mean
}

func getClusterStarCounts(stars []Star) []int {

	counts := []int{}

	for _, star := range stars {
		if star.ClusterId >= len(counts) {
			counts = append(counts, make([]int, 1+star.ClusterId-len(counts))...)
		}
		counts[star.ClusterId]++
	}
	return counts
}

func moveStarsAwayFromEdges(stars []Star, repulsion, maxX, maxY float64) {
	for i := range stars {
		var dx float64 = 0
		var dy float64 = 0

		if stars[i].X <= 0 {
			stars[i].X = 1
		}
		if stars[i].X >= maxX {
			stars[i].X = maxX - 1
		}
		if stars[i].Y <= 0 {
			stars[i].Y = 1
		}
		if stars[i].Y >= maxY {
			stars[i].Y = maxY - 1
		}

		// repel from left

		dx += 1 / (stars[i].X * stars[i].X)

		// repel from right

		dist := maxX - stars[i].X
		dx -= 1 / (dist * dist)

		//repel from top

		dy += 1 / (stars[i].Y * stars[i].Y)

		// repel from bottom

		dist = maxY - stars[i].Y
		dy -= 1 / (dist * dist)

		// x := stars[i].X
		// y := stars[i].Y

		// fmt.Printf("X: %f, Y: %f, dX: %f, dY: %f", x, y, dx, dy)

		stars[i].X += dx * repulsion
		stars[i].Y += dy * repulsion
	}
}

func moveClusterCore(stars []Star, clusterCount int) {
	newLocArr := make([]Star, clusterCount)
	starsInCluster := make([]int, clusterCount)
	for i := 0; i < len(stars); i++ {
		clusterId := stars[i].ClusterId
		newLocArr[clusterId].X += stars[i].X
		newLocArr[clusterId].Y += stars[i].Y
		starsInCluster[clusterId]++
	}

	for i := range newLocArr {
		stars[i].X = newLocArr[i].X / float64(starsInCluster[i])
		stars[i].Y = newLocArr[i].Y / float64(starsInCluster[i])
	}
}

func moveStarsAwayFromEachOther(stars []Star, distanceGrid [][]float64, repulsion float64) {

	for i := range stars {
		var dx float64 = 0
		var dy float64 = 0
		for j := range stars {
			if i != j {
				distance := distanceGrid[i][j]
				repulsionX := stars[i].X - stars[j].X
				repulsionY := stars[i].Y - stars[j].Y

				dx += (repulsionX / (distance * distance * distance))
				dy += (repulsionY / (distance * distance * distance))
			}
		}

		stars[i].X += dx * repulsion
		stars[i].Y += dy * repulsion
	}

}

func updateStarClusters(stars []Star, distanceGrid [][]float64, clusterCount int) {

	for i := range stars {
		if i < clusterCount {
			stars[i].ClusterId = i
		} else {
			closestClusterCoreId := getClosestClusterCoreId(stars[i], distanceGrid, stars[:clusterCount])
			stars[i].ClusterId = closestClusterCoreId
		}
	}
}

func getClosestClusterCoreId(star Star, distanceGrid [][]float64, clusterCores []Star) int {
	clusterId := clusterCores[0].Id
	minDist := distanceGrid[star.Id][0]

	for _, coreStar := range clusterCores {
		dist := distanceGrid[star.Id][coreStar.Id]
		if dist < minDist {
			minDist = dist
			clusterId = coreStar.Id
		}
	}
	return clusterId
}

func getDistanceGrid(configs MapGenConfigs, stars []Star) [][]float64 {
	distanceGrid := make([][]float64, configs.StarCount)
	for i := 0; i < configs.StarCount; i++ {
		distanceGrid[i] = make([]float64, configs.StarCount)
	}

	for i := 0; i < configs.StarCount; i++ {
		for j := i; j < configs.StarCount; j++ {
			distanceGrid[i][j] = getDist(stars[i], stars[j])
			distanceGrid[j][i] = distanceGrid[i][j]
		}
	}
	return distanceGrid
}

func getDist(location1, location2 Star) float64 {
	dx := location1.X - location2.X
	dy := location1.Y - location2.Y

	return math.Sqrt(dx*dx + dy*dy)
}

func initStarsAtRandomLocations(starCount int, maxX float64, maxY float64) []Star {
	locations := make([]Star, starCount)

	for i := 0; i < starCount; i++ {

		loc := Star{
			Vector2: cartesian.Vector2{
				X: randGen.Float64() * maxX,
				Y: randGen.Float64() * maxY,
			},

			ClusterId:     -1,
			IsClusterCore: false,
			Id:            i,
		}

		locations[i] = loc
	}
	return locations
}
