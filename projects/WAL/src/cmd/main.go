package main

import (
	"crypto/rand"
	"fmt"
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

		rnd := rand.Text()
		index, err = w.Write(wal.Command{"select * from " + rnd, []string{rand.Text(), rand.Text()}})
		if err != nil {
			panic("Cant write log")
		}
	}
	command, _ := w.Read(index)
	fmt.Println(command.Command, command.Values)
	w.Close()
}
