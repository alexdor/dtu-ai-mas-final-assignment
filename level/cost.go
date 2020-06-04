package level

type Cost interface {
	Calculate(*CurrentState) int
}

type ManhattanDistance struct{}

func CalculateAggregatedCost(currentState *CurrentState) int {
	aggregatedCost := 0

	for i, box := range currentState.Boxes {
		goals, ok := currentState.LevelInfo.GoalCoordinates[box.Letter]
		if !ok {
			continue
		}

		goalIndex := -1

		for j, goalCoordinates := range goals {
			if goalCoordinates == currentState.LevelInfo.BoxGoalAssignment[i] {
				goalIndex = j
				break
			}
		}

		if goalIndex == -1 {
			continue
		}
		goal := goals[goalIndex]

		manHattanCost := ManhattanPlusPlus(box.Coordinates, goal, currentState, &box)
		aggregatedCost += manHattanCost

		isManhattanCostZero := manHattanCost == 0
		if isManhattanCostZero {
			continue
		}

		aggregatedCost += calculateAgentsToBoxCost(currentState, &box)
	}

	return aggregatedCost
}

func ManhattanPlusPlus(first, second Coordinates, state *CurrentState, box *NodeOrAgent) int {
	diff := manhattenDistance(first, second)
	if diff == 0 {
		return 0
	}

	diff += calculateWallsCost(first, second, state)

	return diff
}

func manhattenDistance(first, second Coordinates) int {
	return abs(first[0]-second[0]) + abs(first[1]-second[1])
}

func calculateWallsCost(boxCoordinates Coordinates, goalCoordinates Coordinates, currentState *CurrentState) int {
	wallPenaltySize := 4

	isXcoordOfBoxSmallest := boxCoordinates[0] < goalCoordinates[0]
	isYcoordOfBoxSmallest := boxCoordinates[1] < goalCoordinates[1]

	smallXcoord := boxCoordinates[0]
	bigXcoord := goalCoordinates[0]

	if !isXcoordOfBoxSmallest {
		smallXcoord, bigXcoord = bigXcoord, smallXcoord
	}

	smallYcoord := boxCoordinates[1]
	bigYcoord := goalCoordinates[1]

	if !isYcoordOfBoxSmallest {
		smallYcoord, bigYcoord = bigYcoord, smallYcoord
	}

	isAreaCheckable := smallXcoord < bigXcoord && smallYcoord < bigYcoord

	if !isAreaCheckable {
		return 0
	}

	cost := 0

	var isWallXcoordWithinRectangle, isWallYcoordWithinRectangle, isWallCoordWithinRectangle bool

	for _, wallCoordinate := range currentState.LevelInfo.InGameWallsCoordinates {

		isWallXcoordWithinRectangle = wallCoordinate[0] > smallXcoord && wallCoordinate[0] < bigXcoord

		isWallYcoordWithinRectangle = wallCoordinate[1] > smallYcoord && wallCoordinate[1] < bigYcoord

		isWallCoordWithinRectangle = isWallXcoordWithinRectangle && isWallYcoordWithinRectangle

		if !isWallCoordWithinRectangle {
			if wallCoordinate[0] > bigXcoord {
				break
			}
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

func calculateAgentsToBoxCost(state *CurrentState, box *NodeOrAgent) int {
	cost := 0
	boxColor := state.LevelInfo.BoxColor[box.Letter]

	for _, agent := range state.Agents {
		isAgentAndBoxTogetherLikeBros := state.LevelInfo.AgentColor[agent.Letter] == boxColor
		if !isAgentAndBoxTogetherLikeBros {
			continue
		}

		cost += manhattenDistance(box.Coordinates, agent.Coordinates)
	}

	return cost
}
