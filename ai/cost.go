package ai

import (
	"math"

	"github.com/alexdor/dtu-ai-mas-final-assignment/types"
)

type Cost interface {
	Calculate(types.LevelInfo) int
}

type ManhattanDistance struct{}

func (ManhattanDistance) Calculate(levelInfo *types.LevelInfo) int {
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

func isBoxFree(coor types.Coordinates, levelInfo *types.LevelInfo) bool {
	if _, ok := levelInfo.WallsCoordinates[coor]; ok {
		return true
	}

	for _, coord := range levelInfo.AgentCoordinates {
		if _, ok := coord[coor]; ok {
			return true
		}
	}

	for _, coord := range levelInfo.BoxCoordinates {
		if _, ok := coord[coor]; ok {
			return true
		}
	}

	return false
}

func abs(x int) int {
	if x < 0 {
		return int(-x)
	}

	return int(x)
}
