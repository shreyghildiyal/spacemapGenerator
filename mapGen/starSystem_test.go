package mapGen_test

import (
	"testing"

	"github.com/shreyghildiyal/spacemapGenerator/mapGen"
)

func TestGenerationDistances(t *testing.T) {

	minDistance := 0
	maxDistance := 50
	desiredCount := 8

	distances := mapGen.GetHeavenlyBodyDistances(desiredCount, minDistance, maxDistance)

	if len(distances) != desiredCount {
		t.Error("the number of returned distances is wrong")
	}
}
