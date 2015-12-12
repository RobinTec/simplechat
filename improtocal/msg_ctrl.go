package improtocal

import (
	"fmt"
	"log"
	"math/rand"
	"net"
)

func SendLogin(conn net.Conn, nickname string) (err error) {
	var pkg Package
	pkg.Msg_type = LOGIN
	pkg.Payload = nickname
	err = Write(conn)
	return
}

func SendLoginAck(conn net.Conn, nickname string) (err error) {
	var pkg improtocal.Package
	pkg.Msg_type = LOGIN_ACK
	pkg.Src_id = uint64(1)
	pkg.Dst_id = uint64(rand.Int63())
	pkg.Payload = fmt.Sprintf("Welcome %s , your temporary id is %d", nickname, pkg.Dst_id)
	log.Printf("%s login success from %s\n", nickname, conn.RemoteAddr())
	err = improtocal.Write(conn, pkg)
	return err
}

func SendPersonalMsg(conn net.Conn, from uint64, to uint64, msgContext string) (err error) {
	if err = sendMsg(conn, SEND_TO_GROUP, from, to, msgContext); err != nil {
		log.Printf("Send personal msg faild\n")
	}
	return
}

func SendPersonalMsgAck(conn net.Conn, cli_id uint64) (err error) {
	if err = sendMsgAck(conn, cli_id); err != nil {
		log.Printf("Send personal msg ack to %s failed\n", conn.RemoteAddr())
	}
}

func SendGroupMsg(conn net.Conn, from uint64, to uint64, msgContext string) (err error) {
	if err = sendMsg(conn, SEND_TO_GROUP, from, to, msgContext); err != nil {
		log.Printf("Send group msg faild\n")
	}
	return
}

func SendGroupMsgAck(conn net.Conn, cli_id uint64) (err error) {
	if err = sendMsgAck(conn, cli_id); err != nil {
		log.Printf("Send group msg ack to %s failed\n", conn.RemoteAddr())
	}
}

func sendMsg(conn net.Conn, msg_type uint64, from uint64, to uint64, msgContext string) (err error) {
	var pkg Package
	pkg.Msg_type = msg_type
	pkg.Src_id = from
	pkg.Dst_id = to
	pkg.Payload = msgContext
	err = Write(conn, pkg)
	return
}

func sendMsgAck(conn net.Conn, cli_id uint64) (err error) {
	var pkg improtocal.Package
	pkg.Dst_id = cli_id
	pkg.Error = SVR_RECV_MSG_OK
	pkg.Payload = "Server received msg success"
	err = improtocal.Write(conn, pkg)
	return
}
