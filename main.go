package main

import (
	"github.com/alexdor/dtu-ai-mas-final-assignment/communication"
	"github.com/alexdor/dtu-ai-mas-final-assignment/parser"
	"github.com/alexdor/dtu-ai-mas-final-assignment/types"
)

func main() {
	communication.Init()
	levelInfo, err := parser.ParseLevel()
	if err != nil {
		communication.Error(err)
	}
	printThings(levelInfo)
}

func printThings(levelInfo types.LevelInfo) {
	communication.Log("\n", "boxColor")
	for k, v := range levelInfo.BoxColor {
		communication.Log(string(k), v)
	}
	communication.Log("\n", "agentColor")
	for k, v := range levelInfo.AgentColor {
		communication.Log(string(k), v)
	}
	communication.Log("\n", "agentCoordinates")
	for k, v := range levelInfo.AgentCoordinates {
		communication.Log(string(k), v)
	}
	communication.Log("\n", "boxCoordinates")
	for k, v := range levelInfo.BoxCoordinates {
		communication.Log(string(k), v)
	}
	communication.Log("\n", "walls")
	for k, v := range levelInfo.WallsCoordinates {
		communication.Log(k, v)
	}

	communication.Log("\n", "goals")
	for k, v := range levelInfo.GoalCoordinates {
		communication.Log(string(k), v)
	}
}
