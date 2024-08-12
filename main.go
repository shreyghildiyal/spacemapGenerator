package main

import (
	"fmt"
	"image/color"
	"log"

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

		drawNeighbourLines(star, screen)

		drawStar(star, screen, g, debugY)

	}

}

func drawDomain(star mapGen.Star, g *Game, screen *ebiten.Image) {
	cornerCount := len(star.BoundaryCorners)

	fmt.Println("StarId", star.Id, "corner count", cornerCount)
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
	if star.IsClusterCore {

		vector.DrawFilledRect(screen, float32(star.X)-4, float32(star.Y)-4, 8, 8, g.clusterColours[star.ClusterId], true)

		debugY += 20
	} else {
		vector.DrawFilledCircle(screen, float32(star.X), float32(star.Y), 2, g.clusterColours[star.ClusterId], true)
	}
}

func drawNeighbourLines(star mapGen.Star, screen *ebiten.Image) {
	for i := range star.Neighbours {
		if star.Id < star.Neighbours[i].Id {
			vector.StrokeLine(screen, float32(star.Neighbours[i].X), float32(star.Neighbours[i].Y), float32(star.X), float32(star.Y), 1, color.White, false)
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

	if testMode == "ALL" {
		starGeneration(game)
	} else if testMode == "DEV" {
		dev(game)
	}
}

func boundaryGenTest(game Game) {

	star1 := mapGen.Star{
		Vector2: cartesian.Vector2{
			X: 100,
			Y: 100,
		},
		Id:            0,
		ClusterId:     0,
		IsClusterCore: true,
	}

	game.stars = append(game.stars, star1)

	star2 := mapGen.Star{
		Vector2: cartesian.Vector2{
			X: 200,
			Y: 200,
		},
		Id:            1,
		ClusterId:     1,
		IsClusterCore: true,
	}

	game.stars = append(game.stars, star2)

	err := mapGen.AddStarBoundaries(game.stars, WIDTH, HEIGHT)
	// err = mapGen.AddDummyStarBoundaries(game.stars, WIDTH, HEIGHT)
	mapGen.AddDummyNeighbours(game.stars)

	if err != nil {
		log.Fatal("There was an error in creating star boundaries", err.Error())
		return
	}

	game.clusterColours = mapGen.GetClusterColours(2)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

func dev(game Game) {

	clusterCount := 5

	coords := []cartesian.Vector2{
		{
			X: 200,
			Y: 200,
		},
		{
			X: 100,
			Y: 200,
		},
		{
			X: 200,
			Y: 100,
		},
		{
			X: 300,
			Y: 200,
		},
		{
			X: 200,
			Y: 300,
		},
	}

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

	err := mapGen.AddStarBoundaries(game.stars, WIDTH, HEIGHT)
	fmt.Println(stars)
	// err = mapGen.AddDummyStarBoundaries(game.stars, WIDTH, HEIGHT)
	mapGen.AddDummyNeighbours(stars)

	game.stars = stars

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
