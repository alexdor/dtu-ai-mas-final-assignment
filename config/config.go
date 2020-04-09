package config

type config struct {
	Name string
}

var Config = config{
	Name: "NeverAI",
}

const (
	WallsSymbol     = '+'
	FreeSpaceSymbol = ' '
)
