package room

import (
	"log/slog"
	"net"
	"sync"

	"github.com/google/uuid"
)

type Message struct {
	From string
	Text string
	Date string
}

type Room struct {
	// ID комнаты (порт)
	ID int

	// Адрес комнаты (host:port)
	Addr string

	// UUID создателя
	Creator uuid.UUID

	// Список пользователей
	users map[uuid.UUID]net.TCPAddr

	// История сообщений
	history []Message

	tcp net.Listener
	mu  *sync.RWMutex
}

// New создает новую комнату
func New(host string, creatorID uuid.UUID) Room {
	// создаем слушатель
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP: net.ParseIP(host),
	})
	if err != nil {
		slog.Error("can't listen", slog.String("error", err.Error()))
		return Room{}
	}

	// сохраняем порт в структуру
	id := listener.Addr().(*net.TCPAddr).Port

	// сохраняем адрес в структуру
	addr := listener.Addr().String()

	// создаем комнату
	room := Room{
		ID:   id,
		Addr: addr,

		Creator: creatorID,

		tcp: listener,

		users:   make(map[uuid.UUID]net.TCPAddr),
		history: make([]Message, 0),
		mu:      &sync.RWMutex{},
	}

	// запускаем слушатель в отдельном потоке
	go room.listen()

	return room
}

// listen слушает входящие запросы
func (r *Room) listen() {

	// запускаем бесконечный цикл прослушивания
	for {
		// ждем подключений к слушателю
		conn, err := r.tcp.Accept()
		if err != nil {
			slog.Error("can't accept connection", slog.String("error", err.Error()))
			return
		}

		// обрабатываем каждый запрос в отдельном потоке, обеспечивая конкуретную обработку запросов
		go r.handleConnection(conn)
	}
}
