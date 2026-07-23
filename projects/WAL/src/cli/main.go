package main

import (
	"fmt"
	"wal/src/wal"
)

func main() {
	index, _ := wal.Write(wal.Entry{"select * from tmp", nil})
	entry, _ := wal.Read(index)
	fmt.Println(entry.Command)
}
