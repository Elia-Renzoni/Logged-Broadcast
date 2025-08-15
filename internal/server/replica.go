package server

import (
	"net/http"
	"net"
	"log"
	"log-b/cache"
)

type Replica interface {
	BindTCP()
	ServeConns()
}

type LoggedServer struct {
	Router
	lst net.Listener
	completeAddr *net.TCPAddr
	addr string
	port string
	signaler chan struct{}
	wakeUp chan struct{}
	networkErr error
	inMemoryDB cache.Bcache
}

func NewLoggedServer(addr, port string, c cache.Bcache) *LoggedServer {
	tcpAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(addr, port))
	if err != nil {
		return nil
	}

	return &LoggedServer{
		addr: addr,
		port: port,
		completeAddr: tcpAddr,
		signaler: make(chan struct{}),
		wakeUp: make(chan struct{}),
		inMemoryDB: c,
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

func (ls *LoggedServer) ServeConns() {
	<- ls.wakeUp
	r := InitRouter(ls.inMemoryDB)
	http.Serve(ls.lst, http.HandlerFunc(r.ServeRequest))
}

