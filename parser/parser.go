package parser

import (
	"strings"

	"github.com/alexdor/dtu-ai-mas-final-assignment/communication"
	"github.com/alexdor/dtu-ai-mas-final-assignment/config"
	"github.com/alexdor/dtu-ai-mas-final-assignment/types"
)

func ParseLevel() (types.LevelInfo, error) {
	levelInfo := types.GetLevelInfo()
	mode := ""
	row := uint(0)

	for {
		msg, err := communication.ReadNextMessages()
		if err != nil {
			communication.Error(err)
			return levelInfo, err
		}

		msg = strings.TrimSpace(msg)
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

		parseMode(mode, msg, row, &levelInfo)
		row++
	}

	return levelInfo, nil
}

func parseMode(mode, msg string, row uint, levelInfo *types.LevelInfo) {
	cor := [2]uint{row, 0}

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
		agentOrBoxCoordinatesMap := levelInfo.BoxCoordinates

		for j := range msg {
			cor[1] = uint(j)
			char := msg[j]

			switch {
			case char == config.FreeSpaceSymbol:
				continue

			case char == config.WallsSymbol:
				levelInfo.WallsCoordinates[cor] = struct{}{}

			case '0' <= char && char <= '9':
				agentOrBoxCoordinatesMap = levelInfo.AgentCoordinates
				fallthrough

			default:
				if _, ok := agentOrBoxCoordinatesMap[char]; ok {
					agentOrBoxCoordinatesMap[char][cor] = struct{}{}
					continue
				}

				agentOrBoxCoordinatesMap[char] = types.CoordinatesLookup{cor: struct{}{}}
			}
		}

	case "goal":
		for j := range msg {
			char := msg[j]
			if char != config.FreeSpaceSymbol && char != config.WallsSymbol {
				if _, ok := levelInfo.GoalCoordinates[char]; ok {
					cor[1] = uint(j)
					levelInfo.GoalCoordinates[char][cor] = struct{}{}

					continue
				}

				levelInfo.GoalCoordinates[char] = types.CoordinatesLookup{cor: struct{}{}}
			}
		}

	default:
		levelInfo.LevelInfo[mode] = msg
	}
}
