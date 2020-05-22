package level

type Cost interface {
	Calculate(*CurrentState) int
}

type ManhattanDistance struct{}

func CalculateManhattanDistance(currentState *CurrentState) int {
	aggregatedCost := 0

	for i, box := range currentState.Boxes {
		goals := currentState.LevelInfo.GoalCoordinates[box.Letter]
		goalIndex := -1

		for j, goalCoordinates := range goals {
			if goalCoordinates == currentState.LevelInfo.BoxGoalAssignment[i] {
				goalIndex = j
				break

		if goalIndex == -1 {
			continue
		}
		goal := goals[goalIndex]

		aggregatedCost += ManhattanPlusPlus(box.Coordinates, goal, currentState)
	}

	return aggregatedCost
}

func ManhattanPlusPlus(first, second Coordinates, state *CurrentState) int {
	diff := abs(first[0]-second[0]) + abs(first[1]-second[1])
	if diff == 0 {
		return 0
	}

	diff += calculateWallsCost(first, second, state)

	return diff
}

func calculateWallsCost(boxCoordinates Coordinates, goalCoordinates Coordinates, currentState *CurrentState) int {
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

	isAreaCheckable := smallXcoord < bigXcoord && smallYcoord < bigYcoord
	if !isAreaCheckable {
		return 0
	}

	for _, wallCoordinate := range currentState.LevelInfo.InGameWallsCoordinates {
		isWallXcoordWithinRectangle := wallCoordinate[0] > smallXcoord && wallCoordinate[0] < bigXcoord
		isWallYcoordWithinRectangle := wallCoordinate[1] > smallYcoord && wallCoordinate[1] < bigYcoord
		isWallCoordWithinRectangle := isWallXcoordWithinRectangle && isWallYcoordWithinRectangle

		if !isWallCoordWithinRectangle {
			continue
		}

		cost += wallPenaltySize
	}

	// account for size of area checked
	cost /= (bigXcoord - smallXcoord) * (bigYcoord - smallYcoord)

	return cost
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}
