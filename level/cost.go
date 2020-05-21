package level

import (
	"math"
)

type Cost interface {
	Calculate(*CurrentState) int
}

type ManhattanDistance struct{}

func CalculateManhattanDistance(currentState *CurrentState) int {
	// TODO: preassign boxes to goals (by creating a rect), add the amount of walls in that rect in the calculation.
	distance := 0

	for _, box := range currentState.Boxes {
		min := math.MaxInt64
		goals := currentState.LevelInfo.GoalCoordinates[box.Letter]

		for _, goalCoordinates := range goals {
			cost := abs(box.Coordinates[0]-goalCoordinates[0]) + abs(box.Coordinates[1]-goalCoordinates[1])
			if cost < min {
				min = cost
			}
		}

		distance += min
	}

	return distance
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}
