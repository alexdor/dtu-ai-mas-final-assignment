package level

type Cost interface {
	Calculate(*CurrentState) int
}

// This is the cost that should be added to the total cost
// when there is a full row of walls
const additionalWallCostForFullRow = 2

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

	boxColor := state.LevelInfo.BoxColor[box.Letter]

	for _, agent := range state.Agents {
		isAgentAndBoxTogetherLikeBros := state.LevelInfo.AgentColor[agent.Letter] == boxColor
		if !isAgentAndBoxTogetherLikeBros {
			continue
		}

		diff += manhattenDistance(box.Coordinates, agent.Coordinates)
		diff += calculateWallsCost(agent.Coordinates, box.Coordinates, state)
	}

	diff += calculateWallsCost(first, second, state)

	return diff
}

func manhattenDistance(first, second Coordinates) int {
	return abs(first[0]-second[0]) + abs(first[1]-second[1])
}

func calculateWallsCost(firstCoordinates, secondCoordinates Coordinates, currentState *CurrentState) int {
	isXcoordOfBoxSmallest := firstCoordinates[0] < secondCoordinates[0]
	isYcoordOfBoxSmallest := firstCoordinates[1] < secondCoordinates[1]

	smallXcoord := firstCoordinates[0]
	bigXcoord := secondCoordinates[0]

	if !isXcoordOfBoxSmallest {
		smallXcoord, bigXcoord = bigXcoord, smallXcoord
	}

	smallYcoord := firstCoordinates[1]
	bigYcoord := secondCoordinates[1]

	if !isYcoordOfBoxSmallest {
		smallYcoord, bigYcoord = bigYcoord, smallYcoord
	}

	isAreaCheckable := smallXcoord < bigXcoord && smallYcoord < bigYcoord

	if !isAreaCheckable {
		return 0
	}

	for x := smallXcoord; x <= bigXcoord; x++ {
		wallColumns, ok := currentState.LevelInfo.WallRows[x]
		if !ok {
			continue
		}
		for _, wallY := range wallColumns {
			// If the smallest column of the wall is bigger than
			// the biggest column of the target then break
			if wallY[0] > bigYcoord {
				break
			}

			// Wall expands the full length
			if wallY[0] <= smallYcoord && wallY[1] >= bigYcoord {
				return min(smallYcoord-wallY[0], wallY[1]-bigYcoord) + additionalWallCostForFullRow
			}

		}
	}

	return 0
}

func min(x, y int) int {
	if x < y {
		return x
	}

	return y
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

		cost += calculateWallsCost(box.Coordinates, agent.Coordinates, state)
	}

	return cost
}
