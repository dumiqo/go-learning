package wal

import (
	"fmt"
	"strings"
)

type Command struct {
	Command string
	Values  []string
}

func (c *Command) Encode() ([]byte, error) {
	value := fmt.Sprintf("%s;%s",
		c.Command,
		strings.Join(c.Values, ","),
	)
	return []byte(value), nil
}

func Decode(str string) (Command, error) {

	cmdParts := strings.SplitN(str, ";", 2)
	if len(cmdParts) != 2 {
		return Command{}, fmt.Errorf("invalid command")
	}
	return Command{
		Command: cmdParts[0],
		Values:  strings.Split(cmdParts[1], ","),
	}, nil
}
