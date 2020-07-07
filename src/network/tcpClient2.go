package network

import (
	"fmt"
	"global"
	"log"
	"net"
	"time"
)

//---------------------------------------------------------------------------
// Protocol
//---------------------------------------------------------------------------
const (
	TCP_BUFFER_SIZE = 1024
	LCST_STX_CHAR   = 0x7E
	MIN_PACKET      = 5
)

//---------------------------------------------------------------------------
const (
	LCPT_STX    = 0
	LCPT_SIZE1  = LCPT_STX + 1
	LCPT_SIZE2  = LCPT_SIZE1 + 1
	LCPT_SEQ    = LCPT_SIZE2 + 1
	LCPT_OPCODE = LCPT_SEQ + 1
	LCPT_DATA   = LCPT_OPCODE + 1
)

//---------------------------------------------------------------------------
const (
	LCOPCD_STATE = 0x10
)

//---------------------------------------------------------------------------
// Struct
//---------------------------------------------------------------------------
type TCP_DESC struct {
	name   string
	port   string
	ipAddr string
	conn   net.Conn
}

//---------------------------------------------------------------------------
// Global
//---------------------------------------------------------------------------
var TcpDesc TCP_DESC = InitComm("WEB", "6000", "127.0.0.1") // name, port, ip

//---------------------------------------------------------------------------
// InitComm
//---------------------------------------------------------------------------
func InitComm(name string, port string, ipAddr string) TCP_DESC {
	tcpC := TCP_DESC{}
	tcpC.name = name
	tcpC.port = port
	tcpC.ipAddr = ipAddr

	return tcpC
}

//---------------------------------------------------------------------------
// Method
//---------------------------------------------------------------------------
// connect
//---------------------------------------------------------------------------
func Connect(chConnected chan bool) {
	log.Println("connect..")
	//netSrc := fmt.Sprintf("%s:%d", TcpDesc.ipAddr, TcpDesc.port)
	netSrc := TcpDesc.ipAddr + ":" + TcpDesc.port
	fmt.Println(netSrc)

	// connect
	conn, err := net.Dial("tcp", netSrc)

	if err != nil {
		global.Tcplog.Write(err.Error())
		log.Println("connect fail")
		chConnected <- false
		return

	}

	TcpDesc.conn = conn

	global.Tcplog.Write("Server ip: " + TcpDesc.ipAddr + ", port : " + string(TcpDesc.port) + "connected.!")
	log.Println("connect ok")

	chConnected <- true
}

//---------------------------------------------------------------------------
// rxHandler
//---------------------------------------------------------------------------
func rxHandler(n int, data []byte) {

}

//---------------------------------------------------------------------------
// msgHandler
//---------------------------------------------------------------------------
func msgHandler(data []byte) {

}

//---------------------------------------------------------------------------
// processLcState
//---------------------------------------------------------------------------
func processLcState(data []byte) {
	log.Println("processLcState")
}

//---------------------------------------------------------------------------
// manageRX
//---------------------------------------------------------------------------
func manageRX(chRecvData chan []byte, chConnected chan bool) {
	data := make([]byte, 1024)

	log.Println("rx (1)")

	n, err := TcpDesc.conn.Read(data)
	log.Println("rx (2)")
	if err != nil {
		log.Println("rx (3)")
		log.Println(err)
		TcpDesc.conn.Close()
		chConnected <- false
		return
	}
	log.Println("rx (4)")

	log.Println("Server send : ", string(data[:n]))
	chRecvData <- data
}

//---------------------------------------------------------------------------
// manageTX
//---------------------------------------------------------------------------
func manageTX(rcvData []byte, chConnected chan bool) {

	TcpDesc.conn.Write(rcvData)
}

//---------------------------------------------------------------------------
// Routine
//---------------------------------------------------------------------------
func Routine() {
	// channel
	chRecvData := make(chan []byte) // recv data
	chTcpConn := make(chan bool)    // connChk data

	// connect
	Connect(chTcpConn)

	// rx routine
	go func() {
		for {
			log.Println("rx routine (1)")
			manageRX(chRecvData, chTcpConn)
			log.Println("rx routine (2)")
			time.Sleep(time.Duration(1) * time.Second)
		}
	}()

	// manage connect, tx
	go func() {
		for {
			log.Println("manage routine (1)")
			//var rcvData []byte
			if rcvData, success := <-chRecvData; success {
				// send process
				manageTX(rcvData, chTcpConn)
			}

			log.Println("manage routine (2)")

			if connected, success := <-chTcpConn; success {
				// connect : false면 connect 수행
				if !connected {
					Connect(chTcpConn)
				}
			}

			log.Println("manage routine (3)")

			time.Sleep(time.Duration(1) * time.Second)
		}
	}()
}
