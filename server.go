package main

import (
	"fmt"
	"net"
)

type Server struct {
	IP   string
	Port string
}

func NewServer(ip, port string) *Server {
	server := &Server{
		IP:   ip,
		Port: port,
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
	fmt.Printf("%v\n", "建立连接成功")
	fmt.Printf("客户端信息: %v\n", conn.RemoteAddr().String())
}
