package global

import "genLib"

//var DBlog genLib.InitOLog("../log", "MARIDB")

var DBlog *genLib.OLog = genLib.InitOLog("../log", "MARIDB")
