package main

import "net"

type User struct {
	Addr net.Addr
	C    chan string
	Conn net.Conn
}

func NewUser(conn net.Conn) *User {

	user := &User{
		Addr: conn.RemoteAddr(),
		C:    make(chan string),
		Conn: conn,
	}
	go user.CListener()

	return user
}

func (this *User) CListener() {
	for {
		msg := <-this.C
		this.Conn.Write([]byte(msg + "\n"))
	}
}
