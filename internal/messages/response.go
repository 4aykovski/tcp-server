package messages

type BaseResponse struct {
	Status string `json:"status"`
}

type CreateRoomResponse struct {
	BaseResponse
	Addr string `json:"addr"`
}
