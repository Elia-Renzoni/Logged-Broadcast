package db

import (
	"database/sql"
	"time"
	"log-b/model"
)

type Storage interface {
	WriteMessage()
	RetrieveMessage()
}

type LogDB struct {
	instance *sql.DB
	pingerTime time.Duration
	tick *time.Ticker
	faultyStatus bool 
}

func NewDB() *LogDB {
	l := &LogDB{
		pingerTime: 3 * time.Second,
		faultyStatus: false,
	}
	
	go l.pinger()
	return l
}

func (l *LogDB) WriteMessage() {

}

// only for the recovery session
func (l *LogDB) RetrieveMessage() {

}

func (l *LogDB) pinger() {
	l.tick = time.NewTicker(l.pingerTime)
	defer l.tick.Stop()

	var err error

	for range l.tick.C {
		err = l.instance.Ping()
		if err != nil {
			l.faultyStatus = true
		} else {
			if l.faultyStatus {
				l.faultyStatus = false
			}
		}
	}
}
