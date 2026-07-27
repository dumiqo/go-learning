package wal

import (
	"errors"
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
	indexOffset    map[uint64]Positions // индекс → смещение в файле
	maxIndex       uint64
	mu             sync.RWMutex
	writer, reader *os.File
}

func NewWal(dirPath, fileName string) (*WAL, error) {
	filePath := dirPath + "\\" + fileName
	writer, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, errors.New("Invalid file path")
	}
	reader, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, errors.New("Invalid file path")
	}

	return &WAL{make(map[uint64]Positions), 0, sync.RWMutex{}, writer, reader}, nil
}

func (w *WAL) Write(command Command) (uint64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	entry, err := NewEntry(w.maxIndex, command)
	if err != nil {
		return 0, err
	}
	log, err := entry.ToLog()
	if err != nil {
		return 0, err
	}
	offset, err := w.writer.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	n, err := fmt.Fprintln(w.writer, log)
	if err != nil {
		return 0, err
	}
	err = w.writer.Sync() // медленно, лучше делать раз в n мс
	if err != nil {
		return 0, err
	}
	w.indexOffset[w.maxIndex] = Positions{offset, n}
	w.maxIndex = w.maxIndex + 1
	return entry.Index, nil
}

func (w *WAL) Read(index uint64) (Command, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if index >= w.maxIndex {
		return Command{}, fmt.Errorf("invalid index, max index %d", w.maxIndex)
	}
	position := w.indexOffset[index]
	buf := make([]byte, position.Count)

	_, err := w.reader.ReadAt(buf, position.Offset)

	if err != nil {
		return Command{}, err
	}

	entry, err := FromLog(strings.TrimSpace(string(buf)))
	if err != nil {
		return Command{}, err
	}
	return entry.Command, nil
}
