package trffic_obj

import (
	"data"
)

func GetLcStateSummary() data.LcStateSummary {
	lcStats := GetLcObjectsValue()

	lcSummary := data.LcStateSummary{}

	lcSummary.LcCount = len(lcStats.mapLc)
	lcSummary.ScuFixCyc = 0
	lcSummary.LocalAct = 0
	lcSummary.LocalNonAct = 0
	lcSummary.CenterAct = 0
	lcSummary.CenterNonAct = 0
	lcSummary.KeepPhase = 0
	lcSummary.ConflictImpossible = 0

	var comm, light, flash, door, conflict int = 0, 0, 0, 0, 0

	for k, _ := range lcStats.mapLc {
		lcObj := lcStats.mapLc[k]

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

	return lcSummary
}
