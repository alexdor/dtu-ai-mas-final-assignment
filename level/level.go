package level

type (
	SimpleMap map[byte]string

	Point = int

	// Coordinates {row,column}
	Coordinates [2]Point

	CoordinatesLookup map[Coordinates]struct{}

	IntrestingCoordinates map[byte][]Coordinates

	NodeOrAgent struct {
		Letter byte
		Coordinates
	}

	Visited map[ID]struct{}

	ID string

	Info struct {
		LevelInfo              map[string]string
		GoalCount              int
		AgentColor             SimpleMap
		BoxColor               SimpleMap
		MaxCoord               Coordinates
		WallsCoordinates       CoordinatesLookup
		InGameWallsCoordinates []Coordinates // Sorted based on row
		GoalCoordinates        IntrestingCoordinates
		BoxGoalAssignment      []Coordinates
	}
)

func (levelInfo Info) IsWall(coor Coordinates) bool {
	_, ok := levelInfo.WallsCoordinates[coor]
	return ok
}

func (levelInfo Info) IsCellFree(coor Coordinates, currentState *CurrentState) bool {
	if _, ok := levelInfo.WallsCoordinates[coor]; ok {
		return false
	}

	for _, v := range currentState.Boxes {
		if v.Coordinates == coor {
			return false
		}
	}
	for _, v := range currentState.Agents {
		if v.Coordinates == coor {
			return false
		}
	}

	return true
}

func (i *Info) Init() {
	i.LevelInfo = make(map[string]string, 2)
	i.AgentColor = SimpleMap{}
	i.BoxColor = SimpleMap{}
	i.WallsCoordinates = make(CoordinatesLookup, 15)
	i.GoalCoordinates = IntrestingCoordinates{}
}
