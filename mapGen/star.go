package mapGen

import "github.com/shreyghildiyal/spacemapGenerator/cartesian"

type Star struct {
	cartesian.Vector2
	Id              int
	ClusterId       int
	IsClusterCore   bool
	BoundaryCorners []cartesian.Vector2      // used for display purposes. Defines a convext polygon. Polygon contains all points such that the star is closest to them
	Neighbours      []int                    // the stars that have a warp line connecting to this star. have to be among the keys of Border
	Borders         map[int]cartesian.Line2D // useful when knowing that the domain is constrained by a star and the corresponding line. There might not be any connections between the 2
	HeavenlyBodies  []int                    // ids of the heavenly bodies in the star system
	Type            int
	SystemRadius    int
}

func (star *Star) GetX() float64 {
	return star.X
}

func (star *Star) GetY() float64 {
	return star.Y
}
