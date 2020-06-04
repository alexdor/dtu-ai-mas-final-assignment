package parser

import (
	"math"
	"sort"
	"strings"
	"sync"

	"github.com/alexdor/dtu-ai-mas-final-assignment/communication"
	"github.com/alexdor/dtu-ai-mas-final-assignment/config"
	"github.com/alexdor/dtu-ai-mas-final-assignment/level"
)

func ParseLevel() (level.Info, level.CurrentState, error) {
	levelInfo := level.Info{}
	levelInfo.Init()

	currentState := level.CurrentState{LevelInfo: &levelInfo}
	mode := ""
	row := level.Point(0)
	col := 0

	for {
		msg, err := communication.ReadNextMessages()
		if err != nil {
			communication.Error(err)
			return levelInfo, currentState, err
		}

		msg = strings.TrimSpace(msg)
		if len(msg)-1 > col {
			col = len(msg) - 1
		}

		// Handle mode switching
		if strings.HasPrefix(msg, "#") {
			msg = strings.TrimPrefix(msg, "#")
			if msg == "end" {
				break
			}

			mode = msg
			row = 0

			continue
		}

		parseMode(mode, msg, row, &levelInfo, &currentState)
		row++
	}

	levelInfo.MaxCoord = level.Coordinates{row, col}
	levelInfo.IsSingleAgent = len(currentState.Agents) == 1

	preprocessLvl(&levelInfo, &currentState)

	return levelInfo, currentState, nil
}

func findCloserBox(coords level.Coordinates, char byte, boxes []level.NodeOrAgent, assignedBoxes map[int]struct{}, state *level.CurrentState) int {
	minDist := math.MaxInt64
	pos := -1

	for i, box := range boxes {
		if _, ok := assignedBoxes[i]; ok || box.Letter != char {
			continue
		}

		cost := level.ManhattanPlusPlus(coords, box.Coordinates, state, &box, i)
		if cost < minDist {
			minDist = cost
			pos = i
		}
	}

	return pos
}

func preprocessLvl(levelInfo *level.Info, state *level.CurrentState) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Make sure agents are sorted
	go func() {
		defer wg.Done()
		sort.Slice(state.Agents, func(i, j int) bool {
			return state.Agents[i].Letter < state.Agents[j].Letter
		})
	}()

	var moveableBoxes []level.NodeOrAgent
	go func() {
		defer wg.Done()
		for _, box := range state.Boxes {
			boxColor := levelInfo.BoxColor[box.Letter]
			isBoxColorMoveable := false

			for _, agentColor := range levelInfo.AgentColor {
				isBoxAndAgentColorEqual := agentColor == boxColor
				if isBoxAndAgentColorEqual {
					isBoxColorMoveable = true
					break
				}
			}

			if !isBoxColorMoveable {
				levelInfo.WallsCoordinates[box.Coordinates] = struct{}{}
				delete(levelInfo.GoalCoordinates, box.Letter)
				continue
			}

			moveableBoxes = append(moveableBoxes, box)
		}
	}()

	goalCount := 0
	inGameWalls := []level.Coordinates{}
	boxGoalAssignment := make([]level.Coordinates, len(state.Boxes))
	agentBoxAssignment := make(level.AgentToBoxesLookup)
	boxIndexToAgentIndex := make(level.IndexToIndexMapping)

	wg.Wait()
	state.Boxes = moveableBoxes

	wg.Add(2)

	wallRows := make(level.ContinuosWalls)
	go computeInGameWallsAndWallRows(wg, levelInfo, &inGameWalls, wallRows)

	go func() {
		defer wg.Done()

		for i := range boxGoalAssignment {
			boxGoalAssignment[i] = level.Coordinates{-1, -1}
		}

		assignedBoxes := make(map[int]struct{})

		for char, goals := range levelInfo.GoalCoordinates {
			for _, coord := range goals {

				boxIndex := findCloserBox(coord, char, state.Boxes, assignedBoxes, state)
				if boxIndex == -1 {
					continue
				}
				assignedBoxes[boxIndex] = struct{}{}

				boxGoalAssignment[boxIndex] = coord
			}
		}
	}()

	for _, v := range levelInfo.GoalCoordinates {
		goalCount += len(v)
	}

	wg.Wait()

	assignAgentsToBoxes(levelInfo, state, boxGoalAssignment, agentBoxAssignment, boxIndexToAgentIndex)

	levelInfo.GoalCount = goalCount
	levelInfo.InGameWallsCoordinates = inGameWalls
	levelInfo.BoxGoalAssignment = boxGoalAssignment
	levelInfo.AgentBoxAssignment = agentBoxAssignment
	levelInfo.WallRows = wallRows
	levelInfo.BoxIndexToAgentIndex = boxIndexToAgentIndex
}

