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

	// Create a list to hold possible actions for each agent
	intents := make([][]agentIntents, numOfAgents)

	// Calculate actions for each agent
	for agentIndex := range c.Agents {
		goroutineLimiter <- struct{}{}
		go c.figureOutAgentMovements(agentIndex, &intents[agentIndex])
	}

	hasConflict := false

	wg.Wait()

	// Start merging all agent actions
	// TODO: merge things in parallel
	mergedIntents := make([][]agentIntents, len(intents[0]))

	// Copy actions from 1st agent to the merged intents
	for i := 0; i < len(intents[0]); i++ {
		mergedIntents[i] = []agentIntents{intents[0][i]}
	}

	if len(mergedIntents) == 0 {
		mergedIntents = [][]agentIntents{{noopIntent}}
	}

	// Merge the actions for the rest of the agents on the first agent
	for currentAgentIndex := 1; currentAgentIndex < len(intents); currentAgentIndex++ {
		// tmp list to hold the new merged intents
		localIntents := [][]agentIntents{}

		// If the current agent doesn't have any actions, add noop action
		if len(intents[currentAgentIndex]) == 0 {
			intents[currentAgentIndex] = []agentIntents{noopIntent}
		}
		for _, intentsFromOtherAgents := range mergedIntents {

			for _, intentToAdd := range intents[currentAgentIndex] {
				hasConflict = !bytes.Equal(intentToAdd.action, noopIntent.action)
				if hasConflict {
					for i, action := range intentsFromOtherAgents {
						// Check if there is a conflict and handle it
						hasConflict = !bytes.Equal(action.action, noopIntent.action) &&
							(action.agentNewCoor == intentToAdd.agentNewCoor ||
								action.boxNewCoor == intentToAdd.agentNewCoor ||
								(intentToAdd.boxNewCoor != noopIntent.boxNewCoor &&
									(action.agentNewCoor == intentToAdd.boxNewCoor ||
										action.boxNewCoor == intentToAdd.boxNewCoor)))

						if hasConflict {
							newIntent := append(intentsFromOtherAgents, intentToAdd)
							newIntent[i] = noopIntent
							localIntents = append(localIntents, append(intentsFromOtherAgents, noopIntent), newIntent)
							break
						}
					}
				}
				// If we didn't find a conflict, then just merge the moves
				if !hasConflict {
					localIntents = append(localIntents, append(intentsFromOtherAgents, intentToAdd))
				}
			}
		}
		mergedIntents = localIntents
	}

	nextStates := make([]CurrentState, len(mergedIntents))
	statesCreated := 0
	for _, intent := range mergedIntents {
		waitGoroutineToFreeUp()
		go calcNewState(c, &nextStates[statesCreated], intent, nodesInFrontier)
		statesCreated++
	}

	wg.Wait()

	return nextStates[:statesCreated]
}

func calcNewState(currentState, newState *CurrentState, currentIntents []agentIntents, nodesInFrontier Visited) {
	defer cleanupAfterGoroutine()
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

func (c *CurrentState) figureOutAgentMovements(agentIndex int, intentToUpdate *[]agentIntents) {
	defer cleanupAfterGoroutine()

	agent := c.Agents[agentIndex]

	if c.LevelInfo.ZeroInGameWalls {
		everythinInGoal := true
		for _, boxIndex := range c.LevelInfo.AgentBoxAssignment[agent.Letter] {
			if c.LevelInfo.BoxGoalAssignment[boxIndex] != c.Boxes[boxIndex].Coordinates {
				everythinInGoal = false
				break
			}
		}
		if everythinInGoal {
			return
		}
	}

	agentCoor := agent.Coordinates
	localIntents := []agentIntents{}

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

	*intentToUpdate = localIntents
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

	boxIndex := state.findBoxAt(*boxCoorToMove)
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
