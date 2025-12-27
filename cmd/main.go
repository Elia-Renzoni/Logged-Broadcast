package main

import (
	"flag"
	"log"
	"log-b/internal/cache"
	"log-b/internal/cluster"
	"log-b/internal/db"
	"log-b/internal/server"
	"net"
)

type StartUp struct {
	address    string
	listenPort string
	seed       bool
	secret     string
}

var conf StartUp

func init() {
	var (
		addr   = flag.String("host", "127.0.0.1", "a string")
		port   = flag.String("port", "8080", "a string")
		sFlag  = flag.Bool("seed", false, "a bool")
		secret = flag.String("secret", "foo", "a string")
	)
	flag.Parse()

	conf = StartUp{
		address:    *addr,
		listenPort: *port,
		seed:       *sFlag,
		secret:     *secret,
	}
}

func startServer(node server.Replica, joiner chan struct{}) {
	go node.BindTCP()
	go node.ServeConns(joiner)
}

func main() {
	var (
		inMemoryMap cache.MemoryCache = &cache.Bcache{}
		diskStorage db.Storage        = db.NewDB()
	)
	defer func() {
		if r := recover(); r != nil {
			log.Println("recovering...")
			log.Fatal(r)
		}
	}()
	if !conf.seed {
		go cluster.RegisterToSeed("127.0.0.1:6767", net.JoinHostPort(conf.address, conf.listenPort))
	}

	inMemoryMap.OpenDB()
	defer inMemoryMap.CloseDB()
	if err := diskStorage.StartDB(); err != nil {
		log.Fatal(err)
		return
	}
	// only for debugging
	log.Printf("[Interface]: %v\n", diskStorage)
	defer diskStorage.ShutdownDB()
	node := server.NewLoggedServer(
		conf.address,
		conf.listenPort,
		inMemoryMap,
		diskStorage,
		conf.secret,
	)
	joiner := make(chan struct{})
	startServer(node, joiner)
	<-joiner
}
