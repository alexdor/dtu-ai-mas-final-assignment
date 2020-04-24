package main

import (
	"github.com/alexdor/dtu-ai-mas-final-assignment/ai"
	"github.com/alexdor/dtu-ai-mas-final-assignment/communication"
	"github.com/alexdor/dtu-ai-mas-final-assignment/parser"
)

func main() {
	communication.Init()

	levelInfo, err := parser.ParseLevel()
	if err != nil {
		communication.Error(err)
		return
	}

	ai.Play(&levelInfo, ai.BiDirectionalBFS{}, ai.ManhattanDistance{})
}
