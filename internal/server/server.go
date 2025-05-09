package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/4aykovski/tcp-kursach/internal/config"
	"github.com/4aykovski/tcp-kursach/internal/room"
)

// Server структура, описывающая tcp сервер
type Server struct {
	// tcp сам tcp слушатель, который принимает все запросы по определенному порта
	tcp net.Listener

	// rooms список созданных в текущий момент комнат в виде словаря, где ключ - порт комнаты, а значение - сама комната
	rooms map[int]room.Room
}

func New(ctx context.Context) *Server {
	return &Server{
		rooms: make(map[int]room.Room, 10),
	}
}

// Run запускает tcp сервер
func (s *Server) Run(ctx context.Context, cfg *config.C) error {

	// создаем tcp слушатель на порту из конфига
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.ParseIP(cfg.Host),
		Port: cfg.Port,
	})
	if err != nil {
		return fmt.Errorf("can't start tcp server: %w", err)
	}

	// сохраняем слушатель в структуру сервера
	s.tcp = listener

	// запускаем слушатель в отдельном потоке
	go s.listen()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.tcp.Close()
}

// listen слушает входящие запросы
func (s *Server) listen() {

	// запускаем бесконечный цикл прослушивания
	for {
		// ждем подключений к слушателю
		conn, err := s.tcp.Accept()
		if err != nil {
			slog.Error("can't accept connection", slog.String("error", err.Error()))
			return
		}

		// обрабатываем каждый запрос в отдельном потоке, обеспечивая конкуретную обработку запросов
		go s.handleConnection(conn)
	}
}
