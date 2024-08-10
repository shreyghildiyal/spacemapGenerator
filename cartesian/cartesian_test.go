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

func TestSolveLinearEquations2(t *testing.T) {

	lhs := [][]float64{
		{0, 3},
		{2, 0},
	}
	rhs := []float64{
		15,
		10,
	}

	expected := []float64{
		5,
		5,
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

func TestPerpendicularVector(t *testing.T) {
	inputVector := cartesian.Vector2{
		X: 100,
		Y: 0,
	}
	outVector := cartesian.GetPerpendicularVector(inputVector)
	if outVector.X != 0 {
		t.Errorf("The returned vector is not perpendicular to input vector")
	}
	if cartesian.DotProduct(inputVector, outVector) != 0 {
		t.Errorf("The returned vector is not perpendicular to input vector")
	}
}

func TestPerpendicularVector2(t *testing.T) {
	inputVector := cartesian.Vector2{
		X: 100,
		Y: 100,
	}
	outVector := cartesian.GetPerpendicularVector(inputVector)

	if cartesian.DotProduct(inputVector, outVector) != 0 {
		t.Errorf("The returned vector is not perpendicular to input vector")
	}
}

func TestGetIntersectionPoint(t *testing.T) {

	bisectingLine := cartesian.Line2D{
		Anchor: cartesian.Vector2{
			X: 150,
			Y: 150,
		},
		Direction: cartesian.Vector2{
			X: 0,
			Y: 100,
		},
	}

	boundaryLine := cartesian.Line2D{
		Anchor: cartesian.Vector2{
			X: 0,
			Y: 0,
		},
		Direction: cartesian.Vector2{
			X: 800,
			Y: 0,
		},
	}

	intersectionPoint, multi1, multi2, err := cartesian.GetIntersectionPoint(bisectingLine, boundaryLine)

	if err != nil {
		t.Error("there was an error in getting intersection point")
	}

	if intersectionPoint.X != 150 || intersectionPoint.Y != 0 {
		t.Error("the intersection point was incorrect")
	}

	if multi1 != -1.5 {
		t.Errorf("the multi for bisecting line was wrong. Found %f", multi1)
	}

	if multi2 != 150.0/800.0 {
		t.Errorf("the boundary line segment multi is wrong. Found %f. Expected %f", multi2, 150.0/800.0)
	}
}

func TestGetBisectingLine(t *testing.T) {
	point1 := cartesian.Vector2{
		X: 100,
		Y: 100,
	}

	point2 := cartesian.Vector2{
		X: 200,
		Y: 200,
	}

	bisectingLine := cartesian.GetBisectingLine(point1, point2)

	if bisectingLine.Anchor.X != 150 {
		t.Errorf("Anchor.X should be 150 but was %f", bisectingLine.Anchor.X)
	}

	if bisectingLine.Anchor.Y != 150 {
		t.Errorf("Anchor.Y should be 150 but was %f", bisectingLine.Anchor.X)
	}

}
