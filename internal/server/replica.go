package server

import (
	"net/http"
	"net"
)

type Replica interface {
	BindTCP()
	ServeTCP()
}

type LoggedServer struct {
	*http.Server
	Router
	ls net.Listener
	completeAddr *net.TCPAddr
	addr string
	port string
	signaler chan struct{}
}

func NewLoggedServer() *LoggedServer {
	return &LoggedServer{
	}
}

func (ls *LoggedServer) BindTCP() {

}

func (ls *LoggedServer) ServeTCP() {

}

func (ls *LoggedServer) acceptConn() net.Conn {

}
