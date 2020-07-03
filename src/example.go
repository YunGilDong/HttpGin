package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"router" //custom pkg
	"time"

	"github.com/gin-gonic/gin"
)

type logFileManage struct {
	fp  *os.File
	yy  int
	mm  int
	dd  int
	min int
}

var logMng logFileManage

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

func main() {

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
