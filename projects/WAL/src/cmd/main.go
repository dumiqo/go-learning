package main

import (
	"crypto/rand"
	"fmt"
	"wal/wal"
)

func main() {
	w := wal.NewWal("D:\\tmp", "file")
	err := w.Init()
	if err != nil {
		panic("Cant init WAL")
	}

	index := uint64(0)
	for i := 0; i < 100; i++ {

		rnd := rand.Text()
		index, err = w.Write(wal.Command{"select * from " + rnd, []string{rand.Text()}})
		if err != nil {
			panic("Cant write log")
		}
	}
	entry, _ := w.Read(index)
	fmt.Println(entry.Command)
	w.Close()
}
