package router

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mariadb"
	"net/http"
	"time"
	"trffic_obj"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
)

type eventData struct {
	message string
	data    string
}

// get : /local
func Local(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("/Local")

	ok, sData := mariadb.Mdb.GetLocal()
	if ok {
		for idx := 0; idx < len(sData); idx++ {
			fmt.Println("ID : ", sData[idx].LOC_ID, "NM : ", sData[idx].LOC_NM)
		}
	}

	//var lcObj trffic_obj.LCobjects
	lcObj := trffic_obj.GetLcObjectsValue()

	for k, v := range lcObj.MapLc {
		fmt.Println("@@>>> : ", k, v)
		sData[k-1].State = v.State // 대입.. 임시방편
	}

	jsonBytes, err := json.Marshal(sData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	jsonString := string(jsonBytes)
	fmt.Println(jsonString)

	json.NewEncoder(w).Encode(sData)
}

// get : /lc_state_summary
func Lc_state_summary(w http.ResponseWriter, r *http.Request) {
	println("Lc_state_summary router!")

	fmt.Println("Lc_state_summary")

	apiData, _ := trffic_obj.GetLcStateSummary()

	jsonBytes, err := json.Marshal(apiData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	jsonString := string(jsonBytes)
	fmt.Println(jsonString)

	fmt.Println(apiData)
	json.NewEncoder(w).Encode(apiData)
}

// get : /lc_event
func Lc_event(c *gin.Context) {
	ch_eventdata := make(chan eventData)

	// event 변화 체크 (cycle : 1sec)
	go func() {
		defer close(ch_eventdata)

		for {
			var evdata eventData
			// check lc_state_summary
			apiData, isChanged := trffic_obj.GetLcStateSummary()
			b, _ := json.Marshal(apiData)
			// log.Println(b)         // [123 34 78 97 109 101 34 58 34 ...]
			// log.Println(string(b)) // {"Name":"Gopher","Age":7}

			evdata.message = "lcstatus"
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
	}()

	// check event recv & event send
	c.Stream(func(w io.Writer) bool {
		if ev_msg, ok := <-ch_eventdata; ok {
			switch ev_msg.message {
			case "lcstatus":

				c.Render(-1, sse.Event{
					Event: ev_msg.message,
					Data:  ev_msg.data,
				})

			default:
				return (false)
			}

			return (true)
		}
		return (false)
	})
}
