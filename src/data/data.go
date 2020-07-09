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

//---------------------------------------------------------------------------
// Group
//---------------------------------------------------------------------------
type Group struct {
	GRP_ID          int
	GRP_NM          string
	GRP_DEFCTRLMODE int
	AUTO_ONLINE     int
	GRP_LAT         string
	GRP_LON         string
	GRP_ZOMMLV      int
}

//---------------------------------------------------------------------------
// LC
//---------------------------------------------------------------------------
type LC struct {
	LC_ID  int
	LC_NM  string
	GRP_ID int
	State  LcState
}
