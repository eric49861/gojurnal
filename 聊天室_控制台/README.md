# gojurnal
go语言学习之旅
该仓库通过一步步构建一个简易的即时通信系统来熟悉go语言的语法和数据结构，通过该项目可以学习以下内容：
- 变量声明
- 循环结构和条件选择结构
- go语言函数声明、匿名函数
- go语言数据结构，切片、map
- 如何使用go语言实现面向对象编程
- 如何创建goroutine以及如何使用channel进行协程间通信
- 网络通信
- 项目依赖管理(go modules)
## 项目功能列表
1. 用户的上下线广播
2. 群聊
3. 查看当前在线用户列表
4. 自定义用户名的功能
5. 长时间不操作自动剔除
6. 私聊功能

## 编译指令
> go build -o ./out/server.exe main.go server.go user.go

## 运行服务端
> ./out/server.exe

##  客户端

### 编译客户端代码
> go build -o ./out/client.exe client.go

## 运行
运行前需要修改服务端的IP地址和监听的端口号

1. 首先运行服务端
   > ./out/server.exe
2. 运行客户端
   > ./out/client.exe ip 服务端IP地址 port 端口号

服务端的ip地址和端口号并不是必须的，但是需要修改默认值