func computeInGameWallsAndWallRows(wg *sync.WaitGroup, levelInfo *level.Info, storeInGameWalls *[]level.Coordinates, wallRows level.ContinuosWalls) {
	defer wg.Done()
	var inGameWalls []level.Coordinates

	for key := range levelInfo.WallsCoordinates {
		isEdgeWall := key[0] == 0 || key[1] == 0 || key[0] == levelInfo.MaxCoord[0] || key[1] == levelInfo.MaxCoord[1]
		if !isEdgeWall {
			inGameWalls = append(inGameWalls, key)
		}
	}

	// Sort inGameWalls based on row and then based on column
	sort.Slice(inGameWalls, func(i, j int) bool {
		return inGameWalls[i][0] < inGameWalls[j][0] || (inGameWalls[i][0] == inGameWalls[j][0] && inGameWalls[i][1] < inGameWalls[j][1])
	})
	figureOutWallRows(inGameWalls, wallRows)

	*storeInGameWalls = inGameWalls
}

func figureOutWallRows(inGameWalls []level.Coordinates, wallRows level.ContinuosWalls) {
outer:
	for _, wall := range inGameWalls {
		rowToAppend, ok := wallRows[wall[0]]
		// Create new key value pair if it doesn't exits
		if !ok {
			wallRows[wall[0]] = level.ContinuosWallCoord{{wall[1], wall[1]}}
			continue
		}

		for i := range rowToAppend {
			// If the current wall is part of the previous box
			// increase tha "max" var of the wall by one
			if wall[1]-1 == rowToAppend[i][1] {
				wallRows[wall[0]][i][1] = rowToAppend[i][1] + 1
				continue outer
			}
		}

		// If the walls isn't connected to previous walls then create a new set
		wallRows[wall[0]] = append(rowToAppend, [2]level.Point{wall[1], wall[1]})

	}
}

func assignAgentsToBoxes(levelInfo *level.Info, state *level.CurrentState, boxGoalAssignment []level.Coordinates, agentBoxAssignment level.AgentToBoxesLookup, boxIndexToAgentIndex level.IndexToIndexMapping) {
	if levelInfo.IsSingleAgent {
		return
	}

	for boxIndex, box := range state.Boxes {
		boxColor := levelInfo.BoxColor[box.Letter]
		minCost := math.MaxInt64

		var agentWithLowestCost byte
		var indexOfAgentWithLowestCost int
		for i, agent := range state.Agents {
			isBoxMoveableByAgent := levelInfo.AgentColor[agent.Letter] == boxColor
			if !isBoxMoveableByAgent {
				continue
			}

			coordsToConsider := agent.Coordinates

			boxesAssignedToAgent, isAnotherBoxAssignedToAgent := agentBoxAssignment[agent.Letter]
			if isAnotherBoxAssignedToAgent {
				coordsToConsider = boxGoalAssignment[boxesAssignedToAgent[len(boxesAssignedToAgent)-1]]
			}

			cost := level.ManhattanPlusPlus(coordsToConsider, box.Coordinates, state, &box, boxIndex)

			if cost < minCost {
				minCost = cost
				agentWithLowestCost = agent.Letter
				indexOfAgentWithLowestCost = i
			}
		}

		agentBoxAssignment[agentWithLowestCost] = append(agentBoxAssignment[agentWithLowestCost], boxIndex)
		boxIndexToAgentIndex[boxIndex] = indexOfAgentWithLowestCost
	}
}

func parseMode(mode, msg string, row level.Point, levelInfo *level.Info, currentState *level.CurrentState) {
	cor := level.Coordinates{row, 0}

	switch mode {
	case "colors":
		colors := strings.Split(msg, ":")

		for _, letter := range strings.Split(colors[1], ",") {
			char := strings.TrimSpace(letter)[0]
			// Check if the char is a number
			if '0' <= char && char <= '9' {
				levelInfo.AgentColor[char] = colors[0]
				continue
			}

			levelInfo.BoxColor[char] = colors[0]
		}

	case "initial":
		for j := range msg {
			cor[1] = level.Point(j)
			char := msg[j]
			agentOrBoxCoordinatesMap := &currentState.Boxes

			switch {
			case char == config.FreeSpaceSymbol:
				continue
			case char == config.WallsSymbol:
				levelInfo.WallsCoordinates[cor] = struct{}{}
			case '0' <= char && char <= '9':
				agentOrBoxCoordinatesMap = &currentState.Agents
				fallthrough
			default:
				*agentOrBoxCoordinatesMap = append(*agentOrBoxCoordinatesMap,
					level.NodeOrAgent{
						Letter:      char,
						Coordinates: cor,
					})
			}
		}

	case "goal":
		for j := range msg {
			char := msg[j]
			if char != config.FreeSpaceSymbol && char != config.WallsSymbol {
				cor[1] = level.Point(j)

				if _, ok := levelInfo.GoalCoordinates[char]; ok {
					levelInfo.GoalCoordinates[char] = append(levelInfo.GoalCoordinates[char], cor)

					continue
				}

				levelInfo.GoalCoordinates[char] = []level.Coordinates{cor}
			}
		}

	default:
		levelInfo.LevelInfo[mode] = msg
	}

}
