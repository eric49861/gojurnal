package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
	Flag       int
}

func (this *Client) Menu() bool {
	var flag int
	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更新用户名")
	fmt.Println("0. 退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		this.Flag = flag
		return true
	} else {
		fmt.Println(">>>输入参数不合法<<<")
		return false
	}
}

func (this *Client) UpdateName() bool {
	fmt.Println(">>>请输入用户名")
	fmt.Scanln(&this.Name)
	sendMsg := "rename|" + this.Name + "\n"
	_, err := this.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("send error")
		return false
	}
	return true
}

func (this *Client) ResponseHandler() {
	//一旦可以从该连接中接收数据，就将数据拷贝到标准输出流
	io.Copy(os.Stdout, this.Conn)
	//相当于以下写法
	// for {
	// 	buf := make([]byte, 1024)
	// 	this.Conn.Read(buf)
	// 	fmt.Println(buf)
	// }
}

func (this *Client) PublicChat() {
	var msg string
	fmt.Println(">>>请输入信息:")
	fmt.Scanln(&msg)
	for msg != "exit" {
		if len(msg) != 0 {
			sendMsg := msg + "\n"
			_, err := this.Conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn write error")
				break
			}
		}
		msg = ""
		fmt.Println(">>>请输入信息:")
		fmt.Scanln(&msg)
	}
}

func (this *Client) SelectUser() {
	this.Conn.Write([]byte("who\n"))
}

func (this *Client) PrivateChat() {
	var remoteName string
	var chatMsg string
	this.SelectUser()
	fmt.Println(">>> 请输入聊天对象[用户名], exit退出")
	fmt.Scanln(&remoteName)
	for remoteName != "exit" {
		fmt.Println(">>>请输入聊天内容, exit退出")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			//消息不为空时则发送
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := this.Conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn write error")
					break
				}
			}
			chatMsg = ""
			fmt.Println(">>>请输入消息内容, exit退出:")
			fmt.Scanln(&chatMsg)
		}
		this.SelectUser()
		fmt.Println(">>> 请输入聊天对象[用户名], exit退出")
		fmt.Scanln(&remoteName)
	}
}
func (this *Client) Run() {
	for this.Flag != 0 {
		for this.Menu() != true {
		}
		switch this.Flag {
		case 1:
			this.PublicChat()
			break
		case 2:
			this.PrivateChat()
			break
		case 3:
			this.UpdateName()
			break
		case 0:
			//退出客户端
			break
		}
	}
}

func NewClient(serverIp string, ServerPort int) *Client {
	// 创建客户端
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: ServerPort,
		Flag:       888,
	}
	//连接服务端
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, ServerPort))
	if err != nil {
		fmt.Println("net.Dial failed")
		return nil
	}
	client.Conn = conn
	//返回client实例
	return client
}

var (
	serverIp   string
	ServerPort int
)

func init() {
	flag.StringVar(&serverIp, "ip", "192.168.0.149", "example: ip 127.0.0.1")
	flag.IntVar(&ServerPort, "port", 8888, "example: port 8888")
}

func main() {

	//解析命令行
	flag.Parse()
	client := NewClient(serverIp, ServerPort)
	if client == nil {
		fmt.Println("客户端创建失败")
		return
	}

	//启动一个goroutine接收服务端的消息
	go client.ResponseHandler()

	client.Run()
}
