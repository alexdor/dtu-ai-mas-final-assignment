package ai

import (
	"github.com/alexdor/dtu-ai-mas-final-assignment/actions"
	"github.com/alexdor/dtu-ai-mas-final-assignment/level"
)

type (
	Heuristic interface {
		Solve(*level.Info, Cost) [][]actions.Action
	}

	visited level.CoordinatesLookup

	BiDirectionalBFS struct{}
)

func (b BiDirectionalBFS) Solve(levelInfo *level.Info, cost Cost) [][]actions.Action {
	panic("BiDirectionBFS not implemented")
	// lenToAllocate := len(levelInfo.WallsCoordinates) / 2
	// srcVisited := make(visited, lenToAllocate)
	// dstVisited := make(visited, lenToAllocate)

	// queue := list.New()

	// queue.PushBack()
	// // add the root node to the map of the visited nodes
	// visited[node.Id] = node

}

func isIntersection(el level.Coordinates, lookupMap visited) bool {
	_, ok := lookupMap[el]
	return ok
}
