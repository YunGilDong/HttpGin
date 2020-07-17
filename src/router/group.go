package router

import (
	"data"
	"encoding/json"
	"fmt"
	"mariadb"

	"net/http"
)

// get : /group
func Group(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("/group")

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

	json.NewEncoder(w).Encode(sData)

}
