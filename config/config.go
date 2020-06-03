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

	BytesUsedForEachPoint      = 4                         // Each point is turned into a uint32 which consumes 4 bytes
	BytesUsedForEachAgentOrBox = BytesUsedForEachPoint * 2 // We have 2 points
)
