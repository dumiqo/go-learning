package wal

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type Positions struct {
	Offset int64
	Count  int
}

type WAL struct {
	File        *os.File
	IndexOffset map[uint64]Positions // индекс → смещение в файле
	MaxIndex    uint64
	Mu          sync.RWMutex
}

func (w *WAL) Write(command Command) (uint64, error) {
	w.Mu.Lock()
	defer w.Mu.Unlock()
	entry := Entry{w.MaxIndex, command, 0}
	log, err := entry.ToLog()
	if err != nil {
		return 0, err
	}
	offset, err := w.File.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	n, err := fmt.Fprintln(w.File, log)
	if err != nil {
		return 0, err
	}
	err = w.File.Sync() // медленно, лучше делать раз в n мс
	if err != nil {
		return 0, err
	}
	w.IndexOffset[w.MaxIndex] = Positions{offset, n}
	w.MaxIndex = w.MaxIndex + 1
	return entry.Index, nil
}

func (w *WAL) Read(index uint64) (Command, error) {
	w.Mu.RLock()
	defer w.Mu.RUnlock()

	if index >= w.MaxIndex {
		return Command{}, fmt.Errorf("invalid index, max index %d", w.MaxIndex)
	}
	position := w.IndexOffset[index]
	buf := make([]byte, position.Count)

	_, err := w.File.ReadAt(buf, position.Offset)

	if err != nil {
		return Command{}, err
	}

	entry, err := EntryFromLog(strings.TrimSpace(string(buf)))

	return entry.Command, nil
}
