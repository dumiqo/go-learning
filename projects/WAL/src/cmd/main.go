package main

import (
	"fmt"
	"math/rand/v2"
	"time"
	"wal/internal/testutil"
	"wal/internal/wal"
)

func main() {
	w := wal.NewWal("D:\\tmp", "file")
	err := w.Init()
	if err != nil {
		panic("Cant init WAL")
	}
	go w.StartPeriodicSnapshot(5 * time.Second)

	count := 1_000_000
	variables := 1_000
	commands := testutil.GenerateTestCommand(count, variables)
	for _, commad := range commands {
		w.Write(commad)
	}
	for i := 0; i < 10; i++ {
		index := uint64(rand.Int32N((int32(count))))
		command, err := w.Read(index)
		if err != nil {
			fmt.Println("%w", err)
		}
		fmt.Println(command.Command, command.Property, command.Value)
	}
	w.Close()
}
