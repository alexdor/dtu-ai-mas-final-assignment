package level

import (
	"encoding/binary"
	"errors"
	"runtime"
	"sync"
	"unsafe"

	"github.com/alexdor/dtu-ai-mas-final-assignment/actions"
	"github.com/alexdor/dtu-ai-mas-final-assignment/communication"
	"github.com/alexdor/dtu-ai-mas-final-assignment/config"
)

var (
	ErrFailedToFindBox = errors.New("Failed to find box")
)

type CurrentState struct {
	Boxes     []NodeOrAgent
	Agents    []NodeOrAgent
	Moves     []byte
	Cost      int
	LevelInfo *Info
	ID
}

// This is an unsafe convertion of a byte array to string
// it's copied from go's string builder implementation
// https://golang.org/src/strings/builder.go#L46
func unsafeByteArrayToID(id []byte) ID {
	return ID(*(*string)(unsafe.Pointer(&id)))
}

func (c *CurrentState) GetID() ID {
	if len(c.ID) == 0 {
		id := make([]byte, c.LevelInfo.TotalBytesForID)

		generateID(id[:c.LevelInfo.BytesUsedForBoxes], c.Boxes)
		generateID(id[c.LevelInfo.BytesUsedForBoxes:], c.Agents)

		c.ID = unsafeByteArrayToID(id)
	}

	return c.ID
}

func generateID(id []byte, agentOrBox []NodeOrAgent) {
	var startIndex int
	for i, value := range agentOrBox {
		startIndex = 1 + i*config.BytesUsedForEachAgentOrBox
		if i == 0 {
			startIndex = 0
		}
		binary.LittleEndian.PutUint32(id[startIndex:], uint32(value.Coordinates[0]))

		binary.LittleEndian.PutUint32(id[startIndex+config.BytesUsedForEachPoint:], uint32(value.Coordinates[1]))

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
