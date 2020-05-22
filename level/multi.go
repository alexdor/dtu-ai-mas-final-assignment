package level

import (
	"sync"

	"github.com/alexdor/dtu-ai-mas-final-assignment/actions"
)

var (
	exploreMutex = sync.Mutex{}
	noopIntent   = agentIntents{action: actions.NoOp(actions.MultiAgentEnd)}
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
		go c.figureOutAgentMovements(agentIndex, &intents)
	}

	nextStates := []*CurrentState{}
	mergedIntents := [][]agentIntents{}
	isLastIntent, isThereAConflict, skipAppend := false, false, false

	wg.Wait()

	for _, intent := range intents[0] {
		mergedIntents = append(mergedIntents, []agentIntents{intent})
	}

	for i := 1; i < len(intents); i++ {
		isLastIntent = i == len(intents)-1
		localIntents := [][]agentIntents{}

		for _, firstElement := range mergedIntents {
			for _, secondElement := range intents[i] {
				skipAppend = false

				if secondElement.agentNewCoor != noopIntent.agentNewCoor {
					// Check if there is a conflict
					for _, action := range firstElement {
						isThereAConflict = action.agentNewCoor == secondElement.agentNewCoor ||
							action.agentNewCoor == secondElement.boxNewCoor ||
							action.boxNewCoor == secondElement.agentNewCoor ||
							action.boxNewCoor == secondElement.boxNewCoor

						if isThereAConflict {
							localIntents = append(
								localIntents,
								[]agentIntents{action, noopIntent},
								[]agentIntents{noopIntent, secondElement},
							)

							skipAppend = true
						}
					}
				}

				if skipAppend {
					localIntents = append(localIntents, append(firstElement, secondElement))
				}
				// If last intent calculate next states
				if isLastIntent {
					var newState CurrentState

					c.copy(&newState)

					for j, action := range firstElement {
						newState.Agents[j].Coordinates = action.agentNewCoor
						newState.Moves = append(newState.Moves, action.action...)

						if action.boxNewCoor != noopIntent.boxNewCoor {
							newState.Boxes[action.boxIndex].Coordinates = action.boxNewCoor
						}
					}

					newState.Agents[i].Coordinates = secondElement.agentNewCoor
					if secondElement.boxNewCoor != noopIntent.boxNewCoor {
						newState.Boxes[secondElement.boxIndex].Coordinates = secondElement.boxNewCoor
					}
					newState.Moves = append(newState.Moves, secondElement.action...)
					// TODO: Create state
					nextStates = append(nextStates, &newState)

					wg.Add(1)

					go calculateCost(&newState, nodesInFrontier)
				}
			}
		}

		mergedIntents = localIntents
	}

	wg.Wait()

	return nextStates
}

func (c *CurrentState) figureOutAgentMovements(agentIndex int, intents *[][]agentIntents) {
	defer wg.Done()

	localIntents := []agentIntents{}
	agent := c.Agents[agentIndex]
	agentCoor := agent.Coordinates
	ending := actions.MultiAgentEnd

	if agentIndex == len(*intents)-1 {
		ending = actions.SingleAgentEnd
	}

	for coordIndex, move := range coordManipulation {
		newCoor := Coordinates{agentCoor[0] + move[0], agentCoor[1] + move[1]}
		if c.LevelInfo.IsWall(newCoor) {
			continue
		}

		if c.LevelInfo.IsCellFree(newCoor, c) {
			localIntents = append(localIntents, agentIntents{
				action:       actions.Move(directionForCoordinates[coordIndex], ending),
				agentNewCoor: newCoor,
			})

			continue
		}
		// Check if the cell has a box that can be moved by this agent
		if !c.IsBoxAndCanMove(newCoor, agent.Letter) {
			continue
		}

		expandMABoxMoves(c, &newCoor, coordIndex, agentIndex, &localIntents, ending)
	}

	// If no possible moves where found then add a noop
	if len(localIntents) == 0 {
		localIntents = append(localIntents, agentIntents{
			action: actions.NoOp(ending),
		})
	}

	exploreMutex.Lock()
	(*intents)[agentIndex] = localIntents
	exploreMutex.Unlock()
}

func expandMABoxMoves(state *CurrentState, boxCoorToMove *Coordinates, boxCoordIndex, agentIndex int, localIntents *[]agentIntents, ending byte) {
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
				action:       action(coordToDirection(currentAgentCoord, agentCoor), boxDirection, ending),
				agentNewCoor: agentCoor,
				boxNewCoor:   boxCoor,
				boxIndex:     boxIndex,
			})
		}
	}
}
