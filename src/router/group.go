package router

import (
	"data"
	"encoding/json"
	"fmt"
	"mariadb"
	"trffic_obj"

	"net/http"
)

// get : /group
func Group(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("/Group")

	//var mData map[int]data.Group
	mData2 := make(map[int]*data.Group)
	ok1 := mariadb.Mdb.GetGroup2(mData2)

	fmt.Println(mData2)

	if ok1 {
		for k, v := range mData2 {
			fmt.Println(k, v)
		}
	}

	//GetGroupState
	ok, sData := mariadb.Mdb.GetGroupState()
	if ok {
		for idx := 0; idx < len(sData); idx++ {

			key := sData[idx].GRP_ID
			grp := mData2[key]
			grp.State = sData[idx]
		}
	}

	//GetGroupPlan
	ok3, sData2 := mariadb.Mdb.GetGroupOprState()
	if ok3 {
		for idx := 0; idx < len(sData2); idx++ {

			key := sData2[idx].GRP_ID
			grp := mData2[key]
			grp.OprState = sData2[idx]
		}
	}

	// GetLocal
	_, locData := mariadb.Mdb.GetLocal()
	//var lcObj trffic_obj.LCobjects
	lcObj := trffic_obj.GetLcObjectsValue()

	for k, v := range lcObj.MapLc {
		fmt.Println("@@>>> : ", k, v)
		locData[k-1].State = v.State // 대입.. 임시방편
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
