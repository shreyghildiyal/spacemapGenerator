package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/shreyghildiyal/spacemapGenerator/mapGen"
)

type Game struct {
	stars          []mapGen.Star
	clusterColours []color.RGBA
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Hello, World!, %d, %d\n", screen.Bounds().Max.X, screen.Bounds().Max.Y), 0, 0)

	debugY := 20

	for _, star := range g.stars {
		if star.IsClusterCore {
			// vector.DrawFilledCircle(screen, float32(star.X), float32(star.Y), 2, RED, true)
			vector.DrawFilledRect(screen, float32(star.X)-4, float32(star.Y)-4, 8, 8, g.clusterColours[star.ClusterId], true)
			// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Id: %d, X: %f, Y: %f\n", star.Id, star.X, star.Y), 0, debugY)
			debugY += 20
		} else {
			vector.DrawFilledCircle(screen, float32(star.X), float32(star.Y), 2, g.clusterColours[star.ClusterId], true)
		}

	}

}

const HEIGHT = 800
const WIDTH = 800

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

func main() {
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Hello, World!")

	game := Game{}

	starCount := 600
	clusterCount := 25

	stars, err := mapGen.InitMap(mapGen.MapGenConfigs{
		MaxX:                    WIDTH,
		MaxY:                    HEIGHT,
		StarCount:               starCount,
		ClusterCount:            clusterCount,
		StarRepulsionFactor:     5,
		ClusterAttractionFactor: 0.50,
		Iterations:              100,
		InnerIterations:         50,
		Seed:                    5,
	})

	if err != nil {
		log.Fatal("There was an error in creating the stars", err.Error())
		return
	} else {
		game.stars = stars
	}

	err = mapGen.AddStarBoundaries(game.stars, WIDTH, HEIGHT)

	if err != nil {
		log.Fatal("There was an error in creating star boundaries", err.Error())
		return
	}

	game.clusterColours = mapGen.GetClusterColours(clusterCount)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
