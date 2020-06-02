package level

import (
	"bytes"

	"github.com/alexdor/dtu-ai-mas-final-assignment/actions"
)

var (
	noopIntent = agentIntents{action: actions.NoOp(actions.MultiAgentEnd)}
)

type agentIntents struct {
	action       actions.Action
	agentNewCoor Coordinates
	boxNewCoor   Coordinates
	boxIndex     int
}

func ExpandMultiAgent(nodesInFrontier Visited, c *CurrentState) []CurrentState {
	numOfAgents := len(c.Agents)
	wg.Add(numOfAgents)

	intents := make([][]agentIntents, numOfAgents)

	for agentIndex := range c.Agents {
		goroutineLimiter <- struct{}{}
		agentIndex := agentIndex
		go c.figureOutAgentMovements(agentIndex, intents)
	}

	hasConflict, skipAppend := false, false

	wg.Wait()
	mergedIntents := make([][]agentIntents, len(intents[0]))

	for i := 0; i < len(intents[0]); i++ {
		mergedIntents[i] = []agentIntents{intents[0][i]}
	}

	for i := 1; i < len(intents); i++ {
		localIntents := [][]agentIntents{}

		for _, firstElement := range mergedIntents {
			if len(intents[i]) == 0 {
				intents[i] = []agentIntents{noopIntent}
			}
			for _, secondElement := range intents[i] {
				skipAppend = false

				if secondElement.agentNewCoor != noopIntent.agentNewCoor {
					// Check if there is a conflict
					for _, action := range firstElement {
						hasConflict = action.agentNewCoor == secondElement.agentNewCoor ||
							action.agentNewCoor == secondElement.boxNewCoor ||
							action.boxNewCoor == secondElement.agentNewCoor ||
							action.boxNewCoor == secondElement.boxNewCoor

						if hasConflict {
							localIntents = append(localIntents, []agentIntents{action, noopIntent}, []agentIntents{noopIntent, secondElement})

							skipAppend = true
						}
					}
				}

				if !skipAppend {
					localIntents = append(localIntents, append(firstElement, secondElement))
				}
			}
		}
		mergedIntents = localIntents
	}

	nextStates := make([]CurrentState, len(mergedIntents))
	i := 0
	for _, agentIntent := range mergedIntents {

		// If all the actions are noop then skip creating them
		skip := true
		for _, action := range agentIntent {
			if !bytes.Equal(action.action, actions.NoOpAction) {
				skip = false
				break
			}
		}
		if skip {
			continue
		}
		wg.Add(1)
		goroutineLimiter <- struct{}{}
		go calcNewState(c, &nextStates[i], agentIntent, nodesInFrontier)
		i++
	}

	wg.Wait()

	return nextStates[:i]
}

func calcNewState(currentState, newState *CurrentState, currentIntents []agentIntents, nodesInFrontier Visited) {
	defer goroutineCleanupFunc()
	currentState.copy(newState)
	for j, action := range currentIntents {
		newState.Moves = append(newState.Moves, action.action...)
		if bytes.Equal(action.action, noopIntent.action) {
			continue
		}
		newState.Agents[j].Coordinates = action.agentNewCoor

		if action.boxNewCoor != noopIntent.boxNewCoor {
			newState.Boxes[action.boxIndex].Coordinates = action.boxNewCoor
		}
	}

	newState.Moves = append(newState.Moves, actions.SingleAgentEnd)

	calculateCost(newState, nodesInFrontier)
}

func goroutineCleanupFunc() {
	wg.Done()
	<-goroutineLimiter
}

func (c *CurrentState) figureOutAgentMovements(agentIndex int, intents [][]agentIntents) {
	defer goroutineCleanupFunc()

	localIntents := []agentIntents{}

	agent := c.Agents[agentIndex]
	agentCoor := agent.Coordinates

	for coordIndex, move := range coordManipulation {
		newCoor := Coordinates{agentCoor[0] + move[0], agentCoor[1] + move[1]}
		if c.LevelInfo.IsWall(newCoor) {
			continue
		}

		if c.LevelInfo.IsCellFree(newCoor, c) {
			localIntents = append(localIntents, agentIntents{
				action:       actions.Move(directionForCoordinates[coordIndex], actions.MultiAgentEnd),
				agentNewCoor: newCoor,
			})

			continue
		}
		// Check if the cell has a box that can be moved by this agent
		if !c.IsBoxAndCanMove(newCoor, agent.Letter) {
			continue
		}

		expandMABoxMoves(c, &newCoor, coordIndex, agentIndex, &localIntents)
	}

	intents[agentIndex] = localIntents
}

func expandMABoxMoves(state *CurrentState, boxCoorToMove *Coordinates, boxCoordIndex, agentIndex int, localIntents *[]agentIntents) {
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

			*localIntents = append(*localIntents, agentIntents{
				action:       action(coordToDirection(currentAgentCoord, agentCoor), boxDirection, actions.MultiAgentEnd),
				agentNewCoor: agentCoor,
				boxNewCoor:   boxCoor,
				boxIndex:     boxIndex,
			})
		}
	}
}
