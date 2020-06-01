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

func ExpandMultiAgent(nodesInFrontier Visited, c *CurrentState) []*CurrentState {
	numOfAgents := len(c.Agents)
	wg.Add(numOfAgents)

	intents := make([][]agentIntents, numOfAgents)

	for agentIndex := range c.Agents {
		agentIndex := agentIndex
		go c.figureOutAgentMovements(agentIndex, intents)
	}

	nextStates := []*CurrentState{}
	isLastIntent, hasConflict, skipAppend := false, false, false

	wg.Wait()
	mergedIntents := make([][]agentIntents, len(intents[0]))

	for i := 0; i < len(intents[0]); i++ {
		mergedIntents[i] = []agentIntents{intents[0][i]}
	}

	for i := 1; i < len(intents); i++ {
		isLastIntent = i == len(intents)-1
		localIntents := [][]agentIntents{}

		for _, firstElement := range mergedIntents {
			if len(intents[i]) == 0 {
				intents[i] = []agentIntents{noopIntent}
			}
		outer:
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

				// If last intent calculate next states
				if isLastIntent {
					// If all the actions are noop then skip creating them
					if bytes.Equal(secondElement.action, actions.NoOpAction) {
						skip := true
						for _, action := range firstElement {
							if !bytes.Equal(action.action, actions.NoOpAction) {
								skip = false
								break
							}
						}
						if skip {
							continue outer
						}
					}

					var newState CurrentState
					nextStates = append(nextStates, &newState)
					// One for the new state and one for the cost
					wg.Add(2)
					go calcNewState(c, &newState, firstElement, secondElement, nodesInFrontier, i)
				}
			}
		}
		mergedIntents = localIntents
	}

	wg.Wait()

	return nextStates
}

func calcNewState(currentState, newState *CurrentState, currentIntents []agentIntents, finalIntent agentIntents, nodesInFrontier Visited, agentIndex int) {
	defer wg.Done()
	currentState.copy(newState)
	for j, action := range currentIntents {
		if bytes.Equal(action.action, noopIntent.action) {
			continue
		}
		newState.Agents[j].Coordinates = action.agentNewCoor
		newState.Moves = append(newState.Moves, action.action...)

		if action.boxNewCoor != noopIntent.boxNewCoor {
			newState.Boxes[action.boxIndex].Coordinates = action.boxNewCoor
		}
	}

	if !bytes.Equal(finalIntent.action, noopIntent.action) {
		newState.Agents[agentIndex].Coordinates = finalIntent.agentNewCoor
		if finalIntent.boxNewCoor != noopIntent.boxNewCoor {
			newState.Boxes[finalIntent.boxIndex].Coordinates = finalIntent.boxNewCoor
		}
	}

	newState.Moves = append(newState.Moves, finalIntent.action...)
	newState.Moves = append(newState.Moves, actions.SingleAgentEnd)

	calculateCost(newState, nodesInFrontier)
}

func (c *CurrentState) figureOutAgentMovements(agentIndex int, intents [][]agentIntents) {
	defer wg.Done()

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
