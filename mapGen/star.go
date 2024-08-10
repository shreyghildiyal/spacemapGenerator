package mapGen

import "github.com/shreyghildiyal/spacemapGenerator/cartesian"

type Star struct {
	cartesian.Vector2
	Id              int
	ClusterId       int
	IsClusterCore   bool
	BoundaryCorners []cartesian.Vector2
	Neighbours      []Star
	Borders         map[int]cartesian.Line2D
}

func (star *Star) GetX() float64 {
	return star.X
}

func (star *Star) GetY() float64 {
	return star.Y
}
