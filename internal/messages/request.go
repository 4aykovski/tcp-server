package messages

import "github.com/google/uuid"

// RequestType описывает возможные типы запроса
type RequestType string

const (
	// CreateRoomRequestType описывает запрос на создание комнаты
	CreateRoomRequestType RequestType = "create"

	// ConnectRoomRequestType описывает запрос на подключение к комнате
	ConnectRoomRequestType RequestType = "connect"

	// DisconnectRoomRequestType описывает запрос на отключение от комнаты
	DisconnectRoomRequestType RequestType = "disconnect"

	// SendMessageRequestType описывает запрос на отправку сообщения
	SendMessageRequestType RequestType = "send"
)

// BaseRequest описывает базовый запрос, который используется в каждом другом запросе
type BaseRequest struct {
	Type RequestType `json:"type"`
	ID   uuid.UUID   `json:"id"`
}

// ConnectRoomRequest описывает запрос на подключение к комнате
type ConnectRoomRequest struct {
	BaseRequest
	Addr string `json:"addr"`
}

// DisconnectRoomRequest описывает запрос на отключение от комнаты
type DisconnectRoomRequest struct {
	BaseRequest
}

// SendMessageRequest описывает запрос на отправку сообщения
type SendMessageRequest struct {
	BaseRequest
	Addr    string `json:"addr"`
	Message string `json:"message"`
	Date    string `json:"date"`
}
