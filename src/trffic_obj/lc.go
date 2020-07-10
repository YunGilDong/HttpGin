package trffic_obj

import (
	"data"
	"log"
	"sync"
)

type LCobjects struct {
	mapLc map[int]data.LC // Lc map obejct
}

var lcobjs LCobjects
var lcMutex *sync.Mutex

func InitLcObjects() {

	log.Println("InitLcObjects")
	lcobjs = LCobjects{}
	lcobjs.mapLc = make(map[int]data.LC)

	// create Lc Mutex
	//lcMutex = &sync.Mutex()
}

func GetLcObjectsValue() LCobjects {
	obj := lcobjs
	obj.mapLc = make(map[int]data.LC)

	for k, v := range lcobjs.mapLc {
		obj.mapLc[k] = v
	}

	return obj
}

func SetLcObjecState(lcdata data.LC) {
	// 추후 뮤텍스 lock 필요할 수 있음
	log.Println(lcdata)
	id := lcdata.LC_ID
	//log.Println(id, lcdata)
	lcobjs.mapLc[id] = lcdata

	//log.Println(lcobjs.mapLc)
}
