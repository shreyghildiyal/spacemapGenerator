package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/shreyghildiyal/spacemapGenerator/drawing"
	"github.com/shreyghildiyal/spacemapGenerator/mapGen"
)

type Game struct {
	Stars          []mapGen.Star
	ClusterColours []color.RGBA
	DummyImage     *ebiten.Image
	Height         int
	Width          int
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Hello, World!, %d, %d\n", screen.Bounds().Max.X, screen.Bounds().Max.Y), 0, 0)

	debugY := 20

	for _, star := range g.Stars {

		drawing.DrawDomain(star, g.ClusterColours, screen, g.DummyImage)

		drawing.DrawNeighbourLines(star, g.Stars, screen)

		drawing.DrawStar(star, screen, g.ClusterColours, debugY)

	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.Width, g.Height
}
