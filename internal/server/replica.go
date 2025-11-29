package server

import (
	"net/http"
	"net"
	"log"
	"log-b/internal/cache"
	"log-b/internal/db"
)

type Replica interface {
	BindTCP()
	ServeConns(joiner chan struct{})
}

type LoggedServer struct {
	Router
	lst          net.Listener
	completeAddr *net.TCPAddr
	addr         string
	port         string
	signaler     chan struct{}
	wakeUp       chan struct{}
	networkErr   error
	inMemoryDB   *cache.Bcache
	persistentDB *db.LogDB
	serverSecret string
}

func NewLoggedServer(addr, port string, c *cache.Bcache, d *db.LogDB, secret string) *LoggedServer {
	tcpAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(addr, port))
	if err != nil {
		return nil
	}

	return &LoggedServer{
		addr:         addr,
		port:         port,
		completeAddr: tcpAddr,
		signaler:     make(chan struct{}),
		wakeUp:       make(chan struct{}),
		inMemoryDB:   c,
		persistentDB: d,
		serverSecret: secret,
	}
}

func (ls *LoggedServer) BindTCP() {
	ls.lst, ls.networkErr = net.ListenTCP("tcp", ls.completeAddr)
	if ls.networkErr != nil {
		ls.signaler <- struct{}{}
		log.Fatal(ls.networkErr.Error())
		return
	}

	ls.wakeUp <- struct{}{}
}

func (ls *LoggedServer) ServeConns(joinSync chan struct{}) {
	<- ls.wakeUp
	r := InitRouter(ls.inMemoryDB, ls.persistentDB, ls.serverSecret)
	log.Println("Server ON...")
	err := http.Serve(ls.lst, http.HandlerFunc(r.ServeRequest))
	if err != nil {
		log.Fatal(err)
	}
	<- joinSync
}
