package mapGen

import "fmt"

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
		stars[bestFromStarIndex].Neighbours = append(stars[bestFromStarIndex].Neighbours, bestStarIndex)
		stars[bestStarIndex].Neighbours = append(stars[bestStarIndex].Neighbours, bestFromStarIndex)
	}
	fmt.Println("Neighbours have been added. There should now be a path from every star to every other star")
}

func ConnectClusterStars(stars []Star) {

}
