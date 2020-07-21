package router

import (
	"data"
	"encoding/json"
	"fmt"
	"mariadb"
	"net/http"
	"trffic_obj"
)

func Loading(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("/loading")

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

	// local
	ok4, lcData := mariadb.Mdb.GetLocal()
	if ok4 {
		for idx := 0; idx < len(lcData); idx++ {
			fmt.Println("ID : ", lcData[idx].LOC_ID, "NM : ", lcData[idx].LOC_NM)
		}
	}

	//var lcObj trffic_obj.LCobjects
	lcObj := trffic_obj.GetLcObjectsValue()

	for k, v := range lcObj.MapLc {
		fmt.Println("@@>>> : ", k, v)
		lcData[k-1].State = v.State // 대입.. 임시방편
	}

	// API
	var rGroup []data.RGroup

	for _, v := range mData2 {
		rGrp := trffic_obj.GetRGroup(*v)
		rGroup = append(rGroup, rGrp)
	}

	fmt.Println("lcData len : ", len(lcData))

	for idx := 0; idx < len(lcData); idx++ {
		rLc := trffic_obj.GetRLoc(lcData[idx])
		id := rLc.GrpId - 1
		rGroup[id].Locs = append(rGroup[id].Locs, rLc)
	}

	jsonBytes, err := json.Marshal(rGroup)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	jsonString := string(jsonBytes)
	fmt.Println(jsonString)

	json.NewEncoder(w).Encode(rGroup)

}
