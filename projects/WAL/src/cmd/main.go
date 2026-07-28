package main

import (
	"fmt"
	"wal/internal/testutil"
	"wal/internal/wal"
)

func main() {
	w := wal.NewWal("D:\\tmp", "file")
	err := w.Init()
	if err != nil {
		panic("Cant init WAL")
	}

	index := uint64(0)
	for i := 0; i < 500_000; i++ {

		index, err = w.Write(testutil.GenerateTestCommand())
		if err != nil {
			panic("Cant write log")
		}
	}
	command, _ := w.Read(index)
	fmt.Println(command.Command, command.Values)
	w.Close()
}
