package ai

import (
	"github.com/alexdor/dtu-ai-mas-final-assignment/actions"
	"github.com/alexdor/dtu-ai-mas-final-assignment/level"
)

func Play(levelInfo *level.Info, heuristic Heuristic, costCalculator Cost) {
	solution := heuristic.Solve(levelInfo, costCalculator)
	for _, actionSet := range solution {
		actions.ExecuteActions(actionSet...)
	}
}
