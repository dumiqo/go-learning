package testutil

import (
	"crypto/rand"
	"wal/internal/wal"
)

func GenerateTestCommand() wal.Command {

	rnd := rand.Text()
	return wal.Command{"select * from " + rnd, []string{rand.Text(), rand.Text()}}
}
