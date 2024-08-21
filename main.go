package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/shreyghildiyal/spacemapGenerator/cartesian"
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

		drawDomain(star, g, screen)

		drawNeighbourLines(star, g.stars, screen)

		drawStar(star, screen, g, debugY)

	}

}

func drawDomain(star mapGen.Star, g *Game, screen *ebiten.Image) {
	cornerCount := len(star.BoundaryCorners)

	// fmt.Println("StarId", star.Id, "corner count", cornerCount)
	if cornerCount > 0 {

		op := &ebiten.DrawTrianglesOptions{
			Blend: ebiten.Blend{},
		}

		vertices := []ebiten.Vertex{}
		indices := []uint16{}

		for i := 0; i < cornerCount; i++ {
			vertex := ebiten.Vertex{
				DstX:   float32(star.BoundaryCorners[i].X),
				DstY:   float32(star.BoundaryCorners[i].Y),
				ColorR: float32(g.clusterColours[star.ClusterId].R) / 255,
				ColorG: float32(g.clusterColours[star.ClusterId].G) / 255,
				ColorB: float32(g.clusterColours[star.ClusterId].B) / 255,
				ColorA: float32(g.clusterColours[star.ClusterId].A) / 255,
			}

			vertices = append(vertices, vertex)

			indices = append(indices, uint16((i)%cornerCount), uint16((i+1)%cornerCount), uint16(cornerCount))
		}

		vertices = append(vertices, ebiten.Vertex{
			DstX:   float32(star.X),
			DstY:   float32(star.Y),
			ColorR: float32(g.clusterColours[star.ClusterId].R) / 255,
			ColorG: float32(g.clusterColours[star.ClusterId].G) / 255,
			ColorB: float32(g.clusterColours[star.ClusterId].B) / 255,
			ColorA: float32(g.clusterColours[star.ClusterId].A) / 255,
		})

		screen.DrawTriangles(vertices, indices, g.dummyImage, op)
	}
}

func drawStar(star mapGen.Star, screen *ebiten.Image, g *Game, debugY int) {
	// fmt.Println("Drawing star", star.Id)
	if star.IsClusterCore {

		vector.DrawFilledRect(screen, float32(star.X)-4, float32(star.Y)-4, 8, 8, g.clusterColours[star.ClusterId], true)

		debugY += 20
	} else {
		vector.DrawFilledCircle(screen, float32(star.X), float32(star.Y), 2, g.clusterColours[star.ClusterId], true)
	}
}

func drawNeighbourLines(star mapGen.Star, stars []mapGen.Star, screen *ebiten.Image) {
	for _, neighbourId := range star.Neighbours {
		if star.Id < neighbourId {
			nStar := stars[neighbourId]
			vector.StrokeLine(screen, float32(nStar.X), float32(nStar.Y), float32(star.X), float32(star.Y), 1, color.White, false)
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

	game.dummyImage = ebiten.NewImage(3, 3)
	// game.dummyImage.Fill(color.White)
	game.dummyImage.Fill(color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 0,
	})

	testMode := "DEV"
	// testMode := "ALL"
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

func dev(game Game) (Game, error) {

	coords := []cartesian.Vector2{}
	starCount := 16

	var randGen *rand.Rand = rand.New(rand.NewSource(100))
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
			ClusterId:     i,
			IsClusterCore: true,
		}
		stars = append(stars, star)
	}
	game.stars = stars

	err := mapGen.AddStarBoundaries(game.stars, WIDTH, HEIGHT)
	// fmt.Println(stars)
	// err = mapGen.AddDummyStarBoundaries(game.stars, WIDTH, HEIGHT)
	mapGen.AddDummyNeighbours(stars)

	if err != nil {
		log.Fatal("There was an error in creating star boundaries", err.Error())
		return Game{}, fmt.Errorf("Error creating star boundaries %w", err)
	}

	game.clusterColours = mapGen.GetClusterColours(clusterCount)

	return game, nil
}

func starGeneration(game Game) (Game, error) {
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
		return Game{}, fmt.Errorf("Error in star creation %w", err)
	} else {
		game.stars = stars
	}

	err = mapGen.AddStarBoundaries(game.stars, WIDTH, HEIGHT)

	if err != nil {
		log.Fatal("There was an error in creating star boundaries", err.Error())
		return Game{}, fmt.Errorf("Error in adding star boundaries %w", err)
	}

	game.clusterColours = mapGen.GetClusterColours(clusterCount)

	mapGen.AddDummyNeighbours(game.stars)

	return game, nil
}

// 	if err := ebiten.RunGame(&game); err != nil {
// 		log.Fatal(err)
// 	}
// }
