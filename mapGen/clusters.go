package mapGen

import (
	"image/color"
	"math/rand"
)

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
