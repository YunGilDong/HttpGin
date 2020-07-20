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
	"net/http"
	"network"
	"os"
	"os/signal"
	"router"
	"syscall"
	"time"
	"trffic_obj"

	"github.com/gin-gonic/gin"
	"gopkg.in/antage/eventsource.v1"
)

type eventmessage struct {
	messagetype string
	data        string
}

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

var terminate bool

//------------------------------------------------------------------------------
// Local
//------------------------------------------------------------------------------
var logMng logFileManage
var mlog genLib.OLog

//var mdb mariadb.MariaDB

//------------------------------------------------------------------------------
// Global
//------------------------------------------------------------------------------
//var DBlog genLib.OLog

//------------------------------------------------------------------------------

//------------------------------------------------------------------------------
// sigHandler
//------------------------------------------------------------------------------
func sigHandler(chSig chan os.Signal) {
	fmt.Println("sigHandler")
	for {
		signal := <-chSig
		switch signal {
		case syscall.SIGHUP:
			fmt.Printf("SIGHUP(%d)\n", signal)
			terminate = true
			//panic(signal)
		case syscall.SIGINT:
			fmt.Printf("SIGINT(%d)\n", signal)
			terminate = true
			//panic(signal)
		case syscall.SIGTERM:
			fmt.Printf("SIGTERM(%d)\n", signal)
			terminate = true
			//panic(signal)
		default:
			fmt.Printf("Unknown signal(%d)\n", signal)
			terminate = true
			//panic(signal)
		}
	}
}

func LoggerFileCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		//println("LoggerFileCheck!")

		c.Next()
	}
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

	trffic_obj.InitLcObjects()
}

func initRoutine() {
	network.Routine()
}

func DbGetGroup() {

	global.DBlog.Write("main", "DbGetGroup")

	//InitDBSrc(user string, passwd string, dbName string, hostAddr string) {

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

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
}

func initRoute2() {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/group", router.Group)
	http.HandleFunc("/local", router.Local)
	http.HandleFunc("/login", handler)
}

func eventHandler(ch_eventdata chan eventmessage) {
	for {
		var evdata eventmessage
		// check lc_state_summary
		apiData, isChanged := trffic_obj.GetLcStateSummary()
		b, _ := json.Marshal(apiData)
		// log.Println(b)         // [123 34 78 97 109 101 34 58 34 ...]
		// log.Println(string(b)) // {"Name":"Gopher","Age":7}

		evdata.messagetype = "lcstatus"
		evdata.data = string(b)
		//log.Println(evdata.data)

		if isChanged {
			ch_eventdata <- evdata
			log.Println("changed")
		} else {
			//ch_eventdata <- evdata
			log.Println("not changed")
		}
		time.Sleep(time.Second * 1)
	}
}

func evnetSend(es eventsource.EventSource) {
	ch_eventdata := make(chan eventmessage)
	go eventHandler(ch_eventdata)

	for data := range ch_eventdata {
		switch data.messagetype {
		case "lcstatus":
			es.SendEventMessage(data.data, "lcstatus", "1")
		}
	}
}

//------------------------------------------------------------------------------
// initSignal
//------------------------------------------------------------------------------
func initSignal() {
	fmt.Println("iniSignal")
	// signal handler
	ch_signal := make(chan os.Signal, 1)
	signal.Notify(ch_signal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go sigHandler(ch_signal)
}

func main() {
	terminate = false
	mlog := genLib.InitOLog("../log", "MAIN")
	mlog.Write("main", "start")

	initVariable()
	initRoutine()
	initSignal()

	// event
	var es eventsource.EventSource
	es = eventsource.New(nil, nil)
	defer es.Close()

	http.Handle("/events", es)
	go evnetSend(es)
	initRoute2()
	fmt.Println(">>>>>>1")
	go func() {
		http.ListenAndServe(":5000", nil)
	}()

	fmt.Println(">>>>>>2")
	for {
		if terminate {
			break
		}

		time.Sleep(time.Second * 2)
	}

	fmt.Println("exit")

}
