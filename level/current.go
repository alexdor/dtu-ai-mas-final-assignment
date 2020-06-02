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

func calculateCostWithGoroutine(newState *CurrentState, nodesVisited Visited) {
	defer wg.Done()
	calculateCost(newState, nodesVisited)
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
	for _, box := range c.Boxes {
		if box.Coordinates == coor {
			return c.LevelInfo.AgentColor[agentChar] == c.LevelInfo.BoxColor[box.Letter]
		}
	}

	return false
}

func (c *CurrentState) CalculateCost() {
	if c.Cost == 0 {
		c.Cost = CalculateAggregatedCost(c)
	}
}
