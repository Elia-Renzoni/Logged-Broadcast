package db

type Storage interface {
	WriteMessage()
	RetrieveMessage()
}

type LogDB struct {

}
