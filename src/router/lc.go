package router

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
	"trffic_obj"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
)

type eventData struct {
	message string
	data    string
}

// get : /lc
func Lc(c *gin.Context) {
	println("Lc router!")

	fmt.Println("lc")

	c.JSON(200, gin.H{
		"lc": "lc ok",
	})
}

// get : /lc_state_summary
func Lc_state_summary(c *gin.Context) {
	println("Lc_state_summary router!")

	fmt.Println("Lc_state_summary")

	apiData, _ := trffic_obj.GetLcStateSummary()

	fmt.Println(apiData)
	c.JSON(200, apiData)
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

			evdata.message = "message"
			evdata.data = string(b)
			//log.Println(evdata.data)

			if isChanged {
				ch_eventdata <- evdata
				log.Println("changed")
			} else {
				log.Println("not changed")
			}
			time.Sleep(time.Second * 1)
		}
	}()

	// check event recv & event send
	c.Stream(func(w io.Writer) bool {
		if ev_msg, ok := <-ch_eventdata; ok {
			switch ev_msg.message {
			case "message":

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
