package mapGen

import (
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/shreyghildiyal/spacemapGenerator/cartesian"
)

// We are going to add some points around the star as the boundary to test our display
func AddDummyStarBoundaries(stars []Star, maxX, maxY float64) error {

	println("The dummy corners are going to be generated")
	const CORNERCOUNT = 6
	const DISTANCE = 20

	for i := 0; i < len(stars); i++ {
		// println(stars[i].Id)
		for j := 0; j < CORNERCOUNT; j++ {
			angle := float64(j) * 2 * math.Pi / CORNERCOUNT
			dx := DISTANCE * math.Cos(angle)
			dy := DISTANCE * math.Sin(angle)
			// println(dx, math.Cos(angle), dy, math.Sin(angle), float64(j), angle)
			point := cartesian.Vector2{
				X: stars[i].X + dx,
				Y: stars[i].Y + dy,
			}
			stars[i].BoundaryCorners = append(stars[i].BoundaryCorners, point)
		}
		// println("boundary corner count", len(stars[i].BoundaryCorners))
	}
	// println("pi/4", math.Pi/4)
	// println("3pi/4", 3*math.Pi/4)

	// println("sin(pi/4)", math.Sin(math.Pi/4))
	// println("sin(pi/2)", math.Sin(math.Pi/2))
	return nil
}

func AddDummyNeighbours(stars []Star) {
	addedStars := []int{}     // index in stars of the star
	remainingStars := []int{} // index in stars of the star

	distanceGrid := make([][]float64, len(stars))
	for i := range stars {
		distanceGrid[i] = make([]float64, len(stars))
	}

	for i := 0; i < len(stars); i++ {
		for j := 0; j < len(stars); j++ {
			distanceGrid[i][j] = getDist(stars[i], stars[j]) // distance using the index within stars slice
		}
	}

	addedStars = append(addedStars, 0) // the first star in stars is where we start
	for i := 1; i < len(stars); i++ {
		remainingStars = append(remainingStars, i) // all other indexes are 'remaining'
	}

	for len(remainingStars) > 0 {
		bestRemIndex := 0
		bestStarIndex := remainingStars[bestRemIndex]
		bestDistance, bestFromStarIndex := getMinDist(bestStarIndex, addedStars, distanceGrid)

		for remIndex := 1; remIndex < len(remainingStars); remIndex++ {
			starIndex := remainingStars[remIndex]
			dist, fromStarIndex := getMinDist(starIndex, addedStars, distanceGrid)
			if dist < bestDistance {
				bestDistance = dist
				bestRemIndex = remIndex
				bestStarIndex = starIndex
				bestFromStarIndex = fromStarIndex
			}
		}

		remainingStars = append(remainingStars[:bestRemIndex], remainingStars[bestRemIndex+1:]...)
		addedStars = append(addedStars, bestStarIndex)
		stars[bestFromStarIndex].Neighbours = append(stars[bestFromStarIndex].Neighbours, stars[bestStarIndex])
		stars[bestStarIndex].Neighbours = append(stars[bestStarIndex].Neighbours, stars[bestFromStarIndex])
	}
	fmt.Println("Neighbours have been added. There should now be a path from every star to every other star")
}

func getMinDist(bestStarIndex int, addedStars []int, distanceGrid [][]float64) (float64, int) {

	bestFromStar := 0
	minDist := distanceGrid[bestStarIndex][addedStars[bestFromStar]]
	for i := 1; i < len(addedStars); i++ {
		dist := distanceGrid[bestStarIndex][addedStars[i]]
		if dist < minDist {
			minDist = dist
			bestFromStar = i
		}
	}
	return minDist, addedStars[bestFromStar]

}

func AddStarBoundaries(stars []Star, maxX, maxY float64) error {

	// we divide the entire map into grids. The number of cells in the gridneeds to depend on the number of stars we are dealing with
	// we will create a map for all the stars in a particular grid
	// for each star we will create the boundary using the stars in its cell and the stars in the 24 cells surrounding it.
	// if a neighbour cell doesnt contain any stars, we will expand to neighbours of the cell as many times as needed
	// for this approach to give any advantage at all, we need the grid to be 6x6 atleast. More if possible.
	// Maybe even clusterCount x clusterCount

	grid := getPopulatedGrid(stars, maxX, maxY)

	wg := sync.WaitGroup{}

	for i, _ := range stars {
		wg.Add(1)
		// go populateBorders(star, grid, maxX, maxY, &wg)
		populateBorders(i, stars, grid, maxX, maxY, &wg)
	}

	wg.Wait()

	return nil
}

