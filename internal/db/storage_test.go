package db_test


import (
	"testing"
	"log-b/internal/db"
	"log-b/model"
)

func TestWriteMessage(t *testing.T) {
	instance := db.NewDB()
	err := instance.StartDB()
	defer instance.ShutdownDB()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	msg := model.PersistentMessage{
		Sinfo: model.Sender{
			Addr: "127.0.0.1",
			Port: "8080",
		},
		Cinfo: model.MessageContent{
			Endpoint: "/addbk",
			Key:      "foo",
			Value:    "bar",
		},
	}
	wErr := instance.WriteMessage(msg, 0)
	if wErr != nil {
		t.Fatalf("%s", wErr.Error())
	}
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
