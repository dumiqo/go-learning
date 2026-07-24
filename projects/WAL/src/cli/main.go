package main

import (
	"fmt"
	"os"
	"sync"
	"wal/wal"
)

func main() {
	file, err := os.OpenFile("D:\\tmp\\file", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error")
		return
	}
	log := wal.WAL{file, make(map[uint64]wal.Positions), 0, sync.RWMutex{}}
	index, _ := log.Write(wal.Command{"select * from dtp", nil})
	index, _ = log.Write(wal.Command{"select * from tmp", nil})
	index, _ = log.Write(wal.Command{"select * from zxtmp", nil})
	index, _ = log.Write(wal.Command{"select * from tmp", nil})
	index, _ = log.Write(wal.Command{"select * from fd", nil})
	index, _ = log.Write(wal.Command{"select * from zzz", nil})
	entry, _ := log.Read(index)
	fmt.Println(entry.Command)
}
