package ai

import (
	"github.com/alexdor/dtu-ai-mas-final-assignment/actions"
	"github.com/alexdor/dtu-ai-mas-final-assignment/types"
)

type Heuristic interface {
	Solve(types.LevelInfo, Cost) [][]actions.Action
}

type AStar struct{}

func (AStar) Solve(levelInfo types.LevelInfo, cost Cost) [][]actions.Action {
	panic("Astar isn't implimented")
}
