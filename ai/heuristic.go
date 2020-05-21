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

	AStart struct{}
)

func (a AStart) Solve(levelInfo *level.Info, currentState *level.CurrentState) actions.Action {
	if len(os.Getenv("DEBUG")) > 0 {
		time.Sleep(10 * time.Second)
		communication.Log("Starting search")
	}

	lenToAllocate := len(levelInfo.WallsCoordinates) / 2
	nodesVisited := make(level.Visited, lenToAllocate)
	// add the root node to the map of the visited nodes
	nodesVisited[currentState.GetID()] = struct{}{}

	// Double linked list
	queue := list.New()
	queue.PushBack(*currentState)

	for node := queue.Front(); node != nil; node = queue.Front() {
		value := node.Value.(level.CurrentState)
		queue.Remove(node)

		if value.IsGoalState() {
			communication.Log("FOUND GOAL")
			communication.Log("Explored ", len(nodesVisited))
			communication.Log("Moves", string(value.Moves))

			return value.Moves
		}

	outer:
		for _, child := range value.Expand(nodesVisited) {
			// The only writer to the map (this happens after all goroutines are done)
			// If the above changes, this is going to lead to a race condition
			if _, ok := nodesVisited[child.ID]; !ok {
				nodesVisited[child.ID] = struct{}{}
				cost := child.Cost

				for el := queue.Front(); el != nil; el = el.Next() {
					if cost < el.Value.(level.CurrentState).Cost {
						queue.InsertBefore(*child, el)
						continue outer
					}
				}
				queue.PushBack(*child)
			}
		}
	}

	communication.Log("Explored ", len(nodesVisited), "states and failed to find a solution")
	os.Exit(1)

	return nil
}
