package main

import (
	"fmt"
	"math/rand/v2"
	"wal/internal/testutil"
	"wal/internal/wal"
)

func main() {
	w := wal.NewWal("D:\\tmp", "file")
	err := w.Init()
	if err != nil {
		panic("Cant init WAL")
	}

	for i := 0; i < 12_000; i++ {

		w.Write(testutil.GenerateTestCommand())
	}
	for i := 0; i < 10; i++ {
		index := uint64(rand.Int32N(12_000))
		command, err := w.Read(index)
		if err != nil {
			fmt.Println("%w", err)
		}
		fmt.Println(command.Command, command.Values)
	}
	w.Close()
}