func populateBorders(starId int, stars []Star, grid [][][]Star, maxX, maxY float64, wg *sync.WaitGroup) {
	boundaryCorners, borders, err := neighbourMap(stars[starId], grid, maxX, maxY)

	if err == nil {
		fmt.Println(starId, "Boundary corner count", len(boundaryCorners), stars[starId].X, stars[starId].Y)
		stars[starId].BoundaryCorners = boundaryCorners
		stars[starId].Borders = borders

	} else {
		fmt.Println(err.Error())
	}
	wg.Done()
}

// func getStarsFromBorders(borders map[int]cartesian.Line, stars []Star) []Star {
// 	selectedStars := []Star{}
// 	for starId, _ := range borders {
// 		selectedStars = append(selectedStars, stars[starId])
// 	}
// 	return selectedStars
// }

func getPopulatedGrid(stars []Star, maxX, maxY float64) [][][]Star {

	gridSide := math.Floor(math.Sqrt(float64(len(stars))))

	grid := make([][][]Star, int(gridSide))
	for i := 0; i < int(gridSide); i++ {
		grid[i] = make([][]Star, int(gridSide))
	}

	cellWidth := maxX / gridSide
	cellHeight := maxY / gridSide

	for _, star := range stars {
		cellX := int(math.Floor(star.GetX() / cellWidth))
		cellY := int(math.Floor(star.GetY() / cellHeight))

		grid[cellX][cellY] = append(grid[cellX][cellY], star)
	}

	return grid
}

func neighbourMap(star Star, grid [][][]Star, maxX, maxY float64) ([]cartesian.Vector2, map[int]cartesian.Line2D, error) {

	borders := getInitialBoundarySegments(maxX, maxY)

	cellToCheck := getCellsToCheck(star, maxX, maxY, grid) // 0 = dont check; 1 = check; -1 = already checked

	for row := 0; row < len(cellToCheck); row++ {
		for col := 0; col < len(cellToCheck); col++ {
			if cellToCheck[row][col] == 1 {
				cellStars := grid[row][col]
				for _, cellStar := range cellStars {
					if star.Id != cellStar.Id {
						updateBoundary(star, cellStar, borders)
					}

				}
			}
		}
	}

	// fmt.Println("Anchor", borderLine.Anchor)
	// fmt.Println("Endpoint", borderLine.EndPoint())
	// fmt.Println(boundaryPoints)
	boundaryPoints, borders, err := cleanupBorders(borders)
	if err == nil {
		// fmt.Println("boundary star ", star.Id, boundaryPoints)
		return boundaryPoints, borders, nil
	} else {
		return []cartesian.Vector2{}, map[int]cartesian.Line2D{}, err
	}

}

func cleanupBorders(borders map[int]cartesian.Line2D) ([]cartesian.Vector2, map[int]cartesian.Line2D, error) {
	boundaryLines := []cartesian.Line2D{}

	for starId, borderLine := range borders {

		boundaryLines = append(boundaryLines, borderLine)
		// boundaryPoints = append(boundaryPoints, borderLine.Anchor)

		// boundaryPoints = append(boundaryPoints, borderLine.EndPoint())

		if starId < 0 {
			delete(borders, starId)
		}
	}

	boundaryPoints, err := getOrderedBoundaryPoints(boundaryLines)

	// boundaryPoints = cartesian.GetUniquepoints(boundaryPoints)

	return boundaryPoints, borders, err
}

func getOrderedBoundaryPoints(boundaryLines []cartesian.Line2D) ([]cartesian.Vector2, error) {

	orderedPoints := []cartesian.Vector2{}

	orderedPoints = append(orderedPoints, boundaryLines[0].Anchor)
	nextPoint := boundaryLines[0].EndPoint()
	addedIndexes := make([]bool, len(boundaryLines))
	addedIndexes[0] = true
	for !nextPoint.Equals(boundaryLines[0].Anchor) {
		fmt.Println(nextPoint)
		nextIndex := -1
		reverse := false

		for i := 0; i < len(boundaryLines); i++ {
			if !addedIndexes[i] {
				line := boundaryLines[i]
				if nextPoint.Equals(line.Anchor) {
					nextIndex = i
					reverse = false
					break
				} else if nextPoint.Equals(line.EndPoint()) {
					nextIndex = i
					reverse = true
					break
				}
			}
		}
		if nextIndex == -1 {
			return nil, errors.New("issue in ordering the boundary points")
		} else {
			if reverse {
				orderedPoints = append(orderedPoints, boundaryLines[nextIndex].EndPoint())
				nextPoint = boundaryLines[nextIndex].Anchor

			} else {
				orderedPoints = append(orderedPoints, boundaryLines[nextIndex].Anchor)
				nextPoint = boundaryLines[nextIndex].EndPoint()
			}
			addedIndexes[nextIndex] = true
		}
	}

	return orderedPoints, nil
}

