package mapGen

import (
	"image/color"
	"math/rand"
)

type StarConnection struct {
	FromStar int
	ToStar   int
}

type Cluster struct {
	Id                int
	Stars             []int
	NeighbourClusters []int
	Boundaries        map[int]int // key: starId within cluster, val: starId in another cluster

}

func GetClusterColours(clusterCount int) []color.RGBA {
	colors := []color.RGBA{}

	for range clusterCount {
		color := color.RGBA{
			R: uint8(100 + rand.Int()%156),
			G: uint8(100 + rand.Int()%156),
			B: uint8(100 + rand.Int()%156),
			A: 100,
		}
		colors = append(colors, color)
	}
	return colors
}
