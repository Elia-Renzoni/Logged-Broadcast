package model

type BasicMessage struct {
	Endpoint string `json:"endpoint"`
	Key      string `json:"key,omitempty"`
	Value    string `json:"value,omitempty"`
	Node     string `json:"add,omitempty"`
	Secret   string `json:"secret,omitempty"`
}

type BasicError struct {
	Error string `json:"err"`
}

type BasicPositiveAck struct {
	Succ string `json:"succ"`
}

type PersistentMessage struct {
	Sinfo Sender
	Cinfo MessageContent
}

type Sender struct {
	Addr string
	Port string
}

type MessageContent struct {
	Endpoint string
	Key      string
	Value    string
}

type BasicJoinMessage struct {
	SeedFlagOption      bool     `json:"seed"`
	NodeCompleteAddress []string `json:"address"`
}
