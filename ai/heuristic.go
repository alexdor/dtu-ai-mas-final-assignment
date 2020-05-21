package ai

import (
	"container/list"
	"os"
	"time"

	"github.com/alexdor/dtu-ai-mas-final-assignment/actions"
	"github.com/alexdor/dtu-ai-mas-final-assignment/communication"
	"github.com/alexdor/dtu-ai-mas-final-assignment/level"
)

type (
	Heuristic interface {
		Solve(*level.Info, *level.CurrentState) actions.Action
	}

	visited map[level.ID]struct{}

	BiDirectionalBFS struct{}

	AStart struct{}
)

// func (b BiDirectionalBFS) Solve(levelInfo *level.Info, currentState level.CurrentState, cost Cost) [][]actions.Action {
// 	panic("smth")
// 	lenToAllocate := len(levelInfo.WallsCoordinates) / 2
// 	srcVisited := make(visited, lenToAllocate)
// 	dstVisited := make(visited, lenToAllocate)
// 	queue := list.New()
// 	queue.PushBack(currentState)
// 	// add the root node to the map of the visited nodes
// 	visited[currentState.GetID()] = currentState
// }

func (a AStart) Solve(levelInfo *level.Info, currentState *level.CurrentState) actions.Action {
	if len(os.Getenv("DEBUG")) > 0 {
		time.Sleep(10 * time.Second)
		communication.Log("Starting")
	}
	// Double linked list
	queue := list.New()
	queue.PushBack(*currentState)

	lenToAllocate := len(levelInfo.WallsCoordinates) / 2
	nodesVisited := make(visited, lenToAllocate)
	// add the root node to the map of the visited nodes
	nodesVisited[currentState.GetID()] = struct{}{}

	for node := queue.Front(); node != nil; node = node.Next() {
		value := node.Value.(level.CurrentState)
		if value.IsGoalState() {
			communication.Log("FOUND GOAL")
			communication.Log(string(value.Moves))
			return value.Moves
		}

		for _, child := range value.Expand() {
			if _, ok := nodesVisited[child.GetID()]; !ok {
				nodesVisited[child.GetID()] = struct{}{}
				cost := child.Cost

				placed := false

				for el := queue.Front(); el != nil; el = el.Next() {
					if cost < el.Value.(level.CurrentState).Cost {
						placed = true

						queue.InsertBefore(child, el)
					}
				}

				if !placed {
					queue.PushBack(child)
				}
			}
		}
	}
	panic("My ass is on fire")
}

func isIntersection(id level.ID, lookupMap visited) bool {
	_, ok := lookupMap[id]
	return ok
}
