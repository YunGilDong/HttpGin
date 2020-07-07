package global

import (
	"genLib"
)

//var DBlog genLib.InitOLog("../log", "MARIDB")

var DBlog *genLib.OLog = genLib.InitOLog("../log", "MARIDB")
var Tcplog *genLib.OLog = genLib.InitOLog("../log", "NETWORK")

//var Mdb *mariadb.MariaDB = mariadb.InitDBSrc("dev", "dev", "sbrt_test", "192.168.1.74")
