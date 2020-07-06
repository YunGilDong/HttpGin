package router

import (
	"data"
	"encoding/json"
	"fmt"
	"global"
	"mariadb"

	"github.com/gin-gonic/gin"
)

//var DBlog *genLib.OLog = genLib.InitOLog("../log", "MARIDB")

var mdb *mariadb.MariaDB = mariadb.InitDBSrc("dev", "dev", "sbrt_test", "192.168.1.74")

// post : /login
func Login(c *gin.Context) {
	println("Login router!")
	// c.JSON(200, gin.H{
	// 	"message": "ping",
	// })

	fmt.Println("Login")
	global.DBlog.Write("main", "Login")

	//InitDBSrc(user string, passwd string, dbName string, hostAddr string) {

	// mdb := mariadb.InitDBSrc("dev", "dev", "sbrt_test", "192.168.1.74")

	var sData []data.Group
	ok, sData := mdb.GetGroup(sData)
	if ok {
		for idx := 0; idx < len(sData); idx++ {
			fmt.Println("ID : ", sData[idx].GRP_ID, "NM : ", sData[idx].GRP_NM)
		}
	}

	jsonBytes, err := json.Marshal(sData)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	jsonString := string(jsonBytes)
	fmt.Println(jsonString)

	c.JSON(200, sData)
}
