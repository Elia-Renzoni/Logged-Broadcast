package db

import (
	"database/sql"
)


func (l *LogDB) getSenderId() int {
	return 0
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
	if err != nil {
		return err
	}

	return nil
}

// this method is responsible of
// deleting the content from
// the Buffer Table
func (l *LogDB) deleteFromBuffer() {

}

func (l *LogDB) getMessageID(messageKey string) int {
	rows, err := l.tx.QueryContext(l.dbCtx, fetchMessageIDStmt, messageKey)
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
