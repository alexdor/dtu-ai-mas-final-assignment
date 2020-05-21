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
			cost += CalculateWallsCost(box.Coordinates, goalCoordinates, currentState)

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

func CalculateWallsCost(boxCoordinates Coordinates, goalCoordinates Coordinates, currentState *CurrentState) int {
	isXcoordOfBoxSmallest := boxCoordinates[0] < goalCoordinates[0]
	isYcoordOfBoxSmallest := boxCoordinates[1] < goalCoordinates[1]

	var smallXcoord, bigXcoord, smallYcoord, bigYcoord int

	if isXcoordOfBoxSmallest {
		smallXcoord = boxCoordinates[0]
		bigXcoord = goalCoordinates[0]
	} else {
		smallXcoord = goalCoordinates[0]
		bigXcoord = boxCoordinates[0]
	}

	if isYcoordOfBoxSmallest {
		smallYcoord = boxCoordinates[1]
		bigYcoord = goalCoordinates[1]
	} else {
		smallYcoord = goalCoordinates[1]
		bigYcoord = boxCoordinates[1]
	}

	cost := 0

	for x := 0; x < bigXcoord-smallXcoord; x++ {
		for y := 0; y < bigYcoord-smallYcoord; y++ {
			newCoor := Coordinates{smallXcoord + x, smallYcoord + y}
			if currentState.LevelInfo.IsWall(newCoor) {
				cost = +1
			}
		}
	}

	return cost
}
