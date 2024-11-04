package mapGen

import (
	"math"
	"math/rand"
)

type SystemGenConfig struct {
	SystemMinHeavenlyBodyCount int
	SystemMaxHeavenlyBodyCount int
	PlanetHabitabilityChance   float64
	MaxHeavenlyBodyType        int
	MinHeavenlyBodyDistance    int
	MaxStarType                int
	MaxSystemRadius            int
	MinSystemRadius            int
}

func GenerateStarSystems(stars []Star, config SystemGenConfig) []HeavenlyBody {
	// we want each system to have a random number of 'planets'
	// we want atleast one habitable planet in each cluster

	// do we want the planet generation to be more even?

	heavenlyBodies := []HeavenlyBody{}

	id := 0

	for starId := range stars {
		stars[starId].Type = rand.Intn(config.MaxStarType)
		stars[starId].SystemRadius = rand.Intn(config.MaxSystemRadius-config.MinSystemRadius) + config.MinSystemRadius
		desiredHeavenlyBodyCount := config.SystemMinHeavenlyBodyCount + rand.Intn(config.SystemMaxHeavenlyBodyCount-config.SystemMinHeavenlyBodyCount)
		distances := GetHeavenlyBodyDistances(desiredHeavenlyBodyCount, config.MinHeavenlyBodyDistance, config.MaxSystemRadius)
		for i := 0; i < desiredHeavenlyBodyCount; i++ {
			heavelyBody := HeavenlyBody{}
			if rand.Float64() <= config.PlanetHabitabilityChance {
				heavelyBody.IsHabitable = true
			} else {
				heavelyBody.IsHabitable = false
			}
			heavelyBody.Angle = rand.Float64() * math.Pi
			heavelyBody.StarDistance = distances[i]
			heavelyBody.Type = rand.Intn(config.MaxHeavenlyBodyType)
			// HeavenlyBody.StarDistance = 10
			heavenlyBodies = append(heavenlyBodies, heavelyBody)
			heavelyBody.Id = id
			stars[starId].HeavenlyBodies = append(stars[starId].HeavenlyBodies, heavelyBody.Id)
		}
	}

	return heavenlyBodies
}

func GetHeavenlyBodyDistances(desiredHeavenlyBodyCount, minDist, maxDist int) []int {
	interBodyMinDistance := (maxDist - minDist) / (2 * desiredHeavenlyBodyCount)

	generatedDistances := []int{}

	for len(generatedDistances) < desiredHeavenlyBodyCount {
		dist := minDist + rand.Intn(maxDist-minDist)
		if len(generatedDistances) == 0 {
			generatedDistances = append(generatedDistances, dist)
		} else {
			isTooClose := false
			for _, prevDist := range generatedDistances {
				nDist := dist - prevDist
				if nDist < 0 {
					nDist = nDist * -1
				}
				if nDist < interBodyMinDistance {
					isTooClose = true
					break
				}
			}
			if !isTooClose {
				generatedDistances = append(generatedDistances, dist)
			}
		}
	}
	return generatedDistances
}
