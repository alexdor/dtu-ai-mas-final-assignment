package level

import (
	"github.com/alexdor/dtu-ai-mas-final-assignment/actions"
	"github.com/alexdor/dtu-ai-mas-final-assignment/communication"
)

func ExpandSingleAgent(nodesInFrontier Visited, c *CurrentState) []CurrentState {
	agentIndex := 0
	agent := c.Agents[0]
	nextStates := []CurrentState{}
	agentCoor := agent.Coordinates

	for coordIndex, move := range coordManipulation {
		newCoor := Coordinates{agentCoor[0] + move[0], agentCoor[1] + move[1]}
		if c.LevelInfo.IsWall(newCoor) {
			continue
		}

		var newState CurrentState

		c.copy(&newState)

		if c.LevelInfo.IsCellFree(newCoor, c) {
			newState.Agents[agentIndex].Coordinates = newCoor
			newState.Moves = append(newState.Moves, actions.Move(directionForCoordinates[coordIndex], actions.SingleAgentEnd)...)
			addStateToStatesToExplore(&nextStates, newState, nodesInFrontier)

			continue
		}
		// Check if the cell has a box that can be moved by this agent
		if !newState.IsBoxAndCanMove(newCoor, agent.Letter) {
			continue
		}

		expandBoxMoves(&newState, &nextStates, &newCoor, coordIndex, agentIndex, nodesInFrontier)
	}

	wg.Wait()

	return nextStates
}

func expandBoxMoves(state *CurrentState, nextStates *[]CurrentState, boxCoorToMove *Coordinates, boxCoordIndex, agentIndex int, nodesVisited Visited) {
	// Prealloc variables
	var (
		isPush bool
		coordsToCheck,
		cellToMoveInto,
		agentCoor,
		boxCoor Coordinates
	)

	boxIndex := state.FindBoxAt(*boxCoorToMove)
	currentBoxCoor := state.Boxes[boxIndex].Coordinates

	currentAgentCoord := state.Agents[agentIndex].Coordinates

	for action_i, action := range pullOrPush {
		isPush = action_i == 1

		for _, move := range coordManipulation {
			coordsToCheck = currentAgentCoord
			if isPush {
				coordsToCheck = currentBoxCoor
			}

			cellToMoveInto = Coordinates{coordsToCheck[0] + move[0], coordsToCheck[1] + move[1]}
			if !state.LevelInfo.IsCellFree(cellToMoveInto, state) {
				continue
			}

			agentCoor, boxCoor = cellToMoveInto, currentAgentCoord
			boxDirection := directionForCoordinates[boxCoordIndex]

			if isPush {
				agentCoor, boxCoor = currentBoxCoor, agentCoor
				boxDirection = coordToDirection(currentBoxCoor, boxCoor)
			}

			var copyOfState CurrentState

			state.copy(&copyOfState)

			moveAction := action(coordToDirection(currentAgentCoord, agentCoor), boxDirection, actions.SingleAgentEnd)
			copyOfState.Moves = append(copyOfState.Moves, moveAction...)
			copyOfState.Agents[agentIndex].Coordinates = agentCoor
			copyOfState.Boxes[boxIndex].Coordinates = boxCoor

			addStateToStatesToExplore(nextStates, copyOfState, nodesVisited)
		}
	}
}

func addStateToStatesToExplore(nextStates *[]CurrentState, newState CurrentState, nodesVisited Visited) {
	wg.Add(1)

	*nextStates = append(*nextStates, newState)

	go calculateCostWithGoroutine(&(*nextStates)[len(*nextStates)-1], nodesVisited)
}

func (c *CurrentState) FindBoxAt(coord Coordinates) int {
	for i, box := range c.Boxes {
		if box.Coordinates == coord {
			return i
		}
	}

	communication.Error(ErrFailedToFindBox)

	return -1
}
