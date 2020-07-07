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

var Mdb *MariaDB = InitDBSrc("dev", "dev", "sbrt_test", "192.168.1.74")

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
								, IFNULL(AUTO_ONLINE,0) 
								, IFNULL(GRP_LAT,0) 
								, IFNULL(GRP_LON,0) 
								, IFNULL(GRP_ZOMMLV,0)  
							from INT_GRP`)

	if err != nil {
		fmt.Println("[GetGroup Query error(1)]", err.Error())
		global.DBlog.Write("[GetGroup Query error(1)]" + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var grp data.Group
		err := rows.Scan(&grp.GRP_ID, &grp.GRP_NM, &grp.GRP_DEFCTRLMODE, &grp.AUTO_ONLINE, &grp.GRP_LAT, &grp.GRP_LON, &grp.GRP_ZOMMLV)

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
