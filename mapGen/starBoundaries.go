package mapGen

import (
	"github.com/shreyghildiyal/spacemapGenerator/cartesian"
)

func AddStarBoundaries(stars []Star, maxX, maxY float64) error {

	// locations := getLocationsFromStars(stars)

	// for each star get its boundary
	for i := 0; i < len(stars); i++ {
		neighbourMap, err := getBoundary(stars[i], stars, maxX, maxY)
		if err != nil {
			return err
		}
		stars[i].NeighbourMap = neighbourMap
	}

	return nil
}

func getLocationsFromStars(stars []Star) []cartesian.Vector2 {
	locArr := make([]cartesian.Vector2, len(stars))
	for i, star := range stars {
		locArr[i] = star.Vector2
	}
	return locArr
}

func getBoundary(loc Star, otherLocs []Star, maxX, maxY float64) (map[int]cartesian.Line, error) {

	boundary := getInitialBoundarySegments(maxX, maxY)

	for otherLocId, otherLoc := range otherLocs {
		bisectingLine := cartesian.GetBisectingLine(loc.Vector2, otherLoc.Vector2)
		intersectionPoints := []cartesian.Vector2{}
		intersectingBoundaryLines := map[int]bool{}
		for boundaryLineId, boundaryLine := range boundary {
			intersectionPoint, err := cartesian.GetIntersectionPoint(bisectingLine, boundaryLine)
			if err != nil {
				var m float64
				if boundaryLine.Direction.X != 0 {
					m = (intersectionPoint.X - boundaryLine.Anchor.X) / boundaryLine.Direction.X
				} else {
					m = (intersectionPoint.Y - boundaryLine.Anchor.Y) / boundaryLine.Direction.Y
				}

				if m > 0 && m < 1 {
					// this is a valid intersection and the intersection point is a potential part of the boundary
					intersectionPoints = append(intersectionPoints, intersectionPoint)
					intersectingBoundaryLines[boundaryLineId] = true
				}
			}
		}

		if len(intersectionPoints) != 0 {
			newBoundary := []cartesian.Line{}
			for boundaryLineId, boundaryLine := range boundary {
				if _, ok := intersectingBoundaryLines[boundaryLineId]; ok {
					// break the boundaryLine  and add only relevant part to the new boundary List
				} else {
					// check if both points are on the valid side or not. If both are on valid side, add to new boundary else delete the boundary.
				}
			}
		}
	}
	return nil, nil
}

func getInitialBoundarySegments(maxX, maxY float64) map[int]cartesian.Line {

	return map[int]cartesian.Line{
		-1: {
			Anchor: cartesian.Vector2{
				X: 0,
				Y: 0,
			},
			Direction: cartesian.Vector2{
				X: maxX,
				Y: 0,
			},
		},
		-2: {
			Anchor: cartesian.Vector2{
				X: maxX,
				Y: 0,
			},
			Direction: cartesian.Vector2{
				X: 0,
				Y: maxY,
			},
		},
		-3: {
			Anchor: cartesian.Vector2{
				X: maxX,
				Y: maxY,
			},
			Direction: cartesian.Vector2{
				X: -maxX,
				Y: 0,
			},
		},
		-4: {
			Anchor: cartesian.Vector2{
				X: 0,
				Y: maxY,
			},
			Direction: cartesian.Vector2{
				X: 0,
				Y: -maxY,
			},
		},
	}
}
