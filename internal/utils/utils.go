package utils

import "bufio"

// Чтение до определенного разделителя (например, новой строки)
func ReadUntilDelimiter(reader *bufio.Reader) ([]byte, error) {
	return reader.ReadBytes('\n')
}
