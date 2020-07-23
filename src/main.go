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
	trffic_obj.InitGrpObject()
}

func initDatabase() {
	// group database set
	mData := make(map[int]*data.Group)
	ok := mariadb.Mdb.GetGroup2(mData)

	fmt.Println(mData)

	if ok {
		for _, v := range mData {
			trffic_obj.SetGrpObject(*v)
		}
	}
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
	http.HandleFunc("/loading", router.Loading)
	http.HandleFunc("/login", handler)

}

func eventHandler(ch_eventdata chan eventmessage) {
	for {
		var evdata eventmessage
		// check lc_state_summary
		apiData, isChanged := trffic_obj.GetLcStateSummary()
		b, _ := json.Marshal(apiData)

		evdata.messagetype = "lcstatus"
		evdata.data = string(b)

		// Total 현황 event
		if isChanged {
			ch_eventdata <- evdata
			log.Println("Total 현황 event changed")
		} else {
			//log.Println("not changed")
		}

		// Local status event

		// Group status event
		// groups := trffic_obj.GetGrpObjectValue()
		// for k, _ := range groups.MapGrp {
		// 	isChanged, groupObj := trffic_obj.CheckGrpStatus(k)
		// 	if isChanged {
		// 		apiGroupSts := trffic_obj.GetGroupStatus(groupObj)
		// 		b, _ := json.Marshal(apiGroupSts)
		// 		evdata.messagetype = "groupstatus"
		// 		evdata.data = string(b)
		// 		//log.Println(evdata.data)
		// 		ch_eventdata <- evdata
		// 		log.Println("group status changed")
		// 	}
		// }

		// Local status event
		select {
		case lcState := <-trffic_obj.QuechLcState:
			rLocdata := trffic_obj.GetRLoc(lcState)
			fmt.Println("lc statue event", rLocdata)
			b, _ := json.Marshal(rLocdata)

			evdata.messagetype = "lcstatusev"
			evdata.data = string(b)

			log.Println("lc status changed")

			ch_eventdata <- evdata
		default:

		}

		time.Sleep(time.Millisecond * 10)
	}
}

func evnetSend(es eventsource.EventSource) {
	ch_eventdata := make(chan eventmessage)
	go eventHandler(ch_eventdata)

	for data := range ch_eventdata {
		log.Println("event msg : ", data.messagetype)
		switch data.messagetype {
		case "lcstatus":
			log.Println("lcstatus event", data.data)
			es.SendEventMessage(data.data, "lcstatus", "1")
		case "groupstatus":
			log.Println("groupstatus event", data.data)
			es.SendEventMessage(data.data, "groupstatus", "2")
		case "lcstatusev":
			log.Println("lcstatusev event", data.data)
			es.SendEventMessage(data.data, "lcstatusev", "3")
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
	initDatabase()
	initRoutine()
	initSignal()

	// event
	var es eventsource.EventSource

	//es = eventsource.New(nil, nil)
	es = eventsource.New(
		eventsource.DefaultSettings(),
		func(req *http.Request) [][]byte {
			return [][]byte{
				[]byte("X-Accel-Buffering: no"),
				[]byte("Access-Control-Allow-Origin: *"),
				[]byte("UTF-8"),
			}
		},
	)
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
