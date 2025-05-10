package room

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net"
	"strconv"
	"strings"

	"github.com/4aykovski/tcp-kursach/internal/messages"
	"github.com/4aykovski/tcp-kursach/internal/utils"
)

// handleConnection обрабатывает входящий запрос
func (r *Room) handleConnection(c net.Conn) {
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

	slog.Debug("check if client in room", slog.String("from", c.RemoteAddr().String()), slog.Any("room clients", r.users))

	// проверяем есть ли клиент в комнате
	r.mu.Lock()
	_, ok := r.users[req.ID]
	r.mu.Unlock()
	if !ok && req.Type != messages.ConnectRoomRequestType {
		slog.Error("client not found", slog.String("from", c.RemoteAddr().String()))
		return
	}

	slog.Debug("client in room", slog.String("from", c.RemoteAddr().String()), slog.Any("room clients", r.users))

	// в зависимости от типа запроса вызываем соответствующий обработчик
	switch req.Type {
	case messages.ConnectRoomRequestType:
		r.handleConnectRoomRequest(c, buf)

	case messages.DisconnectRoomRequestType:
		r.handleDisconnectRoomRequest(c, req)

	case messages.SendMessageRequestType:
		r.handleSendMessageRequest(c, buf)

	default:
		r.handleUnknownRequest(c)

	}

	slog.Info("room request processed", slog.String("from", c.RemoteAddr().String()))
}

// handleUnknownRequest обрабатывает неизвестный запрос
func (r *Room) handleUnknownRequest(c net.Conn) {
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

func (r *Room) handleConnectRoomRequest(c net.Conn, buf []byte) {
	slog.Debug("got connect room request", slog.String("from", c.RemoteAddr().String()))

	// парсим запрос
	var req messages.ConnectRoomRequest
	err := json.Unmarshal(buf, &req)
	if err != nil {
		slog.Error("can't unmarshal request", slog.String("error", err.Error()), slog.String("from", c.RemoteAddr().String()))
		return
	}

	if _, ok := r.users[req.ID]; ok {
		slog.Info("client already in room", slog.String("from", c.RemoteAddr().String()))
		return
	}

	slog.Debug("client not in room, adding")

	// сохраняем пользователя
	portStr := strings.Split(req.Addr, ":")[1]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		slog.Error("can't parse port", slog.String("error", err.Error()), slog.String("from", c.RemoteAddr().String()))
		return
	}

	r.mu.Lock()
	r.users[req.ID] = net.TCPAddr{IP: net.ParseIP(strings.Split(req.Addr, ":")[0]), Port: port}
	r.mu.Unlock()

	// отправляем ответ
	var resp messages.BaseResponse
	resp.Status = "ok"

	err = utils.SendResponse(c, resp)
	if err != nil {
		slog.Error("can't send response", slog.String("error", err.Error()), slog.String("from", c.RemoteAddr().String()))
		return
	}
}

func (r *Room) handleDisconnectRoomRequest(c net.Conn, req messages.BaseRequest) {
	slog.Debug("got disconnect room request", slog.String("from", c.RemoteAddr().String()))

	r.mu.Lock()
	delete(r.users, req.ID)
	r.mu.Unlock()

	// отправляем ответ
	var resp messages.BaseResponse
	resp.Status = "ok"

	err := utils.SendResponse(c, resp)
	if err != nil {
		slog.Error("can't send response", slog.String("error", err.Error()), slog.String("from", c.RemoteAddr().String()))
		return
	}

	slog.Debug("check if room empty")

	if len(r.users) == 0 {
		slog.Debug("room empty, closing")

		err = r.tcp.Close()
		if err != nil {
			slog.Error("can't close tcp listener", slog.String("error", err.Error()))
			return
		}
	}
}

func (r *Room) handleSendMessageRequest(c net.Conn, buf []byte) {
	slog.Debug("got send message request", slog.String("from", c.RemoteAddr().String()))

	// парсим запрос
	var req messages.SendMessageRequest
	err := json.Unmarshal(buf, &req)
	if err != nil {
		slog.Error("can't unmarshal request", slog.String("error", err.Error()), slog.String("from", c.RemoteAddr().String()))
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	addr := r.users[req.ID].IP.String() + ":" + strconv.Itoa(r.users[req.ID].Port)
	req.Addr = addr

	slog.Debug("sending message to clients in room", slog.String("from", addr), slog.Any("message", req.Message))

	// отправляем сообщение всем пользователям
	for user, conn := range r.users {
		if user == req.ID {
			continue
		}

		err = utils.SendRequest(conn, req)
		if err != nil {
			slog.Error("can't send responce to user", slog.String("error", err.Error()), slog.String("from", addr), slog.String("to", c.RemoteAddr().String()))
			delete(r.users, user)
			continue
		}
	}

	slog.Debug("message sent", slog.String("from", addr), slog.Any("message", req.Message))

	// сохраняем сообщение в историю
	msg := Message{
		From: req.Addr,
		Text: req.Message,
		Date: req.Date,
	}

	r.history = append(r.history, msg)

	// отправляем ответ
	var resp messages.BaseResponse
	resp.Status = "ok"

	err = utils.SendResponse(c, resp)
	if err != nil {
		slog.Error("can't send response", slog.String("error", err.Error()), slog.String("from", c.RemoteAddr().String()))
		return
	}

	slog.Debug("checking history", slog.Any("history", r.history))
}
