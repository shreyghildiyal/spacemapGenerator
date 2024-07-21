package cartesian_test

import (
	"testing"

	"github.com/shreyghildiyal/spacemapGenerator/cartesian"
)

func TestGetUniquepoints(t *testing.T) {

	vectors := []cartesian.Vector2{
		{
			X: 1.0,
			Y: 2.0,
		},
		{
			X: 1.0,
			Y: 2.0,
		},
		{
			X: 1.1,
			Y: 2.0,
		},
	}

	uniqueVecs := cartesian.GetUniquepoints(vectors)

	if len(uniqueVecs) != 2 {
		t.Error("The vectors are not unique")
	}
}
