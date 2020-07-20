package trffic_obj

import (
	"data"
)

var existLcSummary data.RLcStateSummary

func initLcSummaray() {
	existLcSummary = data.RLcStateSummary{}

	existLcSummary.LcCount = 0
	existLcSummary.ScuFixCyc = 0
	existLcSummary.LocalAct = 0
	existLcSummary.LocalNonAct = 0
	existLcSummary.CenterAct = 0
	existLcSummary.CenterNonAct = 0
	existLcSummary.KeepPhase = 0
	existLcSummary.ConflictImpossible = 0
	existLcSummary.CommError = 0
	existLcSummary.LightOff = 0
	existLcSummary.Flash = 0
	existLcSummary.DoorOpen = 0
	existLcSummary.Conflict = 0
}

func GetLcStateSummary() (data.RLcStateSummary, bool) {
	lcStats := GetLcObjectsValue() // 두번쨰 파라미터에 변화 여부 추가??

	lcSummary := data.RLcStateSummary{}

	lcSummary.LcCount = len(lcStats.MapLc)
	lcSummary.ScuFixCyc = 0
	lcSummary.LocalAct = 0
	lcSummary.LocalNonAct = 0
	lcSummary.CenterAct = 0
	lcSummary.CenterNonAct = 0
	lcSummary.KeepPhase = 0
	lcSummary.ConflictImpossible = 0

	var comm, light, flash, door, conflict int = 0, 0, 0, 0, 0

	for k, _ := range lcStats.MapLc {
		lcObj := lcStats.MapLc[k]

		if lcObj.State.CommSt == 1 {
			comm++
		}

		if lcObj.State.LightOffSt == 1 {
			light++
		}

		if lcObj.State.FlashSt == 1 {
			flash++
		}

		if lcObj.State.DoorSt == 1 {
			door++
		}

		if lcObj.State.ConflictSt == 1 {
			conflict++
		}
	}

	lcSummary.CommError = comm
	lcSummary.LightOff = light
	lcSummary.Flash = flash
	lcSummary.DoorOpen = door
	lcSummary.Conflict = conflict

	var isChanged = false
	if existLcSummary.CommError != comm {
		isChanged = true
	} else if existLcSummary.LightOff != light {
		isChanged = true
	} else if existLcSummary.Flash != flash {
		isChanged = true
	} else if existLcSummary.DoorOpen != door {
		isChanged = true
	} else if existLcSummary.Conflict != conflict {
		isChanged = true
	}

	existLcSummary = lcSummary

	return lcSummary, isChanged
}
