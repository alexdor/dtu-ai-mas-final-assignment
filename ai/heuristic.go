package ai

import (
	"container/list"
	"os"
	"sort"

	"github.com/alexdor/dtu-ai-mas-final-assignment/actions"
	"github.com/alexdor/dtu-ai-mas-final-assignment/communication"
	"github.com/alexdor/dtu-ai-mas-final-assignment/config"
	"github.com/alexdor/dtu-ai-mas-final-assignment/level"
)

var (
	NodesVisited   level.Visited
	StatesExplored = 0
)

type (
	Heuristic interface {
		Solve(*level.Info, *level.CurrentState) actions.Action
	}

	AStart struct{}

	StateChildren []level.CurrentState
)

// Sort interface implementation
func (s StateChildren) Less(i, j int) bool { return s[i].Cost < s[j].Cost }
func (s StateChildren) Len() int           { return len(s) }
func (s StateChildren) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (a AStart) Solve(levelInfo *level.Info, currentState *level.CurrentState) actions.Action {
	expand := level.ExpandSingleAgent
	if len(currentState.Agents) > 1 {
		expand = level.ExpandMultiAgent
	}

	lenToAllocate := len(levelInfo.WallsCoordinates) / 2
	NodesVisited = make(level.Visited, lenToAllocate)
	currentState.CalculateCost()
	// add the root node to the map of the visited nodes
	NodesVisited[currentState.GetID()] = struct{}{}

	// Double linked list
	queue := list.New()
	queue.PushBack(*currentState)

	for node := queue.Front(); node != nil; node = queue.Front() {
		StatesExplored++
		value := node.Value.(level.CurrentState)
		queue.Remove(node)

		if value.IsGoalState() {
			communication.Log("Goal was found after exploring", StatesExplored, "states")
			communication.Log("Generated", len(NodesVisited), "states")
			if config.IsDebug {
				communication.Log("Moves", string(value.Moves))
				communication.Log("Agents Final State", value.Agents)
				communication.Log("Boxes Final State", value.Boxes)
				communication.Log("Goals", value.LevelInfo.GoalCoordinates)
			}
			return value.Moves
		}

		el := queue.Front()

		children := expand(NodesVisited, &value)
		sort.Sort(StateChildren(children))

		for _, child := range children {
			// The only writer to the map (this happens after all goroutines are done)
			// If the above changes, this is going to lead to a race condition
			if _, ok := NodesVisited[child.ID]; !ok {
				if el == nil {
					queue.PushBack(child)
					continue
				}
				NodesVisited[child.ID] = struct{}{}
				cost := child.Cost
				for {
					if cost < el.Value.(level.CurrentState).Cost {
						queue.InsertBefore(child, el)
						break
					}

					el = el.Next()
					if el == nil {
						queue.PushBack(child)
						break
					}
				}
			}
		}
	}

	communication.Log("Generated", len(NodesVisited), "states")
	communication.Log("Explored", StatesExplored, "states and failed to find a solution")
	os.Exit(1)

	return nil
}
