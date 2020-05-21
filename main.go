package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexdor/dtu-ai-mas-final-assignment/ai"
	"github.com/alexdor/dtu-ai-mas-final-assignment/communication"
	"github.com/alexdor/dtu-ai-mas-final-assignment/parser"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("Got a timeout")
		os.Exit(1)
	}()
	communication.Init()

	levelInfo, currentState, err := parser.ParseLevel()
	currentState.LevelInfo = &levelInfo

	if err != nil {
		communication.Error(err)
		return
	}

	ai.Play(&levelInfo, &currentState, &ai.AStart{})
}
