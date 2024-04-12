package mapGen

import (
	"errors"

	"image/color"
	"math/rand"
)

type Location struct {
	Id            int
	X             float64
	Y             float64
	ClusterId     int
	IsClusterCore bool
}

type Cluster struct {
	Core  Location
	color color.RGBA
}

type MapGenConfigs struct {
	MaxX                    float64
	MaxY                    float64
	StarCount               int
	ClusterCount            int
	StartRepulsionFactor    float64
	ClusterAttractionFactor float64
}

// var gameMap []Location

func InitMap(configs MapGenConfigs) ([]Location, []Cluster, error) {

	if configs.ClusterCount > configs.StarCount {
		return nil, nil, errors.New("the cluster count cant be greater than star count")
	}

	stars := initStarsAtRandomLocations(configs.StarCount, configs.MaxX, configs.MaxY)

	// we will assume that the first N stars are cluster cores.
	for i := 0; i < configs.ClusterCount; i++ {
		stars[i].IsClusterCore = true
	}

	distanceGrid := getDistanceGrid(configs, stars)

	// assign stars to clusters
	updateStarClusters(stars, distanceGrid, configs.ClusterCount)
	// movestars away from each other
	// move stars towards cluster core

	return stars, nil, nil
}

func updateStarClusters(stars []Location, distanceGrid [][]float64, clusterCount int) {

	for i, star := range stars {
		if i < clusterCount {
			star.ClusterId = i
		} else {
			closestClusterCoreId := getClosestClusterCoreId(star, distanceGrid, stars[:clusterCount])
			star.ClusterId = closestClusterCoreId
		}
	}
}

func getClosestClusterCoreId(star Location, distanceGrid [][]float64, clusterCores []Location) int {
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

func getDistanceGrid(configs MapGenConfigs, stars []Location) [][]float64 {
	distanceGrid := make([][]float64, configs.StarCount)
	for i := 0; i < configs.StarCount; i++ {
		distanceGrid[i] = make([]float64, configs.StarCount)
	}

	for i := 0; i < configs.StarCount; i++ {
		for j := i; j < configs.StarCount; j++ {
			distanceGrid[i][j] = getDistSqr(stars[i], stars[j])
			distanceGrid[j][i] = distanceGrid[i][j]
		}
	}
	return distanceGrid
}

func getDistSqr(location1, location2 Location) float64 {
	dx := location1.X - location2.X
	dy := location1.Y - location2.Y

	return dx*dx + dy*dy
}

func initStarsAtRandomLocations(starCount int, maxX float64, maxY float64) []Location {
	locations := make([]Location, starCount)

	for i := 0; i < starCount; i++ {

		loc := Location{
			X:             rand.Float64() * maxX,
			Y:             rand.Float64() * maxY,
			ClusterId:     -1,
			IsClusterCore: false,
		}

		locations[i] = loc
	}
	return locations
}
