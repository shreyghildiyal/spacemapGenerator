package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/shreyghildiyal/spacemapGenerator/cartesian"
	"github.com/shreyghildiyal/spacemapGenerator/game"
	"github.com/shreyghildiyal/spacemapGenerator/mapGen"
)

const HEIGHT = 800
const WIDTH = 800

func main() {

	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Hello, World!")

	game := game.Game{
		Height: HEIGHT,
		Width:  WIDTH,
	}

	game.DummyImage = ebiten.NewImage(3, 3)
	// game.dummyImage.Fill(color.White)
	game.DummyImage.Fill(color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 0,
	})

	// testMode := "DEV"
	testMode := "ALL"
	var err error = nil

	if testMode == "ALL" {
		game, err = starGeneration(game)
	} else if testMode == "DEV" {
		game, err = dev(game)
	}

	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

func dev(gameObj game.Game) (game.Game, error) {

	coords := []cartesian.Vector2{}
	starCount := 10

	var randGen *rand.Rand = rand.New(rand.NewSource(5))
	for i := 0; i < starCount; i++ {
		coords = append(coords, cartesian.Vector2{
			X: randGen.Float64() * WIDTH,
			Y: randGen.Float64() * HEIGHT,
		})
	}

	clusterCount := len(coords)

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
	gameObj.Stars = stars

	err := mapGen.AddStarBoundaries(gameObj.Stars, WIDTH, HEIGHT)
	// fmt.Println(stars)
	// err = mapGen.AddDummyStarBoundaries(game.stars, WIDTH, HEIGHT)
	mapGen.AddStarNeighbours(stars)

	if err != nil {
		log.Fatal("There was an error in creating star boundaries", err.Error())
		return game.Game{}, fmt.Errorf("Error creating star boundaries %w", err)
	}

	gameObj.ClusterColours = mapGen.GetClusterColours(clusterCount)

	return gameObj, nil
}

func starGeneration(gameObj game.Game) (game.Game, error) {
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Hello, World!")

	starCount := 1000
	clusterCount := 30
	maxStarCountRatio := 1.5

	stars, err := mapGen.InitMap(mapGen.MapGenConfigs{
		MaxX:                    WIDTH,
		MaxY:                    HEIGHT,
		StarCount:               starCount,
		ClusterCount:            clusterCount,
		MaxStarCountRatio:       maxStarCountRatio,
		StarRepulsionFactor:     50,
		ClusterAttractionFactor: 0.50,
		Iterations:              500,
		InnerIterations:         50,
		Seed:                    5,
	})

	if err != nil {
		log.Fatal("There was an error in creating the stars", err.Error())

		return game.Game{}, fmt.Errorf("Error in star creation %w", err)
	} else {
		gameObj.Stars = stars
	}

	err = mapGen.AddStarBoundaries(gameObj.Stars, WIDTH, HEIGHT)

	if err != nil {
		log.Fatal("There was an error in creating star boundaries", err.Error())
		return game.Game{}, fmt.Errorf("Error in adding star boundaries %w", err)
	}

	gameObj.ClusterColours = mapGen.GetClusterColours(clusterCount)

	// mapGen.AddDummyNeighbours(gameObj.Stars)
	mapGen.AddStarNeighbours(stars)

	return gameObj, nil
}
