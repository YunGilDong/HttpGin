package data

// /api1/lc_state_summary
type RLcStateSummary struct {
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

type RLocStatus struct {
	OprMode    int
	ConflictSt int
	LightOffSt int
	FlashSt    int
	DoorSt     int
	CommSt     int
}

type RLoc struct {
	LocId  int
	LocNm  string
	GrpId  int
	GrpOrd int
	LocLat string
	LocLon string
	Status RLocStatus
}

type RGrpPlan struct {
	Mode   int
	Cycle  int
	Offset int
	Split  int
}

type RGrpStatus struct {
	CreateTm   string
	GrpCmode   int
	GrpCstatus int
}

type RGroup struct {
	GrpId          int
	GrpNm          string
	GrpLat         string
	GrpLon         string
	GrpDefMode     int
	Status         RGrpStatus
	Plan           RGrpPlan
	LocStatusCount string
	Locs           []RLoc
}
