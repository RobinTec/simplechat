package server

import (
	"log"
	"net"
	"simplechat/improtocal"
)

func closeConn(conn net.Conn) {
	log.Printf("connection closed: %s\n", conn.RemoteAddr())
	conn.Close()
}

func handleConn(conn net.Conn, queue *MsgQueue) {
	defer closeConn(conn)
	log.Printf("new connection from %s\n", conn.RemoteAddr())
	for {
		recv, err := improtocal.Read(conn)
		if err != nil {
			break
		}

		switch recv.Msg_type {
		case improtocal.LOGIN:
			if err = improtocal.SendLoginAck(conn, recv.Payload); err != nil {
				log.Printf("%s", err)
				break
			}
		case improtocal.SEND_TO_ONE:
			if err = improtocal.SendPersonalMsgAck(conn, recv.Payload); err != nil {
				log.Printf("%s", err)
				break
			}
		case improtocal.SEND_TO_GROUP:
			if err = improtocal.SendGroupMsgAck(conn, recv.Payload); err != nil {
				log.Printf("%s", err)
				break
			}
		default:
			log.Printf("Received invalid package: %s", recv)
			break
		}
	}
}

func Start() {
	log.Println("server starting...")
	listener, err := net.Listen("tcp", ":9999")
	if err != nil {
		listenPort := 9999
		log.Println("Listen to port %s failed: %s", listenPort, err)
		return
	}
	defer listener.Close()

	// 实际应用中的在线用户列表和消息队列不能写在这里
	// 这些数据应该是共享的，是一个单独划分出来的模块
	// 有利于多实例部署
	var onlineUsers OnlineUserList
	var msgQueue MsgQueue
	onlineUsers.Init()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept connection error: %s", err)
			continue
		}
		handleConn(conn, &msgQueue)
	}

}
