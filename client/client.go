package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"simplechat/improtocal"
)

func Start() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("username: ")
	username, _ := reader.ReadString('\n')
	username = username[:len(username)-1]

	conn, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
	}

	improtocal.SendLogin(conn, username)
	var pkg improtocal.Package
	pkg, err = improtocal.Read(conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pkg.Payload)
	pkg, err = improtocal.Read(conn)
}
