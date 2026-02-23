package plan

// LoopMode defines the operating mode of the Ralph loop.
type LoopMode string

const (
	ModePlanning LoopMode = "planning"
	ModeBuilding LoopMode = "building"
)
