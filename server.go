package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	IP        string
	Port      string
	OnlineMap map[string]*User
	Lock      sync.RWMutex //对map进行操作时进行加锁
	Message   chan string
}

func NewServer(ip, port string) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Lock:      sync.RWMutex{},
		Message:   make(chan string),
	}
	return server
}

func (this *Server) Start() {
	// 创建监听套接字
	listener, err := net.Listen("tcp4", this.IP+":"+this.Port)
	if err != nil {
		fmt.Printf("%v\n", "套接字创建失败")
		return
	}
	//程序退出时关闭监听套接字
	defer listener.Close()

	go this.MessageListener()
	//循环监听，相当于while(true)
	for {
		//监听连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("%v\n", "接受连接失败")
			return
		}
		//启动一个协程处理连接
		go this.Handler(conn)
	}
}

//处理连接
func (this *Server) Handler(conn net.Conn) {
	//此时该客户端已成功与服务端建立连接，将该客户端加入在线用户列表中
	user := NewUser(conn, this)
	user.Online()
	isAlive := make(chan bool)
	// 处理用户发过来的消息
	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			} else if err != nil && err != io.EOF {
				return
			}
			msg := string(buf[:n-1])
			user.DoMessage(msg)

			isAlive <- true //	标记该用户是活跃的
		}
	}()

	for {
		select {
		case <-isAlive:
			//do nothing
		case <-time.After(time.Second * 100):
			//剔除该用户
			user.SendMsg("You are forced offline because of the long slience")
			//销毁相关的资源
			close(user.C)
			conn.Close()

			this.Lock.Lock()
			delete(this.OnlineMap, user.Name)
			this.Lock.Unlock()
			return
		}
	}
}

func (this *Server) BroadCast(user *User, msg string) {
	//  通过广播管道将消息广播给所有的在线用户
	sendMsg := "[" + user.Addr.String() + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

//创建一个Listener，用于监听广播管道，如果有消息，就遍历所有用户，将广播消息广播出去
func (this *Server) MessageListener() {
	for {
		msg := <-this.Message
		this.Lock.Lock()
		for _, u := range this.OnlineMap {
			u.C <- msg
		}
		this.Lock.Unlock()
	}
}
