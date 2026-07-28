package wal

import (
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
)

type Entry struct {
	Index   uint64
	Data    []byte
	Command Command
	CRC32   uint32
}

func NewEntry(index uint64, command Command) (*Entry, error) {
	data, err := command.Encode()
	if err != nil {
		return nil, err
	}
	crc := crc32.ChecksumIEEE(data)
	return &Entry{index, data, command, crc}, nil
}

func (e *Entry) ToLog() (string, error) {
	return fmt.Sprintf("%d|%s|%d", e.Index, e.Data, e.CRC32), nil
}

func FromLog(line string) (*Entry, error) {
	parts := strings.Split(line, "|")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid log format")
	}

	index, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return nil, err
	}

	crc, err := strconv.ParseUint(parts[2], 10, 32)
	if err != nil {
		return nil, err
	}

	cmd, err := Decode(parts[1])
	if err != nil {
		return nil, err
	}

	data := []byte(parts[1])
	expectedCRC := crc32.ChecksumIEEE(data)

	if uint32(crc) != expectedCRC {
		return nil, fmt.Errorf("crc mismatch")
	}

	entry, err := NewEntry(index, cmd)
	if err != nil {
		return nil, err
	}
	return entry, nil
}
