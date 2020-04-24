package ai

import (
	"math"

	"github.com/alexdor/dtu-ai-mas-final-assignment/level"
)

type Cost interface {
	Calculate(*level.Info) int
}

type ManhattanDistance struct{}

func (ManhattanDistance) Calculate(levelInfo *level.Info) int {
	// Should we add the distance from the agents to the box here?
	distance := 0

	for char, coordinateMap := range levelInfo.BoxCoordinates {
		for boxCoord := range coordinateMap {
			min := math.MaxInt64

			for goalCoord := range levelInfo.GoalCoordinates[char] {
				cost := abs(int(boxCoord[0]-goalCoord[0])) + abs(int(boxCoord[1]-goalCoord[1]))
				if cost < min {
					min = cost
				}
			}

			distance += min
		}
	}

	return distance
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}
