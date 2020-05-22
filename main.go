package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/alexdor/dtu-ai-mas-final-assignment/ai"
	"github.com/alexdor/dtu-ai-mas-final-assignment/communication"
	"github.com/alexdor/dtu-ai-mas-final-assignment/parser"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	isDebug := len(os.Getenv("DEBUG")) > 0
	if isDebug {
		time.Sleep(10 * time.Second)
		communication.Log("Starting search")
	}

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("Got a timeout")

		if *cpuprofile != "" {
			pprof.StopCPUProfile()
		}

		os.Exit(1)
	}()

	communication.Init()

	levelInfo, currentState, err := parser.ParseLevel()

	if err != nil {
		communication.Error(err)
		return
	}

	ai.Play(&levelInfo, &currentState, &ai.AStart{}, isDebug)
}
