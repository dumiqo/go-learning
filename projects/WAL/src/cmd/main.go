package main

import (
	"fmt"
	"wal/wal"
)

func main() {
	log, err := wal.NewWal("D:\\tmp", "file")
	if err != nil {
		panic("Cant start WAL")
	}
	index, _ := log.Write(wal.Command{"select * from dtp", nil})
	_, _ = log.Write(wal.Command{"select * from tmp", nil})
	index, _ = log.Write(wal.Command{"select * from zxtmp", nil})
	index, _ = log.Write(wal.Command{"select * from tm1p", nil})
	_, _ = log.Write(wal.Command{"select * from fd", nil})
	_, _ = log.Write(wal.Command{"select * from zzz", nil})
	entry, _ := log.Read(index)
	fmt.Println(entry.Command)
}
