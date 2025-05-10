package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
)

func SendResponse(c net.Conn, resp any) error {
	slog.Debug("sending response", slog.Any("resp", resp))

	err := json.NewEncoder(c).Encode(resp)
	if err != nil {
		return fmt.Errorf("can't encode response: %w", err)
	}

	err = c.Close()
	if err != nil {
		return fmt.Errorf("can't close connection: %w", err)
	}

	return nil
}

func SendRequest(addr net.TCPAddr, req any) error {
	conn, err := net.DialTCP("tcp", nil, &addr)
	if err != nil {
		return fmt.Errorf("can't dial tcp: %w", err)
	}

	err = json.NewEncoder(conn).Encode(req)
	if err != nil {
		return fmt.Errorf("can't encode request: %w", err)
	}

	return nil
}
