package mapGen

import (
	"errors"

	"github.com/shreyghildiyal/spacemapGenerator/cartesian"
)

func AddStarBoundaries(stars []Star, maxX, maxY float64) error {

	locations := getLocationsFromStars(stars)

	// for each star get its boundary
	for i := 0; i < len(stars); i++ {
		neighbourMap, err := getBoundary(stars[i].Vector2, append(locations[:i], locations[i+1:]...), maxX, maxY)
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

func getBoundary(loc cartesian.Vector2, allOtherLocs []cartesian.Vector2, maxX, maxY float64) (map[int]cartesian.Line, error) {

	boundary := getInitialBoundarySegments(maxX, maxY)

	for _, otherLoc := range allOtherLocs {
		bisectingLine := cartesian.GetBisectingLine(cartesian.Vector2{X: loc.X, Y: loc.Y}, cartesian.Vector2{X: otherLoc.X, Y: otherLoc.Y})
		intersectionPoints := []cartesian.Vector2{}
		for _, boundaryLine := range boundary {
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
				}
			}
		}

		intersectionPoints = cartesian.GetUniquepoints(intersectionPoints)

		if len(intersectionPoints) == 2 {
			// this is the happy case.
		} else if len(intersectionPoints) == 0 {
			// kinda happy case. Nothing to be done
		} else {
			// something has gone horribly wrong?
			return nil, errors.New("the number of intersection points was unexpected")
		}
	}
	return nil, nil
}

func getInitialBoundarySegments(maxX, maxY float64) []cartesian.Line {

	return []cartesian.Line{
		{
			Anchor: cartesian.Vector2{
				X: 0,
				Y: 0,
			},
			Direction: cartesian.Vector2{
				X: maxX,
				Y: 0,
			},
		},
		{
			Anchor: cartesian.Vector2{
				X: maxX,
				Y: 0,
			},
			Direction: cartesian.Vector2{
				X: 0,
				Y: maxY,
			},
		},
		{
			Anchor: cartesian.Vector2{
				X: maxX,
				Y: maxY,
			},
			Direction: cartesian.Vector2{
				X: -maxX,
				Y: 0,
			},
		},
		{
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
