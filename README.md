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

## 编译指令
> go build -o ./out/server.exe main.go server.go user.go

## 运行服务端
> ./out/server.exe