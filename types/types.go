package types

type SimpleMap map[byte]string

type CoordinatesLookup map[[2]uint]struct{}

type IntrestingCoordinates map[byte]CoordinatesLookup
type LevelInfo struct {
	LevelInfo        map[string]string
	AgentColor       SimpleMap
	BoxColor         SimpleMap
	WallsCoordinates map[[2]uint]struct{}
	AgentCoordinates IntrestingCoordinates
	BoxCoordinates   IntrestingCoordinates
	GoalCoordinates  IntrestingCoordinates
}

func GetLevelInfo() LevelInfo {
	return LevelInfo{
		LevelInfo:        make(map[string]string, 2),
		AgentColor:       SimpleMap{},
		BoxColor:         SimpleMap{},
		WallsCoordinates: make(map[[2]uint]struct{}, 15),
		GoalCoordinates:  IntrestingCoordinates{},
		AgentCoordinates: IntrestingCoordinates{},
		BoxCoordinates:   IntrestingCoordinates{},
	}
}
