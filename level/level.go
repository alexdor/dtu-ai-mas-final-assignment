package level

type SimpleMap map[byte]string

// We might be able to use uint16 here
type Point uint32

type Coordinates [2]Point

type CoordinatesLookup map[Coordinates]struct{}

type IntrestingCoordinates map[byte]CoordinatesLookup
type Info struct {
	LevelInfo        map[string]string
	AgentColor       SimpleMap
	BoxColor         SimpleMap
	WallsCoordinates CoordinatesLookup
	GoalCoordinates  IntrestingCoordinates
	// TODO: consider using a simple list
	AgentCoordinates IntrestingCoordinates
	BoxCoordinates   IntrestingCoordinates
}

func (levelInfo Info) IsWall(coor Coordinates) bool {
	_, ok := levelInfo.WallsCoordinates[coor]
	return ok
}

func (levelInfo Info) sBoxFree(coor Coordinates) bool {
	if _, ok := levelInfo.WallsCoordinates[coor]; ok {
		return true
	}

	for _, coord := range levelInfo.AgentCoordinates {
		if _, ok := coord[coor]; ok {
			return true
		}
	}

	for _, coord := range levelInfo.BoxCoordinates {
		if _, ok := coord[coor]; ok {
			return true
		}
	}

	return false
}

func GetLevelInfo() Info {
	return Info{
		LevelInfo:        make(map[string]string, 2),
		AgentColor:       SimpleMap{},
		BoxColor:         SimpleMap{},
		WallsCoordinates: make(CoordinatesLookup, 15),
		GoalCoordinates:  IntrestingCoordinates{},
		AgentCoordinates: IntrestingCoordinates{},
		BoxCoordinates:   IntrestingCoordinates{},
	}
}
