package improtocal

import (
	"errors"
	"net"
	"simplechat/utils"
)

const LOGIN = uint16(0)
const LOGIN_ACK = uint16(1)
const SEND_TO_ONE = uint16(2)
const SEND_TO_GROUP = uint16(3)
const SEND_ACK = uint16(4)

const SVR_RECV_MSG_OK = uint16(1)
const SVR_SEND_MSG_TO_USR_OD = uint16(2)
const SVR_PERMISSION_DENIED = uint16(3)
const SVR_INTERNAL_ERROR = uint16(4)

type Package struct {
	Msg_type uint16
	Src_id   uint64
	Dst_id   uint64
	Error    uint16
	Len      uint16
	Payload  string
}

func Read(conn net.Conn) (pkg Package, err error) {
	pkg.Msg_type, err = getMsgType(conn)
	if err == nil {
		switch pkg.Msg_type {
		case LOGIN:
			pkg.Len, err = getPayloadLen(conn)
			pkg.Payload, err = getPayload(conn, pkg.Len)
		case LOGIN_ACK:
			pkg.Dst_id, err = getDstId(conn)
			pkg.Len, err = getPayloadLen(conn)
			pkg.Payload, err = getPayload(conn, pkg.Len)
		case SEND_TO_ONE:
			pkg.Src_id, err = getSrcId(conn)
			pkg.Dst_id, err = getDstId(conn)
			pkg.Len, err = getPayloadLen(conn)
			pkg.Payload, err = getPayload(conn, pkg.Len)
		case SEND_TO_GROUP:
			pkg.Src_id, err = getSrcId(conn)
			pkg.Len, err = getPayloadLen(conn)
			pkg.Payload, err = getPayload(conn, pkg.Len)
		case SEND_ACK:
			pkg.Dst_id, err = getDstId(conn)
			pkg.Error, err = getError(conn)
			pkg.Len, err = getPayloadLen(conn)
			pkg.Payload, err = getPayload(conn, pkg.Len)
		default:
			err = errors.New("Msg type invalid")
		}
	}
	return
}

func Write(conn net.Conn, pkg Package) (err error) {
	pkg_byte := utils.Uint16ToBytes(pkg.Msg_type)
	pkg.Len = uint16(len(pkg.Payload))
	switch pkg.Msg_type {
	case LOGIN:
		pkg_byte = append(pkg_byte, utils.Uint16ToBytes(pkg.Len)...)
		pkg_byte = append(pkg_byte, []byte(pkg.Payload)...)
	case LOGIN_ACK:
		pkg_byte = append(pkg_byte, utils.Uint64ToBytes(pkg.Dst_id)...)
		pkg_byte = append(pkg_byte, utils.Uint16ToBytes(pkg.Len)...)
		pkg_byte = append(pkg_byte, []byte(pkg.Payload)...)
	case SEND_TO_ONE:
		pkg_byte = append(pkg_byte, utils.Uint64ToBytes(pkg.Src_id)...)
		pkg_byte = append(pkg_byte, utils.Uint64ToBytes(pkg.Dst_id)...)
		pkg_byte = append(pkg_byte, utils.Uint16ToBytes(pkg.Len)...)
		pkg_byte = append(pkg_byte, []byte(pkg.Payload)...)
	case SEND_TO_GROUP:
		pkg_byte = append(pkg_byte, utils.Uint64ToBytes(pkg.Src_id)...)
		pkg_byte = append(pkg_byte, utils.Uint16ToBytes(pkg.Len)...)
		pkg_byte = append(pkg_byte, []byte(pkg.Payload)...)
	case SEND_ACK:
		pkg_byte = append(pkg_byte, utils.Uint64ToBytes(pkg.Dst_id)...)
		pkg_byte = append(pkg_byte, utils.Uint16ToBytes(pkg.Error)...)
		pkg_byte = append(pkg_byte, utils.Uint16ToBytes(pkg.Len)...)
		pkg_byte = append(pkg_byte, []byte(pkg.Payload)...)
	default:
		err = errors.New("Msg type invalid")
		return
	}

	_, err = conn.Write(pkg_byte)
	return
}

func getMsgType(conn net.Conn) (msg_type uint16, err error) {
	msg_type_byte, err := readBuffer(conn, 2)
	msg_type = utils.BytesToUint16(msg_type_byte)
	return
}

func getSrcId(conn net.Conn) (src_id uint64, err error) {
	src_id_byte, err := readBuffer(conn, 8)
	src_id = utils.BytesToUint64(src_id_byte)
	return
}

func getDstId(conn net.Conn) (dst_id uint64, err error) {
	dst_id_byte, err := readBuffer(conn, 8)
	dst_id = utils.BytesToUint64(dst_id_byte)
	return
}

func getError(conn net.Conn) (pkg_err uint16, err error) {
	pkg_err_byte, err := readBuffer(conn, 2)
	pkg_err = utils.BytesToUint16(pkg_err_byte)
	return
}

func getPayloadLen(conn net.Conn) (payload_len uint16, err error) {
	payload_len_byte, err := readBuffer(conn, 2)
	payload_len = utils.BytesToUint16(payload_len_byte)
	return
}

func getPayload(conn net.Conn, payloadLen uint16) (payload string, err error) {
	payload_byte, err := readBuffer(conn, int(payloadLen))
	payload = string(payload_byte)
	return
}

func readBuffer(conn net.Conn, totalBytesToRead int) ([]byte, error) {
	var buf = make([]byte, totalBytesToRead)
	_, err := conn.Read(buf)
	return buf, err
}
