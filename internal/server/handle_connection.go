package server

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net"

	"github.com/4aykovski/tcp-kursach/internal/messages"
	"github.com/4aykovski/tcp-kursach/internal/room"
	"github.com/4aykovski/tcp-kursach/internal/utils"
)

// handleConnection обрабатывает входящий запрос
func (s *Server) handleConnection(c net.Conn) {
	defer c.Close() // nolint:errcheck

	// считываем содержимое запроса в буффер
	buf := make([]byte, 1024)
	_, err := c.Read(buf)
	if err != nil {
		slog.Error("can't read request", slog.String("error", err.Error()), slog.String("from", c.RemoteAddr().String()))
		return
	}

	// удаляем лишние символы
	buf = bytes.TrimSpace(bytes.Trim(buf, "\x00"))

	// парсим запрос
	var req messages.BaseRequest
	err = json.Unmarshal(buf, &req)
	if err != nil {
		slog.Error("can't unmarshal request", slog.String("error", err.Error()), slog.String("from", c.RemoteAddr().String()))
		return
	}

	slog.Debug("got request", slog.Any("request", req), slog.String("from", c.RemoteAddr().String()))

	// в зависимости от типа запроса вызываем соответствующий обработчик
	switch req.Type {
	case messages.CreateRoomRequestType:
		s.handleCreateRoomRequest(c, req)

	default:
		s.handleUnknownRequest(c)

	}

	slog.Info("server request processed", slog.String("from", c.RemoteAddr().String()))
}

// handleCreateRoomRequest обрабатывает запрос на создание комнаты
func (s *Server) handleCreateRoomRequest(c net.Conn, req messages.BaseRequest) {
	slog.Debug("got create room request", slog.String("from", c.RemoteAddr().String()))

	// создаем комнату
	r := room.New(s.tcp.Addr().(*net.TCPAddr).IP.String(), req.ID)

	// сохраняем комнату в структуру сервера
	s.rooms[r.ID] = r

	// отправляем ответ
	var resp messages.CreateRoomResponse
	resp.Status = "ok"
	resp.Addr = r.Addr

	err := utils.SendResponse(c, resp)
	if err != nil {
		slog.Error("can't send response", slog.String("error", err.Error()), slog.String("from", c.RemoteAddr().String()))
		return
	}

	slog.Debug("room created", slog.Any("room-id", r.ID), slog.String("from", c.RemoteAddr().String()), slog.Any("room", s.rooms[r.ID]), slog.Any("rooms", s.rooms))
}

// handleUnknownRequest обрабатывает неизвестный запрос
func (s *Server) handleUnknownRequest(c net.Conn) {
	slog.Info("unknown request", slog.String("from", c.RemoteAddr().String()))

	// отправляем ответ
	var resp messages.BaseResponse
	resp.Status = "error: unknown request"

	err := utils.SendResponse(c, resp)
	if err != nil {
		slog.Error("can't send response", slog.String("error", err.Error()), slog.String("from", c.RemoteAddr().String()))
		return
	}
}
