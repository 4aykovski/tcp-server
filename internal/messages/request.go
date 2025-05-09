package messages

import "github.com/google/uuid"

type RequestType string

const (
	CreateRoomRequestType     RequestType = "create"
	ConnectRoomRequestType    RequestType = "connect"
	DisconnectRoomRequestType RequestType = "disconnect"
	SendMessageRequestType    RequestType = "send"
)

type BaseRequest struct {
	Type RequestType `json:"type"`
	ID   uuid.UUID   `json:"id"`
}

type ConnectRoomRequest struct {
	BaseRequest
	Addr string `json:"addr"`
}

type DisconnectRoomRequest struct {
	BaseRequest
}

type SendMessageRequest struct {
	BaseRequest
	Addr    string `json:"addr"`
	Message string `json:"message"`
	Date    string `json:"date"`
}
