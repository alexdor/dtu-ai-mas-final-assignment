package level

import (
	"math"
)

type Cost interface {
	Calculate(*CurrentState) int
}

type ManhattanDistance struct{}

func CalculateManhattanDistance(currentState *CurrentState) int {
	// TODO: preassign boxes to goals (by creating a rect)
	distance := 0

	for _, box := range currentState.Boxes {
		min := math.MaxInt64
		goals := currentState.LevelInfo.GoalCoordinates[box.Letter]

		for _, goalCoordinates := range goals {
			cost := abs(box.Coordinates[0]-goalCoordinates[0]) + abs(box.Coordinates[1]-goalCoordinates[1])
			cost += calculateWallsCost(box.Coordinates, goalCoordinates, currentState)

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

func calculateWallsCost(boxCoordinates Coordinates, goalCoordinates Coordinates, currentState *CurrentState) int {
	// TODO: Turn this to a percentage
	wallPenaltySize := 4

	isXcoordOfBoxSmallest := boxCoordinates[0] < goalCoordinates[0]
	isYcoordOfBoxSmallest := boxCoordinates[1] < goalCoordinates[1]

	smallXcoord := boxCoordinates[0]
	bigXcoord := goalCoordinates[0]

	if !isXcoordOfBoxSmallest {
		smallXcoord = goalCoordinates[0]
		bigXcoord = boxCoordinates[0]
	}

	smallYcoord := boxCoordinates[1]
	bigYcoord := goalCoordinates[1]

	if !isYcoordOfBoxSmallest {
		smallYcoord = goalCoordinates[1]
		bigYcoord = boxCoordinates[1]
	}

	cost := 0

	for x := 0; x < bigXcoord-smallXcoord; x++ {
		for y := 0; y < bigYcoord-smallYcoord; y++ {
			newCoor := Coordinates{smallXcoord + x, smallYcoord + y}
			if currentState.LevelInfo.IsWall(newCoor) {
				cost += wallPenaltySize
			}
		}
	}

	return cost
}