func getCellsToCheck(star Star, maxX float64, maxY float64, grid [][][]Star) [][]int {
	starCellX := int(math.Floor(star.X / maxX))
	starCellY := int(math.Floor(star.Y / maxY))

	cellToCheck := make([][]int, len(grid))
	for i := range cellToCheck {
		cellToCheck[i] = make([]int, len(grid))
	}

	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {
			addToCheck(starCellX+i, starCellY+j, grid, cellToCheck)
		}
	}
	return cellToCheck
}

func addToCheck(row, col int, grid [][][]Star, cellToCheck [][]int) {

	if row < 0 || col < 0 || row >= len(grid) || col >= len(grid[0]) {
		return
	}

	if cellToCheck[row][col] != 0 {
		return
	}

	if len(grid[row][col]) > 0 {
		cellToCheck[row][col] = 1
	} else {
		cellToCheck[row][col] = -1
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {

				addToCheck(row+i, col+j, grid, cellToCheck)

			}
		}
	}

}

func updateBoundary(star, gStar Star, borders map[int]cartesian.Line2D) error {

	// get the bisecting line between star and gStar
	bisectingLine := cartesian.GetBisectingLine(star.Vector2, gStar.Vector2)

	intersectionPoints := []cartesian.Vector2{}
	for starId, borderLine := range borders {

		intersectionPoint, _, multiplier, err := cartesian.GetIntersectionPoint(bisectingLine, borderLine)

		if err == nil && multiplier >= 0 && multiplier <= 1 {

			err := handleIntersection(bisectingLine, borderLine, star, intersectionPoint, borders, starId)
			if err != nil {
				return fmt.Errorf("issue in handling intersection with borderline %w", err)
			}
			intersectionPoints = append(intersectionPoints, intersectionPoint)
		} else {
			handleNonIntersection(bisectingLine, star, borderLine, borders, starId)
		}
	}

	if len(intersectionPoints) >= 2 {
		borders[gStar.Id] = cartesian.GetLine(intersectionPoints[0], intersectionPoints[1])
	}
	return nil
}

func handleNonIntersection(bisectingLine cartesian.Line2D, star Star, borderLine cartesian.Line2D, borders map[int]cartesian.Line2D, starId int) {
	anchorSameSide := cartesian.IsSameSide(bisectingLine, star.Vector2, borderLine.Anchor)
	endSameSide := cartesian.IsSameSide(bisectingLine, star.Vector2, borderLine.EndPoint())
	if anchorSameSide == cartesian.OPPOSITE_SIDE && endSameSide == cartesian.OPPOSITE_SIDE {
		delete(borders, starId)
	}
}

func handleIntersection(bisectingLine cartesian.Line2D, borderLine cartesian.Line2D, star Star, intersectionPoint cartesian.Vector2, borders map[int]cartesian.Line2D, starId int) error {
	anchorSameSide := cartesian.IsSameSide(bisectingLine, borderLine.Anchor, star.Vector2)
	endPointSameSide := cartesian.IsSameSide(bisectingLine, borderLine.EndPoint(), star.Vector2)

	if anchorSameSide == cartesian.POINT1_ON_LINE && endPointSameSide == cartesian.POINT1_ON_LINE {
		// error condition. this should not be possible
	} else if anchorSameSide == cartesian.POINT1_ON_LINE {
		if endPointSameSide != cartesian.SAME_SIDE {
			// the line needs to be deleted
			delete(borders, starId)
		}
	} else if endPointSameSide == cartesian.POINT1_ON_LINE {
		if anchorSameSide != cartesian.SAME_SIDE {
			// the line needs to be deleted
			delete(borders, starId)
		}
	} else if anchorSameSide == cartesian.SAME_SIDE {
		// crop the line from anchor to intersection point

		newSegment := cartesian.GetLine(borderLine.Anchor, intersectionPoint)
		borders[starId] = newSegment

	} else if endPointSameSide == cartesian.SAME_SIDE {
		newSegment := cartesian.GetLine(intersectionPoint, borderLine.EndPoint())
		borders[starId] = newSegment

	} else {
		return errors.New("Something unexpected is happening")
	}
	return nil
}

func getInitialBoundarySegments(maxX, maxY float64) map[int]cartesian.Line2D {

	return map[int]cartesian.Line2D{
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
