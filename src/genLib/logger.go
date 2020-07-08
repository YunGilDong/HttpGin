package genLib

import (
	"fmt"
	"log"
	"os"
	"time"
)

type OLog struct {
	name     string
	rootDir  string
	filePath string
	stamp    string
}

func InitOLog(rootDir string, name string) *OLog {
	lg := OLog{}
	lg.rootDir = rootDir
	lg.name = name

	return &lg
}

func (lg *OLog) SetRootPath(path string) {
	lg.rootDir = path
	// root 존재 확인

	// root dir 존재 확인
	if _, err := os.Stat(lg.rootDir); os.IsNotExist(err) {
		// root dir 생성
		os.Mkdir(lg.rootDir, os.ModePerm)
	}
}

func (lg *OLog) SetStamp(stamp string) {

}

func (lg OLog) logging(bytes []byte) {
	// file open
	path := lg.rootDir + "/" + lg.filePath + "/" + lg.name + "-" + lg.filePath + ".log"
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	f.Write(bytes)

	defer f.Close()
}

func (lg OLog) Write(str ...string) {
	var stamp, msg string
	if len(str) == 0 {
		return
	}
	if len(str) == 1 {
		stamp = "-"
		msg = str[0]
	}
	if len(str) == 2 {
		msg = str[1]
	}

	filePath := lg.genDateString()
	lg.filePath = filePath
	lg.stamp = stamp
	message := lg.genDateTimeString() + stamp
	message += "  " + msg + "\n"
	lg.logging([]byte(message))
}

func (lg OLog) Dump(stamp string, bytes []byte, length int) {
	//message := lg.genDateTimeString() + "  " + stamp + "  [" + string(length) + "]" + "\n"
	message := fmt.Sprintf("%s  %s  [%d]\n", lg.genDateTimeString(), stamp, length)

	//for idx := 0; idx < len(bytes); idx++ {
	for idx := 0; idx < length; idx++ {
		if idx%20 == 0 {
			if idx != 0 {
				message += message + "\n"
			}
			message += "	"
		}
		message += fmt.Sprintf(" %02X", bytes[idx])
	}

	message += "\n"

	filePath := lg.genDateString()

	lg.filePath = filePath
	lg.logging([]byte(message))
}

func (lg OLog) genDateString() string {
	// log 패키지에서 시간 찍음
	year, month, day := time.Now().Date()

	filePath := fmt.Sprintf("%02d%02d%02d", year, int(month), day)
	dirPath := lg.rootDir + "/" + filePath

	// log directory 존재 확인
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// root dir 생성
		os.Mkdir(dirPath, os.ModePerm)
	}

	return filePath
}

func (lg OLog) genDateTimeString() string {
	year, month, day := time.Now().Date()
	hh := time.Now().Hour()
	min := time.Now().Minute()
	sec := time.Now().Second()
	mills := time.Now().UnixNano() / 1000000

	strTime := fmt.Sprintf("%04d-%02d-%02d  %02d:%02d:%02d.%03d", year, int(month), day,
		hh, min, sec, (mills % 1000))
	return strTime
}
