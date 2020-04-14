package types

type SimpleMap map[byte]string

// We might be able to use uint16 here
type Point uint32

type Coordinates [2]Point

type CoordinatesLookup map[Coordinates]struct{}

type IntrestingCoordinates map[byte]CoordinatesLookup
type LevelInfo struct {
	LevelInfo        map[string]string
	AgentColor       SimpleMap
	BoxColor         SimpleMap
	WallsCoordinates map[Coordinates]struct{}
	AgentCoordinates IntrestingCoordinates
	BoxCoordinates   IntrestingCoordinates
	GoalCoordinates  IntrestingCoordinates
}

func GetLevelInfo() LevelInfo {
	return LevelInfo{
		LevelInfo:        make(map[string]string, 2),
		AgentColor:       SimpleMap{},
		BoxColor:         SimpleMap{},
		WallsCoordinates: make(map[Coordinates]struct{}, 15),
		GoalCoordinates:  IntrestingCoordinates{},
		AgentCoordinates: IntrestingCoordinates{},
		BoxCoordinates:   IntrestingCoordinates{},
	}
}
