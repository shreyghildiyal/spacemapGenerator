package mapGen_test

import (
	"testing"

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
