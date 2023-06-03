package main

func main() {
	server := NewServer("192.168.0.149", "8888")
	server.Start()
}
