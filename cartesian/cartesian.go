package cartesian

import "errors"

type Vector2 struct {
	X float64
	Y float64
}

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

type Line2D struct {
	Anchor    Vector2
	Direction Vector2
}

func GetPerpendicularVector(vec Vector2) Vector2 {
	rotatedVec3 := CrossProduct(AsVector3(vec), Vector3{X: 0, Y: 0, Z: 1})

	return AsVector2(rotatedVec3)
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

func DotProduct(vec1, vec2 Vector2) float64 {
	return vec1.X*vec2.X + vec1.Y*vec2.Y
}

func GetBisectingLine(point1, point2 Vector2) Line2D {
	midPoint := Vector2{
		X: (point1.X + point2.X) / 2,
		Y: (point1.Y + point2.Y) / 2,
	}

	perpendicularVector := GetPerpendicularVector(Vector2{X: point1.X - point2.X, Y: point1.Y - point2.Y})
	return Line2D{
		Anchor:    midPoint,
		Direction: perpendicularVector,
	}
}

func GetIntersectionPoint(line1, line2 Line2D) (Vector2, float64, float64, error) {

	// find the

	equationLhs := [][]float64{
		{line1.Direction.X, -line2.Direction.X},
		{line1.Direction.Y, -line2.Direction.Y},
	}
	equationRhs := []float64{
		line2.Anchor.X - line1.Anchor.X,
		line2.Anchor.Y - line1.Anchor.Y,
	}

	res, err := SolveLinearEquations(equationLhs, equationRhs) // res > [k, l] k is the multiplier for line 1. l is multiplier for line 2

	if err == nil {

		k := res[0]
		// if k < 0 || k > 1 {
		// 	return Vector2{}, 0, 0, errors.New("the line segments dont intersect. the lines do")
		// }
		return line1.Anchor.Add(line1.Direction.Multiply(k)), res[0], res[1], nil
	}

	return Vector2{}, 0, 0, err

}

func SolveLinearEquations(equationLhs [][]float64, equationRhs []float64) ([]float64, error) {

	if len(equationLhs) == 0 {
		return nil, errors.New("nothing on the left hand side")
	}

	for i := 0; i < len(equationLhs); i++ {
		if len(equationLhs[i]) != len(equationLhs) {
			return nil, errors.New("LHS dimensions are not acceptable")
		}
	}

	if len(equationLhs) != len(equationRhs) {
		return nil, errors.New("LHS and RHS dont match")
	}

	for i := 0; i < len(equationLhs); i++ {
		if equationLhs[i][i] == 0 {

			err := swapWithValidRow(i, equationLhs, equationRhs)
			if err != nil {
				return nil, err
			}

		}

		// normalize the row
		normalizeRow(i, equationLhs, equationRhs)

		// zero all the other items in the column
		removeOtherColVals(equationLhs, i, equationRhs)
	}

	result := make([]float64, len(equationLhs))

	for i := 0; i < len(equationLhs); i++ {
		result[i] = equationRhs[i] / equationLhs[i][i]
	}

	return result, nil

}

func removeOtherColVals(equationLhs [][]float64, baseRow int, equationRhs []float64) {
	for otherRow := 0; otherRow < len(equationLhs); otherRow++ {
		if otherRow != baseRow {
			multiplier := equationLhs[otherRow][baseRow] / equationLhs[baseRow][baseRow]
			for col := 0; col < len(equationLhs[otherRow]); col++ {
				equationLhs[otherRow][col] = equationLhs[otherRow][col] - equationLhs[baseRow][col]*multiplier
			}

			equationRhs[otherRow] = equationRhs[otherRow] - equationRhs[baseRow]*multiplier
		}

	}
}

func normalizeRow(row int, equationLhs [][]float64, equationRhs []float64) {
	divideVal := equationLhs[row][row]
	for j := 0; j < len(equationLhs[row]); j++ {
		equationLhs[row][j] = equationLhs[row][j] / divideVal
	}
	equationRhs[row] = equationRhs[row] / divideVal
}

func swapWithValidRow(baseRow int, equationLhs [][]float64, equationRhs []float64) error {
	swapWith := -1
	for j := baseRow + 1; j < len(equationLhs); j++ {
		if equationLhs[j][baseRow] != 0 {
			swapWith = j
			break
		}
	}
	if swapWith >= 0 {
		equationLhs[baseRow], equationLhs[swapWith] = equationLhs[swapWith], equationLhs[baseRow]
		equationRhs[baseRow], equationRhs[swapWith] = equationRhs[swapWith], equationRhs[baseRow]
	} else {
		// if no such row found, we dont have a solution. throw error
		return errors.New("no solution possible")
	}
	return nil
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

func (line Line2D) IsBetween(p1, p2 Vector2) bool {

	diffVec1 := line.Anchor.Subtract(p1)
	diffVec2 := line.Anchor.Subtract(p2)
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

func AsVector2(vec3d Vector3) Vector2 {
	return Vector2{
		X: vec3d.X,
		Y: vec3d.Y,
	}
}

func GetUniquepoints(points []Vector2) []Vector2 {

	uniqueVectors := map[float64]map[float64]bool{}

	for _, vec := range points {
		if _, ok := uniqueVectors[vec.X]; !ok {
			uniqueVectors[vec.X] = map[float64]bool{}
		}
		if _, ok := uniqueVectors[vec.X][vec.Y]; !ok {
			uniqueVectors[vec.X][vec.Y] = true
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

type LineSide int

const (
	SAME_SIDE LineSide = iota
	OPPOSITE_SIDE
	POINT1_ON_LINE
	POINT2_ON_LINE
)

func IsSameSide(line Line2D, point1, point2 Vector2) LineSide {

	product1 := CrossProduct(AsVector3(point1.Subtract(line.Anchor)), AsVector3(line.Direction))
	product2 := CrossProduct(AsVector3(point2.Subtract(line.Anchor)), AsVector3(line.Direction))

	if product1.Z == 0 {
		return POINT1_ON_LINE
	} else if product2.Z == 0 {
		return POINT2_ON_LINE
	} else if product1.Z*product2.Z > 0 {
		return SAME_SIDE
	} else {
		return OPPOSITE_SIDE
	}

}

func GetLine(start, end Vector2) Line2D {
	direction := end.Subtract(start)
	return Line2D{
		Anchor:    start,
		Direction: direction,
	}
}

func (vec Vector2) Subtract(vec2 Vector2) Vector2 {
	return Vector2{
		X: vec.X - vec2.X,
		Y: vec.Y - vec2.Y,
	}
}

func (line Line2D) EndPoint() Vector2 {
	return line.Anchor.Add(line.Direction)
}

func (vec Vector2) Equals(vec2 Vector2) bool {
	return vec.X == vec2.X && vec.Y == vec2.Y
}

func (vec Vector2) EqualsFuzzy(vec2 Vector2, slack float64) (bool, error) {

	xValid, err := slackEquals(vec.X, vec2.X, slack)
	if err != nil {
		return false, err
	}
	yValid, err := slackEquals(vec.Y, vec2.Y, slack)
	if err != nil {
		return false, err
	}

	return xValid && yValid, nil

}

func slackEquals(v1, v2 float64, slack float64) (bool, error) {
	if slack < 0 {
		return false, errors.New("the slack can not be less than 0")
	} else {
		if v1 > v2 {
			return v1-v2 < slack, nil
		} else {
			return v2-v1 < slack, nil
		}
	}
}
