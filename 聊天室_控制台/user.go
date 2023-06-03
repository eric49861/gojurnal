package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr net.Addr
	C    chan string
	Conn net.Conn

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {

	user := &User{
		Name:   conn.RemoteAddr().String(),
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
	this.server.OnlineMap[this.Name] = this
	this.server.Lock.Unlock()

	//广播消息
	this.server.BroadCast(this, "Online")
}

//用户的下线业务
func (this *User) Offline() {
	this.server.BroadCast(this, "Offline")
	//将该用户从Onlinemap中移除
	this.server.Lock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.Lock.Unlock()
}

//给该user对应的客户端发送消息
func (this *User) SendMsg(msg string) {
	this.Conn.Write([]byte(msg + "\n"))
}

// 处理用户消息的功能
func (this *User) DoMessage(msg string) {
	if strings.Compare(msg, "who") == 0 {
		this.server.Lock.Lock()
		this.SendMsg("--------------Online Users Table--------------")
		for _, u := range this.server.OnlineMap {
			onlineMsg := "             [" + u.Name + "]"
			this.SendMsg(onlineMsg)
		}
		this.SendMsg("----------------------------------------------")
		this.server.Lock.Unlock()
	} else if len(msg) > 7 && strings.Compare(msg[:7], "rename|") == 0 {
		_, ok := this.server.OnlineMap[msg[7:]]
		if ok {
			this.SendMsg("rename failed, because name has been occupied")
		} else {
			this.SendMsg("rename successfully")
			this.server.Lock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.Name = msg[7:]
			this.server.OnlineMap[msg[7:]] = this
			this.server.Lock.Unlock()
		}
	} else if len(msg) > 4 && strings.Compare(msg[:3], "to|") == 0 {
		//解析私聊的对象
		format := strings.Split(msg, "|")
		name := format[1]
		if name == "" {
			this.SendMsg("incorrect format\nFor example: to|username|msg")
			return
		}
		//通过用户名获取user
		user, ok := this.server.OnlineMap[name]
		if !ok {
			this.SendMsg("Server hasn't user named " + name)
			return
		}
		//调用user的SendMsg方法转发消息
		user.SendMsg("(private)" + this.Name + ":" + format[2])
	} else {
		this.server.BroadCast(this, msg)
	}
}
