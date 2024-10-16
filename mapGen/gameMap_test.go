package mapGen_test

import (
	"testing"

	"github.com/shreyghildiyal/spacemapGenerator/mapGen"
)

func TestInitMap(t *testing.T) {

	starCount := 5
	clusterCount := 2

	const HEIGHT = 800
	const WIDTH = 800

	stars, err := mapGen.InitMap(mapGen.MapGenConfigs{
		MaxX:         WIDTH - 50,
		MaxY:         HEIGHT - 50,
		StarCount:    starCount,
		ClusterCount: clusterCount,
		// StarRepulsionFactor:     1,
		// ClusterAttractionFactor: 0.50,
		Iterations: 1,
		Seed:       5,
	})

	if err != nil {
		t.Errorf("there was an error generating the stars %s", err.Error())
	}

	for _, star := range stars {
		if star.ClusterId < 0 || star.ClusterId >= clusterCount {
			t.Errorf("The clusterid was invalid. Max is %d, Min is 0. Found %d", clusterCount-1, star.ClusterId)
		}
	}
}
