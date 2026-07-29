package testutil

import (
	"math/rand/v2"
	"strconv"
	"wal/internal/wal"
)

func GenerateTestCommand(count, variable int) []wal.Command {
	variables := make([]string, variable)
	varStat := make(map[string]bool, variable)

	for i := 0; i < len(variables); i++ {
		v := randomString(rand.IntN(13) + 3)
		variables[i] = v
		varStat[v] = false
	}
	commands := make([]wal.Command, count)
	for i := 0; i < count; i++ {
		v := variables[rand.IntN(variable)]
		created := varStat[v]
		t := randomType()
		if !created {
			t = wal.Insert
		}
		commands[i] = wal.Command{Command: t, Property: v, Value: strconv.Itoa(rand.Int())}
		switch t {
		case wal.Insert:
			varStat[v] = true
		case wal.Delete:
			varStat[v] = false
		}
	}

	return commands
}
func randomType() wal.Type {
	if rand.Int32N(100) > 90 {
		return wal.Delete
	}
	return wal.Update
}

func randomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}
