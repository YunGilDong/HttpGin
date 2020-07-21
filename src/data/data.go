package data

//---------------------------------------------------------------------------
// Lc State
//---------------------------------------------------------------------------
type LcState struct {
	OprMode    int
	ConflictSt int
	LightOffSt int
	FlashSt    int
	DoorSt     int
	CommSt     int
}

type GrpOprState struct {
	GRP_ID                int
	CREDATE               string
	GRP_CTRLMODE          int
	GRP_CTRLSTATE         int
	NOW_GRP_CYCLELEN      int
	NOW_LOC_TIMEPLANNUM   int
	NOW_LOC_OFFSETPLANIDX int
	NOW_LOC_PHASEPLANIDX  int
}

type GrpState struct {
	GRP_ID        int
	CREDATE       string
	GRP_CTRLMODE  int
	GRP_CTRLSTATE int
	TRC_POSSIYN   int
	SWCH_REQMODE  int
}

//---------------------------------------------------------------------------
// LOC (Local)
//---------------------------------------------------------------------------
type Loc struct {
	LOC_ID  int
	LOC_NM  string
	GRP_ID  int
	NODE_ID string
	NODELAT string
	NODELON string
	State   LcState
}

//---------------------------------------------------------------------------
// Group
//---------------------------------------------------------------------------
type Group struct {
	GRP_ID          int
	GRP_NM          string
	GRP_DEFCTRLMODE int
	AUTO_ONLINEYN   int
	GRP_LAT         string
	GRP_LON         string
	GRP_ZOMMLV      int
	State           GrpState
	OprState        GrpOprState
}
