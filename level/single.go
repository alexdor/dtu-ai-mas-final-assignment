package level

import (
	"github.com/alexdor/dtu-ai-mas-final-assignment/actions"
)

func ExpandSingleAgent(nodesInFrontier Visited, c *CurrentState) []CurrentState {
	agent := c.Agents[0]
	nextStates := make([]CurrentState, len(coordManipulation)*4)
	i := 0
	agentCoor := agent.Coordinates

	for coordIndex, move := range coordManipulation {
		newCoor := Coordinates{agentCoor[0] + move[0], agentCoor[1] + move[1]}
		if c.LevelInfo.IsWall(newCoor) {
			continue
		}

		if c.LevelInfo.IsCellFree(newCoor, c) {
			waitGoroutineToFreeUp()
			go calculateAgentMove(c, &nextStates[i], &newCoor, coordIndex, nodesInFrontier)
			i++
			continue
		}
		// Check if the cell has a box that can be moved by this agent
		if !c.IsBoxAndCanMove(newCoor, agent.Letter) {
			continue
		}

		expandBoxMoves(c, nextStates, &newCoor, coordIndex, nodesInFrontier, &i)
	}

	wg.Wait()

	return nextStates[:i]
}

func calculateAgentMove(currentState, newState *CurrentState, newCoor *Coordinates, coordIndex int, nodesVisited Visited) {
	defer cleanupAfterGoroutine()
	currentState.copy(newState)
	newState.Agents[0].Coordinates = *newCoor
	newState.Moves = append(newState.Moves, actions.Move(directionForCoordinates[coordIndex], actions.SingleAgentEnd)...)
	calculateCost(newState, nodesVisited)
}

func expandBoxMoves(state *CurrentState, nextStates []CurrentState, boxCoorToMove *Coordinates, boxCoordIndex int, nodesVisited Visited, i *int) {
	// Prealloc variables
	var (
		isPush bool
		coordsToCheck,
		cellToMoveInto,
		agentCoor,
		boxCoor Coordinates
	)

	boxIndex := state.findBoxAt(*boxCoorToMove)
	currentBoxCoor := state.Boxes[boxIndex].Coordinates

	currentAgentCoord := state.Agents[0].Coordinates

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
			waitGoroutineToFreeUp()
			go calculateAgentBoxMove(
				boxMoveCalc{
					currentState:      state,
					newState:          &nextStates[*i],
					nodesVisited:      nodesVisited,
					action:            action,
					agentCoor:         agentCoor,
					boxCoor:           boxCoor,
					boxDirection:      boxDirection,
					boxIndex:          boxIndex,
					currentAgentCoord: currentAgentCoord,
				})
			*i++
		}
	}
}

type boxMoveCalc struct {
	currentState, newState                *CurrentState
	currentAgentCoord, agentCoor, boxCoor Coordinates
	boxDirection                          byte
	boxIndex                              int
	nodesVisited                          Visited
	action                                actions.PullOrPush
}

func calculateAgentBoxMove(params boxMoveCalc) {
	defer cleanupAfterGoroutine()
	params.currentState.copy(params.newState)

	moveAction := params.action(coordToDirection(params.currentAgentCoord, params.agentCoor), params.boxDirection, actions.SingleAgentEnd)
	params.newState.Moves = append(params.newState.Moves, moveAction...)
	params.newState.Agents[0].Coordinates = params.agentCoor
	params.newState.Boxes[params.boxIndex].Coordinates = params.boxCoor

	calculateCost(params.newState, params.nodesVisited)

}
