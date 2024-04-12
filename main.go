package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	mapgen "github.com/shreyghildiyal/spacemapGenerator/mapGen"
)

type Game struct {
	stars []mapgen.Location
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")

	var RED = color.RGBA{
		R: 200,
		G: 0,
		B: 0,
		A: 0,
	}

	for _, star := range g.stars {
		if star.IsClusterCore {
			// vector.DrawFilledCircle(screen, float32(star.X), float32(star.Y), 2, RED, true)
			vector.DrawFilledRect(screen, float32(star.X)-4, float32(star.Y)-4, 8, 8, RED, true)
		} else {
			vector.DrawFilledCircle(screen, float32(star.X), float32(star.Y), 2, color.White, true)
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

	stars, _, err := mapgen.InitMap(mapgen.MapGenConfigs{
		MaxX:                    WIDTH,
		MaxY:                    HEIGHT,
		StarCount:               600,
		ClusterCount:            20,
		StartRepulsionFactor:    50,
		ClusterAttractionFactor: 0.95,
	})
	if err != nil {
		log.Fatal("There was an error in creating the stars", err.Error())
		return
	} else {
		game.stars = stars
	}
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
