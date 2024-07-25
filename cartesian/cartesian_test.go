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

func TestSolveLinearEquations(t *testing.T) {

	lhs := [][]float64{
		{5, 5},
		{2, 3},
	}
	rhs := []float64{
		10,
		5,
	}

	expected := []float64{
		1,
		1,
	}

	results, err := cartesian.SolveLinearEquations(lhs, rhs)

	if err != nil {
		t.Error("there should have been no errors")
	}

	if expected[0] != results[0] {
		t.Error("first val doesnt match", expected[0], results[0])
	}

	if expected[1] != results[1] {
		t.Error("second val doesnt match", expected[1], results[1])
	}

}
