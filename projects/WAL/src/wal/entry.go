package wal

import (
	"fmt"
	"hash/crc32"
)

type Command struct {
	Command string
	Values  []string
}

type Entry struct {
	Index   uint64
	Command Command
	CRC32   uint32
}

func (c *Command) Encode() ([]byte, uint32, error) {
	value := fmt.Sprintf("%s;%s", c.Command, c.Values)
	return []byte(value), crc32.ChecksumIEEE([]byte(value)), nil
}

func (e *Entry) ToLog() (string, error) {
	data, crc, err := e.Command.Encode()
	e.CRC32 = crc
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d|%s|%d", e.Index, data, e.CRC32), nil
}
