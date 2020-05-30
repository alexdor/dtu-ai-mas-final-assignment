package config

import "os"

type config struct {
	Name string
}

var (
	Config = config{
		Name: "NeverAI",
	}

	_, IsDebug = os.LookupEnv("DEBUG")
)

const (
	WallsSymbol      = '+'
	FreeSpaceSymbol  = ' '
	ServersTrueValue = "true"
)
