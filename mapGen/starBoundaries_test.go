package mapGen_test

import (
	"fmt"
	"testing"

	"github.com/shreyghildiyal/spacemapGenerator/cartesian"
	"github.com/shreyghildiyal/spacemapGenerator/mapGen"
)

func TestBoundaryGeneration(t *testing.T) {

	stars := []mapGen.Star{}

	maxX := 800.0
	maxY := 800.0

	// starCount := 10

	coords := [][]float64{
		{643.0758104437056, 415.6953946492434},
		{783.4904284601242, 477.39919489765884},
		{383.4979574598278, 560.2134529598569},
		{375.1536691059741, 686.3475123107872},
		{239.2233540566026, 795.2371977779281},
		{103.01443822191032, 592.2935856308802},
		{435.491394644248, 104.26394873857046},
		{397.26426418066836, 237.66527512170967},
		{502.21563816738455, 329.6547916632327},
		{103.01443822191041, 592.2935856308803},
	}

	for i := 0; i < len(coords); i++ {
		star := mapGen.Star{
			Id: i,
		}
		star.X = coords[i][0]
		star.Y = coords[i][1]
		stars = append(stars, star)
	}

	mapGen.AddStarBoundaries(stars, maxX, maxY)
}

func TestBoundaryGeneration2(t *testing.T) {

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
		{
			X: 600,
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

	mapGen.AddStarBoundaries(stars, 800, 800)

	for _, star := range stars {
		if len(star.BoundaryCorners) < 3 {
			t.Errorf("The number of boundary corners for star with coords %f, %f is too low. Found %d", star.X, star.Y, len(star.BoundaryCorners))
		}
	}
}

func TestDummy(t *testing.T) {

	lines := []cartesian.Line2D{
		{Anchor: cartesian.Vector2{X: 282.2679177162924, Y: 800}, Direction: cartesian.Vector2{X: 66.48150595175082, Y: -75.09348346852596}},
		{Anchor: cartesian.Vector2{X: 478.49607524586384, Y: 799.9999999999995}, Direction: cartesian.Vector2{X: 13.534359392069462, Y: -96.54210481923121}},
		{Anchor: cartesian.Vector2{X: 478.49607524586384, Y: 799.9999999999995}, Direction: cartesian.Vector2{X: -196.22815752957138, Y: 4.547473508864641e-13}},
		{Anchor: cartesian.Vector2{X: 348.7494236680432, Y: 724.906516531474}, Direction: cartesian.Vector2{X: 143.28101096989008, Y: -21.448621350705707}},
	}

	for _, line := range lines {
		fmt.Println("Anchor: ", line.Anchor)
		fmt.Println("EndPoint: ", line.EndPoint())
	}

	points, err := mapGen.GetOrderedBoundaryPoints(lines)

	if err != nil {
		t.Errorf("There was an error")
	}

	if len(points) != 4 {
		t.Errorf("The number of points is incorrect. Expected %d found %d", 4, len(points))
	}
}
