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

	// Make sure agents are sorted
	sort.Slice(currentState.Agents, func(i, j int) bool {
		return currentState.Agents[i].Letter < currentState.Agents[j].Letter
	})

	wg := &sync.WaitGroup{}
	wg.Add(3)

	goalCount := 0
	inGameWalls := []level.Coordinates{}
	boxGoalAssignment := make([]level.Coordinates, len(currentState.Boxes))

	preproccessLvl(wg, &levelInfo, &currentState, &goalCount, &inGameWalls, &boxGoalAssignment)
	wg.Wait()

	levelInfo.GoalCount = goalCount
	levelInfo.InGameWallsCoordinates = inGameWalls
	levelInfo.BoxGoalAssignment = boxGoalAssignment

	return levelInfo, currentState, nil
}
func findCloserBox(coords level.Coordinates, char byte, boxes []level.NodeOrAgent, assignedBoxes map[int]struct{}) int {
	minDist := math.MaxInt64
	pos := -1

	for i, box := range boxes {
		if _, ok := assignedBoxes[i]; ok || box.Letter != char {
			continue
		}

		cost := level.ManhattanPlusPlus(coords, box.Coordinates)
		if cost < minDist {
			minDist = cost
			pos = i
		}
	}

	return pos
}

func preproccessLvl(wg *sync.WaitGroup, levelInfo *level.Info, state *level.CurrentState, goalCount *int, inGameWalls, boxGoalAssignment *[]level.Coordinates) {
	go func() {
		defer wg.Done()

		for i := range *boxGoalAssignment {
			(*boxGoalAssignment)[i] = level.Coordinates{-1, -1}
		}

		assignedBoxes := make(map[int]struct{})

		for char, goals := range levelInfo.GoalCoordinates {
			for _, coord := range goals {
				boxIndex := findCloserBox(coord, char, state.Boxes, assignedBoxes)
				assignedBoxes[boxIndex] = struct{}{}
				(*boxGoalAssignment)[boxIndex] = coord
			}
		}
	}()

	go func() {
		defer wg.Done()

		for key := range levelInfo.WallsCoordinates {
			isEdgeWall := key[0] == 0 || key[1] == 0 || key[0] == levelInfo.MaxCoord[0] || key[1] == levelInfo.MaxCoord[1]
			if !isEdgeWall {
				*inGameWalls = append(*inGameWalls, key)
			}
		}
	}()
	go func() {
		defer wg.Done()

		for _, v := range levelInfo.GoalCoordinates {
			*goalCount += len(v)
		}
	}()

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
