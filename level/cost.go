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
		box := box
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

		manHattanCost := ManhattanPlusPlus(box.Coordinates, goal, currentState, &box, i)

		isManhattanCostZero := manHattanCost == 0
		if isManhattanCostZero {
			continue
		}
		aggregatedCost += manHattanCost

		aggregatedCost += calculateAgentsToBoxCost(currentState, &box, i)
	}

	return aggregatedCost
}

func ManhattanPlusPlus(first, second Coordinates, state *CurrentState, box *NodeOrAgent, boxIndex int) int {
	diff := manhattenDistance(first, second)
	if diff == 0 {
		return 0
	}
	diff += calculateWallsCost(first, second, state.LevelInfo)

	return diff
}

func manhattenDistance(first, second Coordinates) int {
	return abs(first[0]-second[0]) + abs(first[1]-second[1])
}

func calculateWallsCost(firstCoordinates, secondCoordinates Coordinates, levelInfo *Info) int {
	if levelInfo.ZeroInGameWalls {
		return 0
	}
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

	return rowOrColWallCalc(smallXcoord, bigXcoord, smallYcoord, bigYcoord, levelInfo.WallRows)
}

func rowOrColWallCalc(smallXcoord, bigXcoord, smallYcoord, bigYcoord int, walls ContinuosWalls) int {
	cost, tmpCost := 0, 0
	for x := smallXcoord; x <= bigXcoord; x++ {
		wallColumns, ok := walls[x]
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

			tmpCost = min(abs(smallYcoord-wallY[1]), abs(wallY[0]-bigYcoord))

			if tmpCost > cost {
				cost = tmpCost
			}

		}
	}
	return cost
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

func calculateAgentsToBoxCost(state *CurrentState, box *NodeOrAgent, boxIndex int) int {
	cost := 0
	agent := state.Agents[0]
	isAgentAndBoxTogetherLikeBros := true

	if state.LevelInfo.IsSingleAgent {
		isAgentAndBoxTogetherLikeBros = state.LevelInfo.AgentColor[agent.Letter] == state.LevelInfo.BoxColor[box.Letter]
	} else {
		agentIndex, ok := state.LevelInfo.BoxIndexToAgentIndex[boxIndex]

		// Fallback in case agentIndex isn't found
		if !ok || agentIndex < 0 {
		outer:
			for agentLetter, boxIndexes := range state.LevelInfo.AgentBoxAssignment {
				for _, indexOfBox := range boxIndexes {
					if indexOfBox == boxIndex {
						for _, stateAgent := range state.Agents {
							if agentLetter == stateAgent.Letter {
								agent = stateAgent
								break outer
							}
						}
					}
				}
			}
		}
		agent = state.Agents[agentIndex]
	}

	if !isAgentAndBoxTogetherLikeBros {
		return 0
	}

	// If agent is nex to box then the cost is 0
	manh := manhattenDistance(box.Coordinates, agent.Coordinates) // The agent and the box can never be in the same place
	if manh == 1 {
		return 0
	}
	cost += manh
	cost += calculateWallsCost(box.Coordinates, agent.Coordinates, state.LevelInfo)

	return cost
}
