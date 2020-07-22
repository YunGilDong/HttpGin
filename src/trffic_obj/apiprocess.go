package trffic_obj

import (
	"data"
	"fmt"
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
	lcSummary.ScuFixCyc = 1
	lcSummary.LocalAct = 1
	lcSummary.LocalNonAct = 2
	lcSummary.CenterAct = 2
	lcSummary.CenterNonAct = 2
	lcSummary.KeepPhase = 0
	lcSummary.ConflictImpossible = 0
	lcSummary.Trans = 1

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

func GetRGroup(grp data.Group) data.RGroup {
	var rGrp data.RGroup
	rGrp.GrpId = grp.GRP_ID
	rGrp.GrpNm = grp.GRP_NM
	// rGrp.GrpLat = grp.GRP_LAT
	// rGrp.GrpLon = grp.GRP_LON
	rGrp.Position.X = grp.GRP_LAT
	rGrp.Position.Y = grp.GRP_LON
	rGrp.GrpDefMode = grp.GRP_DEFCTRLMODE
	rGrp.Status.CreateTm = grp.State.CREDATE
	rGrp.Status.GrpCmode = grp.State.GRP_CTRLMODE
	rGrp.Status.GrpCstatus = grp.State.GRP_CTRLSTATE
	rGrp.Plan.Cycle = grp.OprState.NOW_GRP_CYCLELEN
	rGrp.Plan.Mode = grp.OprState.GRP_CTRLMODE
	rGrp.Plan.Offset = grp.OprState.NOW_LOC_OFFSETPLANIDX
	rGrp.Plan.Split = grp.OprState.NOW_LOC_PHASEPLANIDX

	return rGrp
}

func GetRLoc(locs data.Loc) data.RLoc {
	var rLoc data.RLoc
	fmt.Println(rLoc)

	rLoc.LocId = locs.LOC_ID
	rLoc.LocNm = locs.LOC_NM
	rLoc.Position.X = locs.NODELAT
	rLoc.Position.Y = locs.NODELON
	// rLoc.LocLat = locs.NODELAT
	// rLoc.LocLon = locs.NODELON
	rLoc.GrpId = locs.GRP_ID
	rLoc.Status.CommSt = locs.State.CommSt
	rLoc.Status.ConflictSt = locs.State.CommSt
	rLoc.Status.DoorSt = locs.State.DoorSt
	rLoc.Status.FlashSt = locs.State.FlashSt
	rLoc.Status.LightOffSt = locs.State.LightOffSt
	rLoc.Status.OprMode = locs.State.OprMode

	return rLoc
}

func GetGroupStatus(grp data.Group) data.RGroup {
	rGrp := GetRGroup(grp)
	rGrp.LocStatusCount = GetGroupLocStateSummary(rGrp.GrpId)
	rGrp.StrStatus = GetGroupOprStatusSummary(rGrp.GrpId)

	return rGrp
}

func GetGroupOprStatusSummary(groupId int) string {

	var strState string = ""
	if groupId == 1 {
		strState = "1/1"

	} else if groupId == 2 {
		strState = "3/3"

	}

	return strState

}

func GetGroupLocStateSummary(groupId int) string {
	var comm, light, flash, door, conflict int = 0, 0, 0, 0, 0

	lcStats := GetLcObjectsValue()

	for k, _ := range lcStats.MapLc {
		lcObj := lcStats.MapLc[k]

		if lcObj.GRP_ID == groupId {
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
	}

	var strStatus string = ""

	strStatus = fmt.Sprintf("%d/%d/%d/%d/%d", comm, light, flash, door, conflict)
	return strStatus
}
