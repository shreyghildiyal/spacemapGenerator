package mapGen

import (
	"math"

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

func AddStarBoundaries(stars []Star, maxX, maxY float64) error {

	// we divide the entire map into grids. The number of cells in the gridneeds to depend on the number of stars we are dealing with
	// we will create a map for all the stars in a particular grid
	// for each star we will create the boundary using the stars in its cell and the stars in the 24 cells surrounding it.
	// if a neighbour cell doesnt contain any stars, we will expand to neighbours of the cell as many times as needed
	// for this approach to give any advantage at all, we need the grid to be 6x6 atleast. More if possible.
	// Maybe even clusterCount x clusterCount

	grid := getPopulatedGrid(stars, maxX, maxY)

	for _, row := range grid {
		for _, cell := range row {
			for _, star := range cell {
				boundaryCorners, borders, err := neighbourMap(star, grid, maxY, maxY)
				if err != nil {
					star.BoundaryCorners = boundaryCorners
					star.Borders = borders
					star.Neighbours = getStarsFromBorders(borders)
				}
			}
		}
	}

	return nil
}

func getStarsFromBorders(borders map[int]cartesian.Line) []Star {
	panic("unimplemented")
}

func getPopulatedGrid(stars []Star, maxX, maxY float64) [][][]Star {

	grid := [][][]Star{}

	gridSide := math.Floor(math.Sqrt(float64(len(stars))))

	cellWidth := maxX / gridSide
	cellHeight := maxY / gridSide

	for _, star := range stars {
		cellX := int(math.Floor(star.GetX() / cellWidth))
		cellY := int(math.Floor(star.GetY() / cellHeight))

		grid[cellX][cellY] = append(grid[cellX][cellY], star)
	}

	return grid
}

func getLocationsFromStars(stars []Star) []cartesian.Vector2 {
	locArr := make([]cartesian.Vector2, len(stars))
	for i, star := range stars {
		locArr[i] = star.Vector2
	}
	return locArr
}

func neighbourMap(star Star, grid [][][]Star, maxX, maxY float64) ([]cartesian.Vector2, map[int]cartesian.Line, error) {

	borders := getInitialBoundarySegments(maxX, maxY)

	starCellX := int(math.Floor(star.X / maxX))
	starCellY := int(math.Floor(star.Y / maxY))

	cellToCheck := make([][]int, len(grid)) // 0 = dont check; 1 = check; -1 = already checked
	for i := range cellToCheck {
		cellToCheck[i] = make([]int, len(grid))
	}
	cellToCheckCount := 0

	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {

			cellToCheck[starCellX+i][starCellY+j] = 1
			cellToCheckCount++

		}
	}

	for cellToCheckCount > 0 {
		cellX, cellY := getCellToCheck(cellToCheck)
		gridStars := grid[cellX][cellY]
		if len(gridStars) == 0 {
			addedCellCount := setMoreCellsToBeChecked(cellX, cellY, cellToCheck)
			cellToCheckCount = cellToCheckCount + addedCellCount - 1
		} else {
			for _, gStar := range gridStars {
				updateBoundary(star, gStar, borders)
			}

		}
	}

	return nil, nil, nil
}

func updateBoundary(star, gStar Star, borders map[int]cartesian.Line) {
	panic("unimplemented")
}

func setMoreCellsToBeChecked(cellX, cellY int, cellToCheck [][]int) int {
	panic("unimplemented")
}

func getCellToCheck(cellToCheck [][]int) (int, int) {

	panic("unimplemented")
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
