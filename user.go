package main

import "net"

type User struct {
	Addr net.Addr
	C    chan string
	Conn net.Conn

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {

	user := &User{
		Addr:   conn.RemoteAddr(),
		C:      make(chan string),
		Conn:   conn,
		server: server,
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

// 用户的上线业务
func (this *User) Online() {
	this.server.Lock.Lock()
	this.server.OnlineMap[this.Addr.String()] = this
	this.server.Lock.Unlock()

	//广播消息
	this.server.BroadCast(this, "Online")
}

//用户的下线业务
func (this *User) Offline() {
	this.server.BroadCast(this, "["+this.Addr.String()+" Offline]")
	//将该用户从Onlinemap中移除
	this.server.Lock.Lock()
	delete(this.server.OnlineMap, this.Addr.String())
	this.server.Lock.Unlock()
}

// 处理用户消息的功能
func (this *User) DoMessage(msg string) {
	this.server.BroadCast(this, msg)
}
