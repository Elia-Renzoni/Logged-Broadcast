package db

import (
	"database/sql"
)

func (l *LogDB) getMessageId() {

}

func (l *LogDB) getSenderId() int {
	return 0
}


func (l *LogDB) setDBInternals() error {
	var (
		result sql.Result
		err error
	)

	result, err = l.instance.ExecContext(l.dbCtx, CREATE_DB, nil)
	if err != nil || result == nil {
		return err
	}

	result, err = l.instance.ExecContext(l.dbCtx, CREATE_TABLE_MESSAGE, nil)
	if err != nil || result == nil {
		return err
	}

	result, err = l.instance.ExecContext(l.dbCtx, CREATE_TABLE_SENDER, nil)
	if err != nil || result == nil {
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
	var messageId int

	rows, err := l.tx.QueryContext(l.dbCtx, FETCH_MESSAGE_ID, messageKey)
	if err != nil {
		return -1
	}

	for rows.Next() {
		if err := rows.Scan(&messageId); err != nil {
			return -1
		}
	}

	cerr := rows.Close()
	if cerr != nil {
		return -1
	}

	return 0
}
