package main

import (
	"strings"

	"github.com/alexdor/dtu-ai-mas-final-assigment/communication"
	"github.com/alexdor/dtu-ai-mas-final-assigment/config"
)

type simpleMap map[byte]string

type coordinatesLookup map[[2]int]struct{}

type intrestingCoordinates map[byte]coordinatesLookup

var (
	lvlInfo          = make(map[string]string, 2)
	agentColor       = simpleMap{}
	boxColor         = simpleMap{}
	walls            = make(map[[2]int]struct{}, 15)
	agentCoordinates = intrestingCoordinates{}
	boxCoordinates   = intrestingCoordinates{}
	goalCoordinates  = intrestingCoordinates{}
)

func main() {
	communication.Init()
	parseLvl()
	printThings()
}

func parseLvl() {
	mode := ""
	row := 0
	for {
		msg, err := communication.ReadNextMessages()
		if err != nil {
			communication.Error(err)
			break
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

		// Parse
		switch mode {
		case "colors":
			colors := strings.Split(msg, ":")

			for _, letter := range strings.Split(colors[1], ",") {
				char := strings.TrimSpace(letter)[0]
				// Check if the char is a number
				if '0' <= char && char <= '9' {
					agentColor[char] = colors[0]
					continue
				}
				boxColor[char] = colors[0]
			}

		case "initial":
			cor := [2]int{row, 0}
			communication.Log(row)
			varToAssign := boxCoordinates
			for j := range msg {
				cor[1] = j
				ch := msg[j]
				switch {
				case ch == config.FreeSpaceSymbol:
					continue

				case ch == config.WallsSymbol:
					walls[cor] = struct{}{}

				case '0' <= ch && ch <= '9':
					varToAssign = agentCoordinates
					fallthrough

				default:
					if _, ok := varToAssign[ch]; ok {
						varToAssign[ch][cor] = struct{}{}
						continue
					}
					varToAssign[ch] = coordinatesLookup{cor: struct{}{}}
				}
			}

		case "goal":
			for j := range msg {
				ch := msg[j]
				if ch != config.FreeSpaceSymbol && ch != config.WallsSymbol {
					cor := [2]int{row, j}
					if _, ok := goalCoordinates[ch]; ok {
						goalCoordinates[ch][cor] = struct{}{}
						continue
					}
					goalCoordinates[ch] = coordinatesLookup{cor: struct{}{}}
				}
			}

		default:
			lvlInfo[mode] = msg
		}

		row++
	}
}

func printThings() {
	communication.Log("\n", "boxColor")
	for k, v := range boxColor {
		communication.Log(string(k), v)
	}
	communication.Log("\n", "agentColor")
	for k, v := range agentColor {
		communication.Log(string(k), v)
	}
	communication.Log("\n", "agentCoordinates")
	for k, v := range agentCoordinates {
		communication.Log(string(k), v)
	}
	communication.Log("\n", "boxCoordinates")
	for k, v := range boxCoordinates {
		communication.Log(string(k), v)
	}
	communication.Log("\n", "walls")
	for k, v := range walls {
		communication.Log(k, v)
	}

	communication.Log("\n", "goals")
	for k, v := range goalCoordinates {
		communication.Log(string(k), v)
	}
}
