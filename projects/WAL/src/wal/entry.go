package wal

import (
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
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
	value := fmt.Sprintf("%s;%s",
		c.Command,
		strings.Join(c.Values, ","),
	)

	crc := crc32.ChecksumIEEE([]byte(value))
	return []byte(value), crc, nil
}

func (e *Entry) ToLog() (string, error) {
	data, crc, err := e.Command.Encode()
	e.CRC32 = crc
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d|%s|%d", e.Index, data, e.CRC32), nil
}
func EntryFromLog(line string) (Entry, error) {
	parts := strings.Split(line, "|")
	if len(parts) != 3 {
		return Entry{}, fmt.Errorf("invalid log format")
	}

	index, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return Entry{}, err
	}

	crc, err := strconv.ParseUint(parts[2], 10, 32)
	if err != nil {
		return Entry{}, err
	}

	cmdParts := strings.SplitN(parts[1], ";", 2)
	if len(cmdParts) != 2 {
		return Entry{}, fmt.Errorf("invalid command")
	}

	cmd := Command{
		Command: cmdParts[0],
		Values:  strings.Split(cmdParts[1], ","),
	}

	data := []byte(parts[1])
	expectedCRC := crc32.ChecksumIEEE(data)

	if uint32(crc) != expectedCRC {
		return Entry{}, fmt.Errorf("crc mismatch")
	}

	return Entry{
		Index:   index,
		Command: cmd,
		CRC32:   uint32(crc),
	}, nil
}
