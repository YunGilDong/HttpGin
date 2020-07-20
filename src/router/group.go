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
	fmt.Println("/Group")

	//var mData map[int]data.Group
	mData := make(map[int]*data.Group)
	ok1, mData2 := mariadb.Mdb.GetGroup2(mData)

	fmt.Println(mData2)

	if ok1 {
		for k, v := range mData2 {
			fmt.Println(k, v)
		}
	}

	//GetGroupState
	var sData []data.GrpState
	ok, sData := mariadb.Mdb.GetGroupState(sData)
	if ok {
		for idx := 0; idx < len(sData); idx++ {

			key := sData[idx].GRP_ID
			grp := mData2[key]
			grp.State = sData[idx]
		}
	}

	//GetGroupPlan
	var sData2 []data.GrpOprState
	ok3, sData2 := mariadb.Mdb.GetGroupOprState(sData2)
	if ok3 {
		for idx := 0; idx < len(sData2); idx++ {

			key := sData2[idx].GRP_ID
			grp := mData2[key]
			grp.OprState = sData2[idx]
		}
	}

	jsonBytes, err := json.Marshal(mData2)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	jsonString := string(jsonBytes)
	fmt.Println(jsonString)

	json.NewEncoder(w).Encode(mData2)
}
