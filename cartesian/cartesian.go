package cartesian

type Vector2 struct {
	X float64
	Y float64
}

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

type Line struct {
	Anchor    Vector2
	Direction Vector2
}

func GetPerpendicularVector(vec Vector2) Vector2 {
	rotatedVec3 := CrossProduct(Vector3{X: vec.X, Y: vec.Y, Z: 0}, Vector3{X: 0, Y: 0, Z: 1})

	return Vector2{X: rotatedVec3.X, Y: rotatedVec3.Y}
}

func CrossProduct(vector1, vector2 Vector3) Vector3 {
	newVec := Vector3{X: 0, Y: 0, Z: 0}

	newVec.Z += vector1.X * vector2.Y
	newVec.Z -= vector1.Y * vector2.X

	newVec.Y += vector1.Z * vector2.X
	newVec.Y -= vector1.X * vector2.Z

	newVec.X += vector1.Y * vector2.Z
	newVec.X -= vector1.Z * vector2.Y

	return newVec
}

func GetBisectingLine(point1, point2 Vector2) Line {
	midPoint := Vector2{
		X: (point1.X + point2.X) / 2,
		Y: (point1.Y + point2.Y) / 2,
	}

	perpendicularVector := GetPerpendicularVector(Vector2{X: point1.X - point2.X, Y: point1.Y - point2.Y})
	return Line{
		Anchor:    midPoint,
		Direction: perpendicularVector,
	}
}

func GetIntersectionPoint(line1, line2 Line) (Vector2, error) {
	equationLhs := [][]float64{
		{-line1.Direction.X, line2.Direction.X},
		{-line1.Direction.Y, line2.Direction.Y},
	}
	equationRhs := []float64{
		line1.Anchor.X - line2.Anchor.X,
		line1.Anchor.Y - line2.Anchor.Y,
	}

	res, err := SolveLinearEquations(equationLhs, equationRhs)
	if err == nil {
		m := res[0]
		return line1.Anchor.Add(line1.Direction.Multiply(m)), nil
	}

	return Vector2{}, err

}

func SolveLinearEquations(equationLhs [][]float64, equationRhs []float64) ([]float64, error) {
	panic("unimplemented")
}

func (v Vector2) Add(v2 Vector2) Vector2 {
	return Vector2{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
	}
}

func (v Vector2) Multiply(m float64) Vector2 {
	return Vector2{
		X: v.X * m,
		Y: v.Y * m,
	}
}

func (line Line) IsBetween(p1, p2 Vector2) bool {

	diffVec1 := line.Anchor.Add(p1.Multiply(-1))
	diffVec2 := line.Anchor.Add(p2.Multiply(-1))
	prod1 := CrossProduct(AsVector3(diffVec1), AsVector3(line.Direction))
	prod2 := CrossProduct(AsVector3(diffVec2), AsVector3(line.Direction))

	if prod1.Z*prod2.Z < 0 {
		return true
	} else {
		return false
	}

}

func AsVector3(vec2d Vector2) Vector3 {
	return Vector3{
		X: vec2d.X,
		Y: vec2d.Y,
		Z: 0,
	}
}

func GetUniquepoints(points []Vector2) []Vector2 {

	uniqueVectors := map[float64]map[float64]bool{}

	for _, vec := range points {
		if x, ok := uniqueVectors[vec.X]; !ok {
			if _, ok2 := x[vec.Y]; !ok2 {
				uniqueVectors[vec.X][vec.Y] = true
			}
		}
	}

	retVectors := []Vector2{}
	for x, yVals := range uniqueVectors {
		for y := range yVals {
			retVectors = append(retVectors, Vector2{X: x, Y: y})
		}
	}

	return retVectors
}
