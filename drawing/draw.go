package drawing

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/shreyghildiyal/spacemapGenerator/mapGen"
)

func DrawDomain(star mapGen.Star, clusterColours []color.RGBA, screen *ebiten.Image, dummyImage *ebiten.Image) {
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
				ColorR: float32(clusterColours[star.ClusterId].R) / 255,
				ColorG: float32(clusterColours[star.ClusterId].G) / 255,
				ColorB: float32(clusterColours[star.ClusterId].B) / 255,
				ColorA: float32(clusterColours[star.ClusterId].A) / 255,
			}

			vertices = append(vertices, vertex)

			indices = append(indices, uint16((i)%cornerCount), uint16((i+1)%cornerCount), uint16(cornerCount))
		}

		vertices = append(vertices, ebiten.Vertex{
			DstX:   float32(star.X),
			DstY:   float32(star.Y),
			ColorR: float32(clusterColours[star.ClusterId].R) / 255,
			ColorG: float32(clusterColours[star.ClusterId].G) / 255,
			ColorB: float32(clusterColours[star.ClusterId].B) / 255,
			ColorA: float32(clusterColours[star.ClusterId].A) / 255,
		})

		screen.DrawTriangles(vertices, indices, dummyImage, op)
	}
}

func DrawStar(star mapGen.Star, screen *ebiten.Image, clusterColours []color.RGBA, debugY int) {
	// fmt.Println("Drawing star", star.Id)
	if star.IsClusterCore {

		vector.DrawFilledRect(screen, float32(star.X)-4, float32(star.Y)-4, 8, 8, clusterColours[star.ClusterId], true)

		debugY += 20
	} else {
		vector.DrawFilledCircle(screen, float32(star.X), float32(star.Y), 2, clusterColours[star.ClusterId], true)
	}
}

func DrawNeighbourLines(star mapGen.Star, stars []mapGen.Star, screen *ebiten.Image) {
	for _, neighbourId := range star.Neighbours {
		if star.Id < neighbourId {
			nStar := stars[neighbourId]
			vector.StrokeLine(screen, float32(nStar.X), float32(nStar.Y), float32(star.X), float32(star.Y), 1, color.White, false)
		}

	}
}
