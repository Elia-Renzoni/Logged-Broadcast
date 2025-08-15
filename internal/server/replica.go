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
	*http.Server
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

func NewLoggedServer() *LoggedServer {
	return &LoggedServer{
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

