package router

import (
	"data"
	"encoding/json"
	"fmt"
	"global"
	"mariadb"

	"github.com/gin-gonic/gin"
)

// get : /index
func Index(c *gin.Context) {
	println("Login router!")
	// c.JSON(200, gin.H{
	// 	"message": "ping",
	// })

	fmt.Println("Login")
	global.DBlog.Write("main", "Login")

	var sData []data.Group
	ok, sData := mariadb.Mdb.GetGroup(sData)
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
