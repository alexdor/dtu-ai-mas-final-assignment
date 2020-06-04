package level

import (
	"errors"
	"runtime"
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
		generateID(&s, c.Boxes)
		generateID(&s, c.Agents)
		c.ID = ID(s.String())
	}

	return c.ID
}

func generateID(s *strings.Builder, agentOrBox []NodeOrAgent) {
	for _, value := range agentOrBox {

		_, err := s.WriteString(strconv.Itoa(int(value.Coordinates[0])))
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
}

var (
	directionForCoordinates = []byte{actions.East, actions.West, actions.North, actions.South}
	coordManipulation       = []Coordinates{{0, 1}, {0, -1}, {-1, 0}, {1, 0}}
	pullOrPush              = []actions.PullOrPush{actions.Pull, actions.Push}
	wg                      = &sync.WaitGroup{}
	goroutineLimiter        = make(chan struct{}, runtime.NumCPU())
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

func (c *CurrentState) findBoxAt(coord Coordinates) int {
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
	for goalChar, goal := range c.LevelInfo.GoalCoordinates {
	inner:
		for _, goalCoor := range goal {
			for _, box := range c.Boxes {
				if goalChar == box.Letter && box.Coordinates == goalCoor {
					goalCounter++
					if goalCounter == c.LevelInfo.GoalCount {
						break outer
					}
					continue inner
				}
			}
			return false
		}
	}

	return true
}

func (c *CurrentState) IsBoxAndCanMove(coor Coordinates, agentChar byte) bool {
	for _, boxIndex := range c.LevelInfo.AgentBoxAssignment[agentChar] {
		isBoxAtCoor := coor == c.Boxes[boxIndex].Coordinates
		if isBoxAtCoor {
			return true
		}
	}

	return false
}

func (c *CurrentState) CalculateCost() {
	if c.Cost == 0 {
		c.Cost = CalculateAggregatedCost(c)
	}
}

func cleanupAfterGoroutine() {

	wg.Done()
	<-goroutineLimiter
}

func waitGoroutineToFreeUp() {
	wg.Add(1)
	goroutineLimiter <- struct{}{}
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

	if _, ok := nodesVisited[newState.GetID()]; !ok {
		newState.CalculateCost()
	}
}
