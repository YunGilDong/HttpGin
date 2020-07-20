package trffic_obj

import (
	"data"
	"log"
	"sync"
)

type LCobjects struct {
	MapLc map[int]data.Loc // Lc map obejct
}

var lcobjs LCobjects
var lcMutex *sync.Mutex

func InitLcObjects() {

	log.Println("InitLcObjects")
	lcobjs = LCobjects{}
	lcobjs.MapLc = make(map[int]data.Loc)

	// create Lc Mutex
	//lcMutex = &sync.Mutex()
}

func GetLcObjectsValue() LCobjects {
	obj := lcobjs
	obj.MapLc = make(map[int]data.Loc)

	for k, v := range lcobjs.MapLc {
		obj.MapLc[k] = v
	}

	return obj
}

func SetLcObjecState(lcdata data.Loc) {
	// 추후 뮤텍스 lock 필요할 수 있음
	log.Println(lcdata)
	id := lcdata.LOC_ID
	//log.Println(id, lcdata)
	lcobjs.MapLc[id] = lcdata

	//log.Println(lcobjs.MapLc)
}
