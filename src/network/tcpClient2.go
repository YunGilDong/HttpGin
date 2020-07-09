package network

import (
	"data"
	"fmt"
	"genLib"
	"global"
	"log"
	"net"
	"time"
	"trffic_obj"
)

//---------------------------------------------------------------------------
// Protocol
//---------------------------------------------------------------------------
const (
	TCP_BUFFER_SIZE = 1024
	LCST_STX_CHAR   = 0x7E
	MIN_PACKET      = 5 // header size
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
// Rx statsu
//---------------------------------------------------------------------------
const (
	RXST_STX    = 0
	RXST_SIZE1  = 1
	RXST_SIZE2  = 2
	RXST_SEQ    = 3
	RXST_OPCODE = 4
	RXST_DATA   = 5
)

//---------------------------------------------------------------------------
// Struct
//---------------------------------------------------------------------------

type RecvData struct {
	data   []byte
	length int
}

type TCPinfo struct {
	name      string
	port      string
	ipAddr    string
	connected bool
	conn      net.Conn
	m_data    []byte
	m_index   int // m_data last index
	m_length  int // data size
	rx_status int // RXST_STX (0) ~ RXST_DATA (5)
}

//---------------------------------------------------------------------------
// Global
//---------------------------------------------------------------------------
var TcpClient TCPinfo = InitComm("WEB", "6000", "127.0.0.1") // name, port, ip

//---------------------------------------------------------------------------
// InitComm
//---------------------------------------------------------------------------
func InitComm(name string, port string, ipAddr string) TCPinfo {
	tcpC := TCPinfo{}
	tcpC.name = name
	tcpC.port = port
	tcpC.ipAddr = ipAddr

	return tcpC
}

func minInt(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

//---------------------------------------------------------------------------
// Method
//---------------------------------------------------------------------------
// connect
//---------------------------------------------------------------------------
func Connect(ch_connected chan bool) {
	log.Println("connect..")
	//netSrc := fmt.Sprintf("%s:%d", TcpClient.ipAddr, TcpClient.port)
	netSrc := TcpClient.ipAddr + ":" + TcpClient.port
	log.Println(netSrc)

	// connect
	conn, err := net.Dial("tcp", netSrc)
	TcpClient.conn = conn
	resetRxStatus()

	if err != nil {
		global.Tcplog.Write(err.Error())
		TcpClient.connected = false
		log.Println("connect fail")
		return
	}

	TcpClient.connected = true
	global.Tcplog.Write("Server ip: " + TcpClient.ipAddr + ", port : " + string(TcpClient.port) + "connected.!")
	log.Println("connect ok")
}

func resetRxStatus() {

	// clear
	TcpClient.m_data = TcpClient.m_data[:0]
	TcpClient.m_index = 0
	TcpClient.rx_status = RXST_STX
}

func SetRxStatus(rxState int, readCount int) {

	TcpClient.rx_status = rxState
	TcpClient.m_index += readCount

	if rxState == RXST_STX {
		resetRxStatus()
	}
}

func rxHandler(data []byte, length int) {

	global.Tcplog.Dump("RX", data, length)
	// check header
	for idx := 0; idx < length; idx++ {

		m_rxState := TcpClient.rx_status

		//log.Println("m_rxState : ", m_rxState)

		switch m_rxState {
		case RXST_STX:
			if data[idx] == LCST_STX_CHAR {
				TcpClient.m_data = append(TcpClient.m_data, data[idx])
				SetRxStatus(RXST_SIZE1, 1)

			} else {
				TcpClient.m_data = append(TcpClient.m_data, data[idx])
				SetRxStatus(RXST_SIZE1, 1)
				//SetRxStatus(RXST_STX, 0)
			}
		case RXST_SIZE1:
			TcpClient.m_data = append(TcpClient.m_data, data[idx])
			SetRxStatus(RXST_SIZE2, 1)

		case RXST_SIZE2:
			TcpClient.m_data = append(TcpClient.m_data, data[idx])
			dataLen := genLib.GetNumber(TcpClient.m_data, LCPT_SIZE1, 2, genLib.ED_BIG)
			TcpClient.m_length = dataLen + MIN_PACKET
			SetRxStatus(RXST_SEQ, 1)

		case RXST_SEQ:
			TcpClient.m_data = append(TcpClient.m_data, data[idx])
			SetRxStatus(RXST_OPCODE, 1)

		case RXST_OPCODE:
			TcpClient.m_data = append(TcpClient.m_data, data[idx])
			SetRxStatus(RXST_DATA, 1)

		case RXST_DATA:
			m_index := TcpClient.m_index

			remainCount := length - idx
			requestCount := TcpClient.m_length - m_index
			requestCount = minInt(requestCount, remainCount)

			lastIdx := idx + requestCount
			TcpClient.m_data = append(TcpClient.m_data, data[idx:lastIdx]...)
			SetRxStatus(RXST_DATA, requestCount)
			idx += requestCount - 1

			if m_index == TcpClient.m_length {
				msgHandler()
				SetRxStatus(RXST_STX, 0)
			} else if m_index < TcpClient.m_length {
				log.Println("(2)", m_index, TcpClient.m_length)
				continue
			} else if m_index > TcpClient.m_length {
				log.Println("(3)", m_index, TcpClient.m_length)
				SetRxStatus(RXST_STX, 0)
			}
		}
	}
}

func msgHandler() {
	// opcode별로 처리
	global.Tcplog.Dump("MSG HND", TcpClient.m_data, TcpClient.m_length)
	code := TcpClient.m_data[LCPT_OPCODE]

	switch code {
	case LCOPCD_STATE:
		processLcStatus()
	default:
		log.Printf("Undefined opcode %02X", code)

	}

}

func processLcStatus() {
	lcdata := TcpClient.m_data[0:] // copy data
	datacount := genLib.GetNumber(lcdata, LCPT_DATA, 2, genLib.ED_BIG)
	const lcStateSize = 28
	const STATE_DATA_IDX = LCPT_DATA + 2

	for idx := 0; idx < datacount; idx++ {
		lcObj := data.LC{}

		delta := 2 // id, group
		offsetIdx := lcStateSize*idx + STATE_DATA_IDX

		lcObj.LC_ID = int(lcdata[offsetIdx+0])
		lcObj.GRP_ID = int(lcdata[offsetIdx+1])

		lcObj.State.OprMode = int(lcdata[offsetIdx+delta+0])

		//id, group, opr, conflict, light, flash, doost, commst
		// log.Printf(" id[%02X] grp[%02X] opr[%02X] status[%02X]   comm[%02X] "
		// 	, lcObj.LC_ID, lcObj.GRP_ID	, lcdata[offsetIdx+delta+0], lcdata[offsetIdx+delta+3], lcdata[offsetIdx+delta+25])

		log.Printf(" id[%02X] grp[%02X] opr[%02X] status[%02X] comm[%02X]",
			lcObj.LC_ID,
			lcObj.GRP_ID,
			lcdata[offsetIdx+delta+0],
			lcdata[offsetIdx+delta+3],
			lcdata[offsetIdx+delta+25])

		// Conflict
		if lcdata[offsetIdx+delta+3]&0x08 > 0 {
			log.Println("conflict 1")
			lcObj.State.ConflictSt = 1
		} else {
			log.Println("conflict 0")
			lcObj.State.ConflictSt = 0
		}

		// Light
		if lcdata[offsetIdx+delta+3]&0x04 > 0 {
			log.Println("Light 1")
			lcObj.State.LightOffSt = 1
		} else {
			log.Println("Light 0")
			lcObj.State.LightOffSt = 0
		}

		// Flash
		if lcdata[offsetIdx+delta+3]&0x02 > 0 {
			log.Println("Flash 1")
			lcObj.State.FlashSt = 1
		} else {
			log.Println("Flash 0")
			lcObj.State.FlashSt = 0
		}

		// Door Status
		if lcdata[offsetIdx+delta+3]&0x01 > 0 {
			log.Println("Door 1")
			lcObj.State.DoorSt = 1
		} else {
			log.Println("Door 0")
			lcObj.State.DoorSt = 0
		}

		// Comm status
		lcObj.State.CommSt = int(lcdata[offsetIdx+delta+25])

		// Set Lcobjects
		trffic_obj.SetLcObjecState(lcObj)
	}
}

//---------------------------------------------------------------------------
// manageRX
//---------------------------------------------------------------------------
func manageRX(ch_connected chan bool, ch_recvdata chan RecvData) {
	for {

		if TcpClient.connected {

			data := make([]byte, 1024)

			n, err := TcpClient.conn.Read(data)
			if err != nil {
				log.Println(err)
				TcpClient.connected = false
				ch_connected <- false
			}

			rcvData := RecvData{}
			rcvData.data = data
			rcvData.length = n
			rxHandler(data, n)

			ch_recvdata <- rcvData
		}
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
}

//---------------------------------------------------------------------------
// manageTX
//---------------------------------------------------------------------------
func manageTX(ch_connected chan bool, ch_recvdata chan RecvData) {

	for {

		if TcpClient.connected {
			s := "hello"
			_, err := TcpClient.conn.Write([]byte(s))
			if err != nil {
				fmt.Println("err")
				ch_connected <- false
			}

			// recv message handler
			select {
			// case rcvdata := <-ch_recvdata:
			// 	log.Println("manageTX recv data! size : ", rcvdata.length)
			case <-ch_recvdata:

			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func manage(ch_connected chan bool) {

	for {
		select {
		case connected := <-ch_connected:
			if !connected {
				if TcpClient.connected {
					TcpClient.connected = false
					TcpClient.conn.Close()
				}
				Connect(ch_connected)
			}
		default:
			if !TcpClient.connected {
				Connect(ch_connected)
			}
		}

		time.Sleep(time.Duration(100) * time.Millisecond)
	}
}

//---------------------------------------------------------------------------
// Routine
//---------------------------------------------------------------------------
func Routine() {
	ch_connected := make(chan bool)
	ch_recvdata := make(chan RecvData)

	// tcp connect
	Connect(ch_connected)

	go manage(ch_connected)                // connect manage
	go manageRX(ch_connected, ch_recvdata) // rx manage
	go manageTX(ch_connected, ch_recvdata) // tx manage
}
