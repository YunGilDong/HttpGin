package data

// /api1/lc_state_summary
type LcStateSummary struct {
	LcCount            int
	ScuFixCyc          int
	LocalAct           int
	LocalNonAct        int
	CenterAct          int
	CenterNonAct       int
	KeepPhase          int
	CommError          int
	LightOff           int
	Flash              int
	Manual             int
	Conflict           int
	DoorOpen           int
	ConflictImpossible int
}
