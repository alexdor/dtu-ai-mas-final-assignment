package level

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/alexdor/dtu-ai-mas-final-assignment/actions"
	"github.com/alexdor/dtu-ai-mas-final-assignment/communication"
)

var ErrFailedToFindBox = errors.New("Failed to find box")

type CurrentState struct {
	Boxes     []NodeOrAgent
	Agents    []NodeOrAgent
	Moves     []byte
	Cost      int
	LevelInfo *Info
	ID
}

func (c *CurrentState) GetID() ID {
	if len(c.ID) == 0 {
		var s strings.Builder
		for _, value := range append(c.Boxes, c.Agents...) {

			err := s.WriteByte(value.Letter)
			if err != nil {
				communication.Error(err)
			}

			_, err = s.WriteString(strconv.Itoa(int(value.Coordinates[0])))
			if err != nil {
				communication.Error(err)
			}

			err = s.WriteByte(',')
			if err != nil {
				communication.Error(err)
			}

			_, err = s.WriteString(strconv.Itoa(int(value.Coordinates[1])))
			if err != nil {
				communication.Error(err)
			}

			err = s.WriteByte(' ')
			if err != nil {
				communication.Error(err)
			}
		}

		c.ID = ID(s.String())
	}

	return c.ID
}

var (
	directionForCoordinates = []byte{actions.East, actions.West, actions.North, actions.South}
	coordManipulation       = []Coordinates{{0, 1}, {0, -1}, {-1, 0}, {1, 0}}
	pullOrPush              = []actions.PullOrPush{actions.Pull, actions.Push}
	wg                      = &sync.WaitGroup{}
)

func (c *CurrentState) copy(newState *CurrentState) {
	*newState = CurrentState{LevelInfo: c.LevelInfo}
	newState.Agents = make([]NodeOrAgent, len(c.Agents))
	copy(newState.Agents, c.Agents)
	newState.Boxes = make([]NodeOrAgent, len(c.Boxes))
	copy(newState.Boxes, c.Boxes)
	newState.Moves = make([]byte, len(c.Moves))
	copy(newState.Moves, c.Moves)

}

func ExpandSingleAgent(nodesInFrontier Visited, c *CurrentState) []*CurrentState {
	agentIndex := 0
	agent := c.Agents[0]
	nextStates := []*CurrentState{}
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
			addStateToStatesToExplore(&nextStates, &newState, nodesInFrontier)

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

func expandBoxMoves(state *CurrentState, nextStates *[]*CurrentState, boxCoorToMove *Coordinates, boxCoordIndex, agentIndex int, nodesVisited Visited) {
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

			addStateToStatesToExplore(nextStates, &copyOfState, nodesVisited)
		}
	}
}
func coordToDirection(oldCoord, newCoord Coordinates) actions.Direction {
	tmp := Coordinates{newCoord[0] - oldCoord[0], newCoord[1] - oldCoord[1]}
	for i := range coordManipulation {
		if coordManipulation[i] == tmp {
			return directionForCoordinates[i]
		}
	}

	panic("Failed to find direction, this should never happen")
}

func calculateCost(newState *CurrentState, nodesVisited Visited) {
	defer wg.Done()

	if _, ok := nodesVisited[newState.GetID()]; !ok {
		newState.CalculateCost()
	}
}

func addStateToStatesToExplore(nextStates *[]*CurrentState, newState *CurrentState, nodesVisited Visited) {
	wg.Add(1)

	*nextStates = append(*nextStates, newState)

	go calculateCost(newState, nodesVisited)
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

func (c *CurrentState) IsGoalState() bool {
	goalCounter := 0
outer:
	for char, goal := range c.LevelInfo.GoalCoordinates {
		for _, box := range c.Boxes {
			if char == box.Letter {
				for _, coor := range goal {
					if box.Coordinates == coor {
						goalCounter++
						if goalCounter == c.LevelInfo.GoalCount {
							continue outer
						}
					}
				}
			}
		}
		return false
	}

	return true
}

func (c *CurrentState) IsBoxAndCanMove(coor Coordinates, agentChar byte) bool {
	for _, box := range c.Boxes {
		if box.Coordinates == coor {
			return c.LevelInfo.AgentColor[agentChar] == c.LevelInfo.BoxColor[box.Letter]
		}
	}

	return false
}

func (c *CurrentState) CalculateCost() {
	if c.Cost == 0 {
		c.Cost = CalculateManhattanDistance(c)
	}
}
