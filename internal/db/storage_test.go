package db_test


import (
	"testing"
	"log-b/internal/db"
)

// TODO
func TestWriteMessage(t *testing.T) {
	instance := db.NewDB()
	instance.WriteMessage()
}

// TODO
func TestDeleteMessage(t *testing.T) {
	instance := db.NewDB()
	instance.DeleteMessage()

}

// TODO
func TestDBFaultRecovery(t *testing.T) {

}

// TODO
func TestChangeStatus(t *testing.T) {

}
