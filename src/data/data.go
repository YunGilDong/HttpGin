package data

//---------------------------------------------------------------------------
// Lc State
//---------------------------------------------------------------------------
type LcState struct {
	oprMode    int
	conflictSt int
	lightOffSt int
	flashSt    int
	doorSt     int
	commSt     int
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
	state  LcState
}
