package mapGen_test

import (
	"log"
	"math/rand"
	"testing"

	"github.com/shreyghildiyal/spacemapGenerator/cartesian"
	"github.com/shreyghildiyal/spacemapGenerator/mapGen"
)

func TestNeighbours(t *testing.T) {

	coords := []cartesian.Vector2{}
	starCount := 10

	WIDTH := 800.0
	HEIGHT := 800.0

	var randGen *rand.Rand = rand.New(rand.NewSource(5))
	for i := 0; i < starCount; i++ {
		coords = append(coords, cartesian.Vector2{
			X: randGen.Float64() * WIDTH,
			Y: randGen.Float64() * HEIGHT,
		})
	}

	// clusterCount := len(coords)

	stars := []mapGen.Star{}

	for i, coord := range coords {
		star := mapGen.Star{
			Vector2:       coord,
			Id:            i,
			ClusterId:     0,
			IsClusterCore: true,
		}
		stars = append(stars, star)
	}
	// gameObj.Stars = stars

	err := mapGen.AddStarBoundaries(stars, WIDTH, HEIGHT)
	// fmt.Println(stars)
	// err = mapGen.AddDummyStarBoundaries(game.stars, WIDTH, HEIGHT)
	mapGen.AddStarNeighbours(stars)

	if err != nil {
		log.Fatal("There was an error in creating star boundaries", err.Error())

	}

	// gameObj.ClusterColours = mapGen.GetClusterColours(clusterCount)

	// return gameObj, nil
}
