package messages

// BaseResponse описывает базовый ответ
type BaseResponse struct {
	Status string `json:"status"`
}

// CreateRoomResponse описывает ответ на запрос на создание комнаты
type CreateRoomResponse struct {
	BaseResponse
	Addr string `json:"addr"`
}
