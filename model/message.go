package model


type BasicMessage struct {
	Endpoint string `json:"endpoint"`
	Key string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
	Node string `json:"addr",omitempty`
}
