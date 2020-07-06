package main

import (
	"data"
	"encoding/json"
	"fmt"
	"genLib"
	"global"
	"io"
	"log"
	"mariadb"
	"os"
	"router"
	"time"

	"github.com/gin-gonic/gin"
)

//------------------------------------------------------------------------------
// struct
//------------------------------------------------------------------------------
type logFileManage struct {
	fp  *os.File
	yy  int
	mm  int
	dd  int
	min int
}

//------------------------------------------------------------------------------
// Local
//------------------------------------------------------------------------------
var logMng logFileManage
var mlog genLib.OLog
var mdb mariadb.MariaDB

//------------------------------------------------------------------------------
// Global
//------------------------------------------------------------------------------
//var DBlog genLib.OLog

//------------------------------------------------------------------------------

func LoggerFileCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		//println("LoggerFileCheck!")

		c.Next()
	}
}

func initRouter() *gin.Engine {
	r := gin.Default()

	r1 := r.Group("/api1")
	r1.Use(LoggerFileCheck(), gin.Logger())
	{
		fmt.Println("login(1)")
		r1.GET("/login", router.Login) // GET => /api1/login
	}

	return r
}

// none
func initLogDirectory() {
	year, month, day := time.Now().Date()
	min := time.Now().Minute()

	logMng.yy = year
	logMng.mm = int(month)
	logMng.dd = day
	logMng.min = min

	dirname := fmt.Sprintf("%02d%02d%02d%02d", year, int(month), day, min)
	dirPath := "../log/" + dirname

	os.Mkdir("../log", os.ModePerm)
	os.Mkdir(dirPath, os.ModePerm)
}

func checkLogFile() {

	// http log file
	f, err := os.OpenFile("./gin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func initVariable() {
	logMng.yy = 0
	logMng.mm = 0
	logMng.dd = 0
}

func DbGetGroup() {

	global.DBlog.Write("main", "DbGetGroup")

	//InitDBSrc(user string, passwd string, dbName string, hostAddr string) {

	mdb := mariadb.InitDBSrc("dev", "dev", "sbrt_test", "192.168.1.74")

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

}

func main() {
	mlog := genLib.InitOLog("../log", "MAIN")

	mlog.Write("test", "hello1")

	var bytes []byte
	bytes = []byte{1, 2, 3}
	mlog.Dump("TEST", bytes)

	DbGetGroup()
	//initVariable()
	checkLogFile()
	r := initRouter()

	defer func() {
		println("main exit!")
	}()

	err := r.Run(":5000")
	if err != nil {
		println("Http run fail.. port 5000,,,")
	}
}