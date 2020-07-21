package trffic_obj

import (
	"data"
	"fmt"
	"log"
	"mariadb"
)

type Grpobjects struct {
	MapGrp map[int]*data.Group // Group map obejct
}

var grpObj Grpobjects

func InitGrpObject() {
	log.Println("InitGrpObject")
	grpObj = Grpobjects{}
	grpObj.MapGrp = make(map[int]*data.Group)

	mariadb.Mdb.GetGroup2(grpObj.MapGrp)
	fmt.Println(grpObj)
}

func GetGrpObjectValue() Grpobjects {
	obj := grpObj
	obj.MapGrp = make(map[int]*data.Group)

	for k, v := range grpObj.MapGrp {
		obj.MapGrp[k] = v
	}

	return obj
}

func SetGrpObject(grpdata data.Group) {
	log.Println(grpdata)

	id := grpdata.GRP_ID
	grpObj.MapGrp[id] = &grpdata
}

func SetGrpStatus(state data.GrpState) {
	id := state.GRP_ID
	pGrp := grpObj.MapGrp[id]
	pGrp.State = state
}

func SetGrpOprStatus(state data.GrpOprState) {
	id := state.GRP_ID
	pGrp := grpObj.MapGrp[id]
	pGrp.OprState = state
}

func CheckGrpStatus(grpId int) (bool, data.Group) {

	var group data.Group = data.Group{}
	var grpSts data.GrpState = data.GrpState{}
	var grpOPrSts data.GrpOprState = data.GrpOprState{}

	var isChaned bool = false

	ok := mariadb.Mdb.GetGroupStateOprState(grpId, &grpSts, &grpOPrSts)

	if ok {
		pObj := grpObj.MapGrp[grpId]
		if pObj.State.CREDATE != grpSts.CREDATE || pObj.OprState.CREDATE != grpOPrSts.CREDATE {
			pObj.State = grpSts
			pObj.OprState = grpOPrSts

			group = *pObj
			isChaned = true

		}
	}

	return isChaned, group
}
