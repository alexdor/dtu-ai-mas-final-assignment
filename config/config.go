package config

import "os"

type config struct {
	Name string
}

var (
	Config = config{
		Name: "NeverAI",
	}

	IsDebug = len(os.Getenv("DEBUG")) > 0
)

const (
	WallsSymbol      = '+'
	FreeSpaceSymbol  = ' '
	ServersTrueValue = "true"
)
