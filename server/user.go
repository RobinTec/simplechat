package server

import (
	"net"
)

type OnlineUser struct {
	Id       uint64
	Nickname string
	Conn     net.Conn
}

type OnlineUserList struct {
	UserList map[string]OnlineUser
	Lock     bool
}

func (this *OnlineUserList) Init() {
	this.UserList = make(map[string]OnlineUser)
	this.Lock = false
}

func (this *OnlineUserList) Append(onlineUser OnlineUser) (err error) {
	for this.Lock {
	}
	this.Lock = true
	this.UserList[string(utils.Uint64ToBytes(onlineUser.Id))] = onlineUser
	this.Lock = false
}

func (this *OnlineUserList) Delete(onlineUser OnlineUser) (err error) {
	for this.Lock {
	}
	this.Lock = true
	delete(this.UserList, string(utils.Uint64ToBytes(onlineUser.Id)))
	this.Lock = false
}
