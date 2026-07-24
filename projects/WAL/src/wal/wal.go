package wal

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"
)

type WAL struct {
	File        *os.File
	IndexOffset map[uint64]int64 // индекс → смещение в файле
	MaxIndex    uint64
	Mu          sync.RWMutex
}

func (w *WAL) Write(command Command) (uint64, error) {
	w.Mu.Lock()
	defer w.Mu.Unlock()
	entry := Entry{w.MaxIndex, command, 0}
	// log, err := entry.ToLog()

	// if err != nil {
	// 	return 0, err
	// }
	err := binary.Write(w.File, binary.BigEndian, entry)
	// _, err = fmt.Fprintln(w.File, log)
	if err != nil {
		return 0, err
	}
	err = w.File.Sync() // вызывает fsync
	if err != nil {
		return 0, err
	}
	w.MaxIndex = w.MaxIndex + 1
	return entry.Index, nil
}

func (w *WAL) Read(index uint64) (Command, error) {
	w.Mu.RLock()
	defer w.Mu.RUnlock()

	if index >= w.MaxIndex {
		return Command{}, fmt.Errorf("invalid index, max index %d", w.MaxIndex)
	}
	// w.File.Seek()
	// offset := w.index[index]

	// w.file.see
	return Command{}, nil
}
