package mariadb

import (
	"data"
	"database/sql"
	"fmt"
	"global"

	_ "github.com/go-sql-driver/mysql"
)

type MariaDB struct {
	user, passwd, dbName, hostAddr string
}

var Mdb *MariaDB = InitDBSrc("dev", "dev", "sbrt", "192.168.1.74")

func InitDBSrc(user string, passwd string, dbName string, hostAddr string) *MariaDB {
	fmt.Println("InitDBSrc")
	mdb := MariaDB{}
	mdb.user = user
	mdb.passwd = passwd
	mdb.hostAddr = hostAddr
	mdb.dbName = dbName

	return &mdb
}

func (mdb *MariaDB) GetGroup(sData []data.Group) (bool, []data.Group) {
	fmt.Println("GetGroup")

	dbSrc := mdb.user + ":" + mdb.passwd + "@tcp(" + mdb.hostAddr + ")/" + mdb.dbName

	// open
	db, err := sql.Open("mysql", dbSrc)

	global.DBlog.Write("DB", "GetGroup")

	if err != nil {
		fmt.Println(err.Error())
		global.DBlog.Write("[DB Open error(1)]" + err.Error())
		return false, sData
	}
	defer db.Close()

	// Query
	rows, err := db.Query(`SELECT GRP_ID 
								, GRP_NM 
								, IFNULL(GRP_DEFCTRLMODE,0)
								, IFNULL(AUTO_ONLINEYN,0)
								, IFNULL(GRP_LAT,0)
								, IFNULL(GRP_LON,0)
								, IFNULL(GRP_ZOMMLV,0)
							from GRP_MST`)

	if err != nil {
		fmt.Println("[GetGroup Query error(1)]", err.Error())
		global.DBlog.Write("[GetGroup Query error(1)]" + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var grp data.Group
		err := rows.Scan(&grp.GRP_ID, &grp.GRP_NM, &grp.GRP_DEFCTRLMODE, &grp.AUTO_ONLINEYN, &grp.GRP_LAT, &grp.GRP_LON, &grp.GRP_ZOMMLV)

		if err != nil {
			fmt.Println("[GetGroup Query error(2)]", err.Error())
			global.DBlog.Write("[GetGroup Query error(2)]" + err.Error())
			return false, sData
		}

		// Return value
		sData = append(sData, grp)
	}

	return true, sData
}

func (mdb *MariaDB) GetGroup2(sData map[int]*data.Group) (bool, map[int]*data.Group) {
	fmt.Println("GetGroup")

	dbSrc := mdb.user + ":" + mdb.passwd + "@tcp(" + mdb.hostAddr + ")/" + mdb.dbName

	// open
	db, err := sql.Open("mysql", dbSrc)

	global.DBlog.Write("DB", "GetGroup")

	if err != nil {
		fmt.Println(err.Error())
		global.DBlog.Write("[DB Open error(1)]" + err.Error())
		return false, sData
	}
	defer db.Close()

	// Query
	rows, err := db.Query(`SELECT GRP_ID 
								, GRP_NM 
								, IFNULL(GRP_DEFCTRLMODE,0)
								, IFNULL(AUTO_ONLINEYN,0)
								, IFNULL(GRP_LAT,0)
								, IFNULL(GRP_LON,0)
								, IFNULL(GRP_ZOMMLV,0)
							from GRP_MST`)

	if err != nil {
		fmt.Println("[GetGroup Query error(1)]", err.Error())
		global.DBlog.Write("[GetGroup Query error(1)]" + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var grp data.Group
		err := rows.Scan(&grp.GRP_ID, &grp.GRP_NM, &grp.GRP_DEFCTRLMODE, &grp.AUTO_ONLINEYN, &grp.GRP_LAT, &grp.GRP_LON, &grp.GRP_ZOMMLV)

		if err != nil {
			fmt.Println("[GetGroup Query error(2)]", err.Error())
			global.DBlog.Write("[GetGroup Query error(2)]" + err.Error())
			return false, sData
		}

		// Return value
		sData[grp.GRP_ID] = &grp
		//sData = append(sData, grp)
	}

	return true, sData
}

func (mdb *MariaDB) GetLocal(sData []data.Loc) (bool, []data.Loc) {
	fmt.Println("GetLocal")

	dbSrc := mdb.user + ":" + mdb.passwd + "@tcp(" + mdb.hostAddr + ")/" + mdb.dbName

	// open
	db, err := sql.Open("mysql", dbSrc)

	global.DBlog.Write("DB", "GetLocal")

	if err != nil {
		fmt.Println(err.Error())
		global.DBlog.Write("[DB Open error(1)]" + err.Error())
		return false, sData
	}
	defer db.Close()

	// Query
	rows, err := db.Query(`SELECT a.LOC_ID
								, a.LOC_NM
								, a.GRP_ID
								, IFNULL(b.NODE_ID, "")
								, IFNULL(b.NODELAT,"")
								, IFNULL(b.NODELON,"")
							FROM LOC_MST as a
							left outer join node as b
							on a.NODE_ID = b.NODE_ID`)

	if err != nil {
		fmt.Println("[GetLocal Query error(1)]", err.Error())
		global.DBlog.Write("[GetLocal Query error(1)]" + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var loc data.Loc
		err := rows.Scan(&loc.LOC_ID, &loc.LOC_NM, &loc.GRP_ID, &loc.NODE_ID, &loc.NODELAT, &loc.NODELON)

		if err != nil {
			fmt.Println("[GetLocal Query error(2)]", err.Error())
			global.DBlog.Write("[GetLocal Query error(2)]" + err.Error())
			return false, sData
		}

		// Return value
		sData = append(sData, loc)
	}

	return true, sData
}

func (mdb *MariaDB) GetGroupState(sData []data.GrpState) (bool, []data.GrpState) {
	fmt.Println("GetGroupState")

	dbSrc := mdb.user + ":" + mdb.passwd + "@tcp(" + mdb.hostAddr + ")/" + mdb.dbName

	// open
	db, err := sql.Open("mysql", dbSrc)

	global.DBlog.Write("DB", "GetGroupState")

	if err != nil {
		fmt.Println(err.Error())
		global.DBlog.Write("[DB Open error(1)]" + err.Error())
		return false, sData
	}
	defer db.Close()

	// Query
	rows, err := db.Query(`SELECT A.GRP_ID
								, B.CREDATE
								, B.GRP_CTRLMODE
								, B.GRP_CTRLSTATE
							FROM grp_mst AS A
							LEFT OUTER JOIN grp_state AS B
							ON A.GRP_ID = B.GRP_ID;`)

	if err != nil {
		fmt.Println("[GetGroupState Query error(1)]", err.Error())
		global.DBlog.Write("[GetGroupState Query error(1)]" + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var grpSts data.GrpState
		err := rows.Scan(&grpSts.GRP_ID, &grpSts.CREDATE, &grpSts.GRP_CTRLMODE, &grpSts.GRP_CTRLSTATE)

		if err != nil {
			fmt.Println("[GetGroupState Query error(2)]", err.Error())
			global.DBlog.Write("[GetGroupState Query error(2)]" + err.Error())
			return false, sData
		}

		// Return value
		sData = append(sData, grpSts)
	}

	return true, sData
}

func (mdb *MariaDB) GetGroupOprState(sData []data.GrpOprState) (bool, []data.GrpOprState) {
	fmt.Println("GetGroupState")

	dbSrc := mdb.user + ":" + mdb.passwd + "@tcp(" + mdb.hostAddr + ")/" + mdb.dbName

	// open
	db, err := sql.Open("mysql", dbSrc)

	global.DBlog.Write("DB", "GetGroupOprState")

	if err != nil {
		fmt.Println(err.Error())
		global.DBlog.Write("[DB Open error(1)]" + err.Error())
		return false, sData
	}
	defer db.Close()

	// Query
	rows, err := db.Query(`SELECT A.GRP_ID
								, B.CREDATE
								, B.GRP_CTRLMODE
								, B.GRP_CTRLSTATE
								, B.NOW_GRP_CYCLELEN
								, B.NOW_LOC_TIMEPLANNUM
								, B.NOW_LOC_OFFSETPLANIDX
								, B.NOW_LOC_PHASEPLANIDX
							FROM grp_mst AS A
							LEFT OUTER JOIN grp_oprstate AS B
							ON A.GRP_ID = B.GRP_ID;`)

	if err != nil {
		fmt.Println("[GetGroupOprState Query error(1)]", err.Error())
		global.DBlog.Write("[GetGroupOprState Query error(1)]" + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var grpOprSts data.GrpOprState
		err := rows.Scan(&grpOprSts.GRP_ID, &grpOprSts.CREDATE, &grpOprSts.GRP_CTRLMODE, &grpOprSts.GRP_CTRLSTATE, &grpOprSts.NOW_GRP_CYCLELEN, &grpOprSts.NOW_LOC_TIMEPLANNUM, &grpOprSts.NOW_LOC_OFFSETPLANIDX, &grpOprSts.NOW_LOC_PHASEPLANIDX)

		if err != nil {
			fmt.Println("[GetGroupOprState Query error(2)]", err.Error())
			global.DBlog.Write("[GetGroupOprState Query error(2)]" + err.Error())
			return false, sData
		}

		// Return value
		sData = append(sData, grpOprSts)
	}

	return true, sData
}
