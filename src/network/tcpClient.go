package network

// import (
// 	"fmt"
// 	"genLib"
// 	"global"
// 	"log"
// 	"net"
// )

// //---------------------------------------------------------------------------
// // Protocol
// //---------------------------------------------------------------------------
// const (
// 	TCP_BUFFER_SIZE = 1024
// 	LCST_STX_CHAR   = 0x7E
// 	MIN_PACKET      = 5
// )

// //---------------------------------------------------------------------------
// const (
// 	LCPT_STX    = 0
// 	LCPT_SIZE1  = LCPT_STX + 1
// 	LCPT_SIZE2  = LCPT_SIZE1 + 1
// 	LCPT_SEQ    = LCPT_SIZE2 + 1
// 	LCPT_OPCODE = LCPT_SEQ + 1
// 	LCPT_DATA   = LCPT_OPCODE + 1
// )

// //---------------------------------------------------------------------------
// const (
// 	LCOPCD_STATE = 0x10
// )

// //---------------------------------------------------------------------------
// // Struct
// //---------------------------------------------------------------------------
// type TcpClient struct {
// 	connected bool
// 	name      string
// 	port      uint16
// 	ipAddr    string
// 	m_data    []byte
// 	m_length  int
// 	conn      net.Conn
// }

// //---------------------------------------------------------------------------
// // Global
// //---------------------------------------------------------------------------
// var Client *TcpClient = InitComm("WEB", 6000, "127.0.0.1") // name, port, ip

// //---------------------------------------------------------------------------
// // InitComm
// //---------------------------------------------------------------------------
// func InitComm(name string, port uint16, ipAddr string) *TcpClient {
// 	tcpC := TcpClient{}
// 	tcpC.name = name
// 	tcpC.port = port
// 	tcpC.ipAddr = ipAddr

// 	return &tcpC
// }

// //---------------------------------------------------------------------------
// // Method
// //---------------------------------------------------------------------------
// // connect
// //---------------------------------------------------------------------------
// func (tcp *TcpClient) Connect() bool {
// 	log.Println("connect..")

// 	//netSrc := tcp.ipAddr + ":" + string(tcp.port)
// 	netSrc := fmt.Sprintf("%s:%d", tcp.ipAddr, tcp.port)
// 	log.Println(netSrc)
// 	conn, err := net.Dial("tcp", netSrc)

// 	if err != nil {
// 		global.Tcplog.Write(err.Error())
// 		tcp.connected = false
// 		log.Println("connect fail")
// 		return (tcp.connected)
// 	}

// 	global.Tcplog.Write("Server ip: " + tcp.ipAddr + ", port : " + string(tcp.port) + "connected.!")

// 	tcp.conn = conn
// 	tcp.connected = true
// 	log.Println("connect ok")
// 	return (tcp.connected)
// }

// //---------------------------------------------------------------------------
// // rxHandler
// //---------------------------------------------------------------------------
// func (tcp *TcpClient) rxHandler(n int, data []byte) {
// 	global.Tcplog.Dump("RX", data)

// 	// check STX
// 	if data[LCPT_STX] != LCST_STX_CHAR {
// 		strStx := fmt.Sprintf("%02X", data[LCPT_STX])
// 		global.Tcplog.Write("Invalid STX : [" + strStx + "]")
// 	}

// 	// check size
// 	dataSize := genLib.GetNumber(data, LCPT_SIZE1, 2, genLib.ED_BIG)
// 	log.Println("size : ", dataSize)

// 	tcp.msgHandler(data)
// }

// //---------------------------------------------------------------------------
// // msgHandler
// //---------------------------------------------------------------------------
// func (tcp *TcpClient) msgHandler(data []byte) {
// 	code := data[LCPT_OPCODE]

// 	switch code {
// 	case LCOPCD_STATE:
// 		tcp.processLcState(data)
// 		break
// 	default:
// 		hCode := fmt.Sprintf("%02X", code)
// 		global.Tcplog.Write("Invalid CODE : [" + hCode + "]")
// 	}
// }

// //---------------------------------------------------------------------------
// // processLcState
// //---------------------------------------------------------------------------
// func (tcp *TcpClient) processLcState(data []byte) {
// 	log.Println("processLcState")
// }

// //---------------------------------------------------------------------------
// // manageRX
// //---------------------------------------------------------------------------
// func (tcp *TcpClient) manageRX(chRecvData chan []byte, chConn chan bool) bool {

// 	data := make([]byte, TCP_BUFFER_SIZE)

// 	n, err := tcp.conn.Read(data)
// 	if err != nil {

// 		global.Tcplog.Write("manageRX err : " + err.Error())
// 		return (false)
// 	}

// 	tcp.rxHandler(n, data)

// 	return (true)
// }

// //---------------------------------------------------------------------------
// // manageTX
// //---------------------------------------------------------------------------
// func (tcp *TcpClient) manageTX(chRecvData chan []byte, chConn chan bool) bool {

// 	data := "hello"
// 	if tcp.connected {
// 		tcp.conn.Write([]byte(data))
// 	}
// 	return (true)
// }

// //---------------------------------------------------------------------------
// // Manage
// //---------------------------------------------------------------------------
// func (tcp *TcpClient) Manage() {
// 	// manage connection
// 	if !tcp.connected {
// 		if !tcp.Connect() {
// 			return
// 		}
// 	}

// 	// // manage send
// 	// if !tcp.manageTX() {
// 	// 	return
// 	// }

// 	// // manage recv
// 	// if !tcp.manageRX() {
// 	// 	return
// 	// }
// }

// //---------------------------------------------------------------------------
// // Routine
// //---------------------------------------------------------------------------
// func (tcp *TcpClient) Routine() {
// 	chRecvData := make(chan []byte)
// 	chTcpConn := make(chan bool)

// 	// thread (go routine)

// 	// for {
// 	// 	log.Println("Routine")
// 	// 	tcp.Manage()
// 	// 	//time.Sleep(time.Duration(100) * time.Millisecond) // 100 msec
// 	// 	time.Sleep(time.Duration(1) * time.Second) // 100 msec
// 	// }

// 	// // rx routine
// 	// go func() {
// 	// 	for {
// 	// 		tcp.manageRX(chRecvData, chTcpConn)
// 	// 		time.Sleep(time.Duration(1) * time.Second)
// 	// 	}
// 	// }()

// 	// for {
// 	// 	if rcvData, success := <-chRecvData; success {
// 	// 		// send process
// 	// 	}

// 	// 	if connChk, success := <-chTcpConn; success {
// 	// 		// connect 수행
// 	// 	}

// 	// }

// }
