package main

import (
	"fmt"
	"image"
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
	dummyImage     *ebiten.Image
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Hello, World!, %d, %d\n", screen.Bounds().Max.X, screen.Bounds().Max.Y), 0, 0)

	debugY := 20

	for _, star := range g.stars {

		cornerCount := len(star.BoundaryCorners)
		// cornerCount := -1

		if cornerCount > 0 {
			// println("Drawing boundaries")
			op := &ebiten.DrawTrianglesOptions{
				Blend: ebiten.Blend{},
			}

			vertices := []ebiten.Vertex{}
			indices := []uint16{}

			for i := 0; i < cornerCount; i++ {
				vertex := ebiten.Vertex{
					DstX:   float32(star.BoundaryCorners[i].X),
					DstY:   float32(star.BoundaryCorners[i].Y),
					ColorR: float32(g.clusterColours[star.ClusterId].R),
					ColorG: float32(g.clusterColours[star.ClusterId].G),
					ColorB: float32(g.clusterColours[star.ClusterId].B),
					ColorA: float32(g.clusterColours[star.ClusterId].A),
				}

				vertices = append(vertices, vertex)

				indices = append(indices, uint16((i)%cornerCount), uint16((i+1)%cornerCount), uint16(cornerCount))
			}

			vertices = append(vertices, ebiten.Vertex{
				DstX:   float32(star.X),
				DstY:   float32(star.Y),
				ColorR: float32(g.clusterColours[star.ClusterId].R),
				ColorG: float32(g.clusterColours[star.ClusterId].G),
				ColorB: float32(g.clusterColours[star.ClusterId].B),
				ColorA: float32(g.clusterColours[star.ClusterId].A),
			})
			// println(indices)
			fmt.Println(indices)
			screen.DrawTriangles(vertices, indices, g.dummyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), op)
		}

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

	testMode := "BOUNDARYGEN"

	game := Game{}

	game.dummyImage = ebiten.NewImage(3, 3)
	// game.dummyImage.Fill(color.White)
	game.dummyImage.Fill(color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 100,
	})

	if testMode == "STARGEN" {
		starGeneration(game)
	} else if testMode == "BOUNDARYGEN" {
		boundaryGeneration(game)
	}
}

func boundaryGeneration(game Game) {
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Hello, World!")

	starCount := 100
	clusterCount := 90
	maxStarCountRatio := 1.5

	stars, err := mapGen.InitMap(mapGen.MapGenConfigs{
		MaxX:                    WIDTH,
		MaxY:                    HEIGHT,
		StarCount:               starCount,
		ClusterCount:            clusterCount,
		MaxStarCountRatio:       maxStarCountRatio,
		StarRepulsionFactor:     5,
		ClusterAttractionFactor: 0.50,
		Iterations:              500,
		InnerIterations:         50,
		Seed:                    5,
	})

	if err != nil {
		log.Fatal("There was an error in creating the stars", err.Error())
		return
	} else {
		game.stars = stars
	}

	// err = mapGen.AddStarBoundaries(game.stars, WIDTH, HEIGHT)
	err = mapGen.AddDummyStarBoundaries(game.stars, WIDTH, HEIGHT)

	if err != nil {
		log.Fatal("There was an error in creating star boundaries", err.Error())
		return
	}

	game.clusterColours = mapGen.GetClusterColours(clusterCount)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

func starGeneration(game Game) {
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Hello, World!")

	starCount := 100
	clusterCount := 10
	maxStarCountRatio := 1.5

	stars, err := mapGen.InitMap(mapGen.MapGenConfigs{
		MaxX:                    WIDTH,
		MaxY:                    HEIGHT,
		StarCount:               starCount,
		ClusterCount:            clusterCount,
		MaxStarCountRatio:       maxStarCountRatio,
		StarRepulsionFactor:     5,
		ClusterAttractionFactor: 0.50,
		Iterations:              500,
		InnerIterations:         50,
		Seed:                    5,
	})

	if err != nil {
		log.Fatal("There was an error in creating the stars", err.Error())
		return
	} else {
		game.stars = stars
	}

	// err = mapGen.AddStarBoundaries(game.stars, WIDTH, HEIGHT)

	// if err != nil {
	// 	log.Fatal("There was an error in creating star boundaries", err.Error())
	// 	return
	// }

	game.clusterColours = mapGen.GetClusterColours(clusterCount)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
