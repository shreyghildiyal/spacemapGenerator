package mapGen

import (
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

	for _, star := range stars {
		wg.Add(1)
		// go populateBorders(star, grid, maxX, maxY, &wg)
		populateBorders(star, grid, maxX, maxY, &wg)
	}

	wg.Wait()

	return nil
}

func populateBorders(star Star, grid [][][]Star, maxX, maxY float64, wg *sync.WaitGroup) {
	boundaryCorners, borders, err := neighbourMap(star, grid, maxX, maxY)
	if err == nil {
		fmt.Println(star.Id, "Boundary corner count", len(boundaryCorners), star.X, star.Y)
		star.BoundaryCorners = boundaryCorners
		star.Borders = borders

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
					updateBoundary(star, cellStar, borders)
				}
			}
		}
	}

	boundaryPoints := []cartesian.Vector2{}

	for starId, borderLine := range borders {

		// fmt.Println("Anchor", borderLine.Anchor)
		boundaryPoints = append(boundaryPoints, borderLine.Anchor)
		// fmt.Println("Endpoint", borderLine.EndPoint())
		boundaryPoints = append(boundaryPoints, borderLine.EndPoint())

		if starId < 0 {
			delete(borders, starId)
		}
	}

	// fmt.Println(boundaryPoints)
	boundaryPoints = cartesian.GetUniquepoints(boundaryPoints)

	return boundaryPoints, borders, nil
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

func updateBoundary(star, gStar Star, borders map[int]cartesian.Line2D) {

	// get the bisecting line between star and gStar
	bisectingLine := cartesian.GetBisectingLine(star.Vector2, gStar.Vector2)

	intersectionPoints := []cartesian.Vector2{}
	for starId, borderLine := range borders {
		intersectionPoint, _, multiplier, err := cartesian.GetIntersectionPoint(bisectingLine, borderLine)
		if err == nil && multiplier >= 0 && multiplier <= 1 {
			if cartesian.IsSameSide(bisectingLine, borderLine.Anchor, star.Vector2) {
				newSegment := cartesian.GetLine(borderLine.Anchor, intersectionPoint)
				borders[starId] = newSegment
			} else {
				newSegment := cartesian.GetLine(intersectionPoint, borderLine.Anchor.Add(borderLine.Direction))
				borders[starId] = newSegment
			}
			// fmt.Println("Intersection Point", intersectionPoint)
			intersectionPoints = append(intersectionPoints, intersectionPoint)
		} else {
			anchorSameSide := cartesian.IsSameSide(bisectingLine, star.Vector2, borderLine.Anchor)
			endSameSide := cartesian.IsSameSide(bisectingLine, star.Vector2, borderLine.EndPoint())
			if !anchorSameSide && !endSameSide {
				delete(borders, starId)
			}
		}
	}

	if len(intersectionPoints) >= 2 {
		borders[gStar.Id] = cartesian.GetLine(intersectionPoints[0], intersectionPoints[1])
	}

}

// func setMoreCellsToBeChecked(cellX, cellY int, cellToCheck [][]int) int {
// 	panic("unimplemented")
// }

// func getCellToCheck(cellToCheck [][]int) (int, int) {

// 	panic("unimplemented")
// }

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
