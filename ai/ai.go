package ai

import (
	"github.com/alexdor/dtu-ai-mas-final-assignment/actions"
	"github.com/alexdor/dtu-ai-mas-final-assignment/level"
)

func Play(levelInfo *level.Info, currentState *level.CurrentState, heuristic Heuristic) {
	solution := heuristic.Solve(levelInfo, currentState)
	actions.ExecuteActions(solution)
}
