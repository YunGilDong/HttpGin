package trffic_obj

import (
	"data"
	"log"
	"mariadb"
	"sync"
)

type LCobjects struct {
	MapLc map[int]*data.Loc // Lc map obejct
}

var lcobjs LCobjects
var lcMutex *sync.Mutex
var QuechLcState chan data.Loc

func InitLcObjects() {

	log.Println("InitLcObjects")
	lcobjs = LCobjects{}
	lcobjs.MapLc = make(map[int]*data.Loc)

	_, lcData := mariadb.Mdb.GetLocal()
	for idx := 0; idx < len(lcData); idx++ {
		key := lcData[idx].LOC_ID
		lcobjs.MapLc[key] = &lcData[idx]
	}

	QuechLcState = make(chan data.Loc, 5)

	// create Lc Mutex
	//lcMutex = &sync.Mutex()
}

func GetLcObjectsValue() LCobjects {
	obj := lcobjs
	obj.MapLc = make(map[int]*data.Loc)

	for k, v := range lcobjs.MapLc {
		var cpVal data.Loc = data.Loc{}
		cpVal = *v
		obj.MapLc[k] = &cpVal
	}

	return obj
}

func SetLcObjecState(lcdata data.Loc) {
	// 추후 뮤텍스 lock 필요할 수 있음
	log.Println(lcdata)
	id := lcdata.LOC_ID
	pLc := lcobjs.MapLc[id]

	var ok bool = false

	if pLc.State.CommSt != lcdata.State.CommSt {
		ok = true
	}
	if pLc.State.ConflictSt != lcdata.State.ConflictSt {
		ok = true
	}
	if pLc.State.DoorSt != lcdata.State.DoorSt {
		ok = true
	}
	if pLc.State.FlashSt != lcdata.State.FlashSt {
		ok = true
	}
	if pLc.State.LightOffSt != lcdata.State.LightOffSt {
		ok = true
	}
	if pLc.State.OprMode != lcdata.State.OprMode {
		ok = true
	}

	if ok {
		QuechLcState <- lcdata
	}

	pLc.State = lcdata.State
}
