package main

import (
	"log-b/internal/server"
	"log-b/internal/cache"
	"log-b/internal/db"
	"log-b/internal/cluster"
	"flag"
	"net"
	"log"
)

type StartUp struct {
	address    string
	listenPort string
	seed       bool 
}

var conf StartUp

func init() {
	var (
		addr = flag.String("host", "127.0.0.1", "a string")
		port = flag.String("port", "8080", "a string")
		sFlag = flag.Bool("seed", false, "a bool")
	)
	flag.Parse()

	conf = StartUp {
		address:    *addr,
		listenPort: *port,
		seed:       *sFlag,
	}
}

func startServer(node server.Replica) {
	go node.BindTCP()
	go node.ServeConns()
}

func main() {
	inMemoryMap := cache.Bcache{}
	diskStorage := db.NewDB()
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(r)
		}
	}()
	if !conf.seed {
		cluster.RegisterToSeed("127.0.0.1:6767", net.JoinHostPort(conf.address, conf.listenPort))
	}
	node := server.NewLoggedServer(conf.address, conf.listenPort, inMemoryMap, diskStorage)
	startServer(node)
}
