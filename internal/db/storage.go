package db

import (
	"database/sql"
	"time"
	"log-b/model"
	"sync/atomic"
	"errors"
	"context"
	"slices"
)

type Storage interface {
	StartDB() error
	WriteMessage(content model.PersistentMessage, opType uint8) error
	RetrieveMessage()
	DeleteMessage(searchKey string) error
	ChangeStatus()
	ShutdownDB()
}

type LogDB struct {
	instance     *sql.DB
	pingerTime   time.Duration
	tick         *time.Ticker
	faultyStatus atomic.Bool
	dbCtx        context.Context
	tx           *sql.Tx
}

func NewDB() *LogDB {
	return &LogDB{
		pingerTime: 3 * time.Second,
		dbCtx: context.Background(),
	}
}

func (l *LogDB) WriteMessage(content model.PersistentMessage, opType uint8) error {
	var (
		senderInfo = content.Sinfo
		messageContent = content.Cinfo
	)

	fCheck, sCheck := isEmpty(senderInfo), isEmpty(messageContent)

	if !((fCheck || sCheck) && !l.faultyStatus.Load()) {
		return errors.New("Write Operations Aborted Before Completion")
	}

	fResult, fErr := l.instance.ExecContext(l.dbCtx, insertSenderStmt, senderInfo.Addr, senderInfo.Port)
	sResult, sErr := l.instance.ExecContext(l.dbCtx, insertMessageStmt, messageContent.Endpoint, messageContent.Key, messageContent.Value)

	if fErr != nil || sErr != nil {
		return errors.New("Impossible to Operate INSERT statements")
	}

	fId, _ := fResult.LastInsertId()
	sId, _ := sResult.LastInsertId()

	_, err := l.instance.ExecContext(l.dbCtx, insertBufferStmt, opType, fId, sId)
	if err != nil {
		return errors.New("Impossible to Operate final INSERT statements")
	}

	return nil
}

func (l *LogDB) DeleteMessage(searchKey string) error {
	if searchKey == "" {
		return errors.New("Delete Message Operation Aborted due to empty search key")
	}

	var err error
	l.tx, err = l.instance.BeginTx(l.dbCtx, nil)

	result, err := l.tx.ExecContext(l.dbCtx, deleteMessageStmt, searchKey)
	if err != nil {
		return errors.New("Some Errors Occured When Tried to Perform a Delete Message Operation " + err.Error())
	}

	if result == nil {
		return errors.New("Some Errors Occured When Tried to Peform a Delete Message Operation")
	}

	if msgId := l.getMessageID(searchKey); msgId == -1 {
		l.tx.Rollback()
	}

	return nil
}

func (l *LogDB) StartDB() error {
	var err error

	l.instance, err = sql.Open("sqlite", "logger.db")
	if err != nil {
		return err
	}
	
	l.instance.SetMaxOpenConns(30)
	l.instance.SetConnMaxIdleTime(2 * time.Second)
	l.instance.SetMaxIdleConns(5)

	errInt := l.setDBInternals()
	go l.pinger()

	return errInt
}

func (l *LogDB) ShutdownDB() {
	l.instance.Close()
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
			l.faultyStatus.Store(true)
		} else {
			if l.faultyStatus.Load() {
				l.faultyStatus.Store(false)
			}
		}
	}
}

func isEmpty(msg any) bool {
	var ok bool 

	switch value := msg.(type) {
	case model.Sender:
		ok = check(value.Addr, value.Port)
	case model.MessageContent:
		ok = check(value.Endpoint, value.Key, value.Value)
	}

	return ok
}

func check(content ...string) bool {
	return slices.Contains(content, "")
}
