package main

import (
	"log-b/internal/server"
)

func startServer(node server.Replica) {
	go node.BindTCP()
	go node.ServeConns()
}

// TODO-> add command line flag parsing
func main() {
	node := server.NewLoggedServer()
	startServer(node)
}
