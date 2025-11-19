package db

import (
	"database/sql"
	"errors"
)

func (l *LogDB) getSenderId(tx *sql.Tx, messageId int) int {
	rows, err := tx.QueryContext(l.dbCtx, fetchSenderIDStmt, messageId)
	if err != nil {
		return -1
	}

	var senderID int

	for rows.Next() {
		if scanErr := rows.Scan(&senderID); scanErr != nil {
			return -1
		}
	}

	if cErr := rows.Close(); cErr != nil {
		return -1
	}

	return senderID
}

func (l *LogDB) setDBInternals() error {
	var (
		result sql.Result
		err error
	)

	result, err = l.instance.ExecContext(l.dbCtx, createTableMessage, nil)
	if err != nil || result == nil {
		return err
	}

	result, err = l.instance.ExecContext(l.dbCtx, createTableSender, nil)
	if err != nil || result == nil {
		return err
	}

	result, err = l.instance.ExecContext(l.dbCtx, createTableBuffer, nil)
	if err != nil || result == nil {
		return err
	}

	return nil
}

func (l *LogDB) deleteFromBuffer(tx *sql.Tx, senderID, messageID int) error {
	result, err := tx.ExecContext(l.dbCtx, deleteEntriesFromBufferStmt, messageID, senderID)
	if err != nil || result == nil {
		return err
	}

	return nil
}

func (l *LogDB) deleteSender(tx *sql.Tx, senderID int) error {
	result, err := tx.ExecContext(l.dbCtx, deleteSenderStmt, senderID)
	if err != nil { 
		return err
	}

	if result == nil {
		return errors.New("Unable to execute the SQL statement")
	}

	return nil
}

func (l *LogDB) getMessageID(tx *sql.Tx, messageKey string) int {
	rows, err := tx.QueryContext(l.dbCtx, fetchMessageIDStmt, messageKey)
	if err != nil {
		return -1
	}

	
	var messageId int
	for rows.Next() {
		if err := rows.Scan(&messageId); err != nil {
			return -1
		}
	}

	cerr := rows.Close()
	if cerr != nil {
		return -1
	}

	return messageId
}
