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
	MaxX              float64
	MaxY              float64
	StarCount         int
	ClusterCount      int
	MaxStarCountRatio float64
	// StarRepulsionFactor     float64
	// ClusterAttractionFactor float64
	Iterations      int
	Seed            int64
	InnerIterations int
	MinStarDistance float64
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

		stars := initStarsAtRandomLocations(configs.StarCount, configs.MaxX, configs.MaxY, configs.MinStarDistance)

		// we will assume that the first N stars are cluster cores.
		for i := 0; i < configs.ClusterCount; i++ {
			stars[i].IsClusterCore = true
		}

		// distanceGrid := getDistanceGrid(stars)
		// movestars away from each other
		// forceStarRepulsion(stars, configs.StarRepulsionFactor, configs.MaxX, configs.MaxY)

		for i := 0; i < configs.InnerIterations; i++ {

			// assign stars to clusters
			updateStarClusters(stars, configs.ClusterCount)

			// move clusterCore
			moveClusterCore(stars, configs.ClusterCount)
		}

		updateStarClusters(stars, configs.ClusterCount)

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

func updateStarClusters(stars []Star, clusterCount int) {

	for i := range stars {
		if i < clusterCount {
			stars[i].ClusterId = i
		} else {
			closestClusterCoreId := getClosestClusterCoreId(stars[i], stars[:clusterCount])
			stars[i].ClusterId = closestClusterCoreId
		}
	}
}

func getClosestClusterCoreId(star Star, clusterCores []Star) int {
	clusterId := clusterCores[0].Id
	// minDist := distanceGrid[star.Id][0]
	minDist := getDist(star, clusterCores[0])

	for _, coreStar := range clusterCores {
		dist := getDist(star, coreStar)
		if dist < minDist {
			minDist = dist
			clusterId = coreStar.Id
		}
	}
	return clusterId
}

func getDist(location1, location2 Star) float64 {
	dx := location1.X - location2.X
	dy := location1.Y - location2.Y

	return math.Sqrt(dx*dx + dy*dy)
}

func initStarsAtRandomLocations(starCount int, maxX, maxY, minDistance float64) []Star {
	locations := []cartesian.Vector2{}

	minX := maxX * 0.03
	minY := maxY * 0.03

	randSizeX := maxX * 0.94
	randSizeY := maxY * 0.94

	failedCount := 0

	for len(locations) < starCount {

		loc := cartesian.Vector2{
			X: minX + randGen.Float64()*randSizeX,
			Y: minY + randGen.Float64()*randSizeY,
		}

		if isValidNewLocation(loc, locations, minDistance) {
			locations = append(locations, loc)
			failedCount = 0
		} else {
			failedCount++

			if failedCount > 500 {
				locations = []cartesian.Vector2{}
				logger.Println("Faced too many failures generating a new valid location. Resetting")
			}
		}
	}

	// stars := make([]cartesian.Vector2, starCount)

	return initStarsAtLocations(locations)
}

func isValidNewLocation(newLoc cartesian.Vector2, locations []cartesian.Vector2, minDistance float64) bool {
	minDistSquared := minDistance * minDistance
	for _, loc := range locations {
		dist := getDistSquared(loc, newLoc)
		if dist < minDistSquared {
			return false
		}
	}
	return true
}

func initStarsAtLocations(locs []cartesian.Vector2) []Star {

	stars := make([]Star, len(locs))
	for i, loc := range locs {
		star := Star{
			Vector2: loc,

			ClusterId:     -1,
			IsClusterCore: false,
			Id:            i,
		}
		stars[i] = star
	}
	return stars
}
