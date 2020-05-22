package level

import (
	"bytes"
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
	isLastIntent, isThereAConflict, skipAppend := false, false, false

	wg.Wait()
	mergedIntents := make([][]agentIntents, len(intents[0]))

	for i := 0; i < len(intents[0]); i++ {
		mergedIntents[i] = []agentIntents{intents[0][i]}
	}

	for i := 1; i < len(intents); i++ {
		isLastIntent = i == len(intents)-1
		localIntents := make([][]agentIntents, len(mergedIntents)*len(intents[i]))

		for mergeIndex, firstElement := range mergedIntents {
		outer:
			for currentIndex, secondElement := range intents[i] {
				skipAppend = false

				if secondElement.agentNewCoor != noopIntent.agentNewCoor {
					// Check if there is a conflict
					for _, action := range firstElement {
						isThereAConflict = action.agentNewCoor == secondElement.agentNewCoor ||
							action.agentNewCoor == secondElement.boxNewCoor ||
							action.boxNewCoor == secondElement.agentNewCoor ||
							action.boxNewCoor == secondElement.boxNewCoor

						if isThereAConflict {
							localIntents[mergeIndex+currentIndex] = []agentIntents{action, noopIntent}
							localIntents = append(localIntents, []agentIntents{noopIntent, secondElement})

							skipAppend = true
						}
					}
				}

				if !skipAppend {
					localIntents[mergeIndex+currentIndex] = append(firstElement, secondElement)
				}
				// If last intent calculate next states
				if isLastIntent {
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

	ending := actions.MultiAgentEnd

	if agentIndex == len(*intents)-1 {
		ending = actions.SingleAgentEnd
	}

	localIntents := []agentIntents{{
		action: actions.NoOp(ending),
	}}

	agent := c.Agents[agentIndex]
	agentCoor := agent.Coordinates

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
