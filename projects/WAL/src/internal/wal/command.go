package wal

import (
	"fmt"
	"strings"
)

type Type string

const (
	Insert Type = "Insert"
	Update Type = "Update"
	Delete Type = "Delete"
)

type Command struct {
	Command  Type
	Property string
	Value    string
}

func (c *Command) Encode() ([]byte, error) {
	value := fmt.Sprintf("%s;%s;%s",
		c.Command,
		c.Property,
		c.Value,
	)
	return []byte(value), nil
}

func Decode(str string) (Command, error) {
	cmdParts := strings.SplitN(str, ";", 3)
	if len(cmdParts) != 3 {
		return Command{}, fmt.Errorf("invalid command")
	}
	t, err := parseStatus(cmdParts[0])
	if err != nil {
		return Command{}, err
	}
	return Command{
		Command:  t,
		Property: cmdParts[1],
		Value:    cmdParts[2],
	}, nil
}
func parseStatus(s string) (Type, error) {
	switch s {
	case "Insert", "insert":
		return Insert, nil
	case "Delete", "delete":
		return Delete, nil
	case "Update", "update":
		return Update, nil
	default:
		return Delete, fmt.Errorf("invalid status: %s", s)
	}
}
