package wal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Position struct {
	Offset int64
	Count  int
}

type WAL struct {
	index map[uint64]Position // индекс → смещение в файле

	mu sync.RWMutex

	direction, file string
	writer, reader  *os.File

	maxSize int64
}

func NewWal(dirPath, fileName string) *WAL {
	return &WAL{make(map[uint64]Position), sync.RWMutex{}, dirPath, fileName, nil, nil, 1024 * 100}
}

func (w *WAL) Init() error {
	filePath := filepath.Join(w.direction, w.file)
	writer, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("open writer: %w", err)
	}
	reader, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("open reader: %w", err)
	}
	w.writer = writer
	w.reader = reader
	return nil
}

func (w *WAL) Close() error {
	if err := w.writer.Close(); err != nil {
		return err
	}
	return w.reader.Close()
}

func (w *WAL) Write(command Command) (uint64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.needRotate() {
		w.rotate()
	}
	entry, err := NewEntry(uint64(len(w.index)), command)
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
	w.index[uint64(len(w.index))] = Position{offset, n}
	return entry.Index, nil
}

func (w *WAL) Read(index uint64) (Command, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if index >= uint64(len(w.index)) {
		return Command{}, fmt.Errorf("invalid index, max index %d", uint64(len(w.index)))
	}
	position := w.index[index]
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

func (w *WAL) rotate() {
	stat, _ := os.Stat(w.writer.Name())

	newFileName := newFileName(stat.Name(), w.file)
	filePath := filepath.Join(w.direction, newFileName)
	writer, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic("Cant rotate")
	}
	w.writer = writer
}

func newFileName(currentName, startName string) string {
	if currentName == startName {
		return startName + "001"
	}

	if len(currentName) > len(startName) {
		numberPart := currentName[len(startName):]
		number, err := strconv.Atoi(numberPart)
		if err != nil {
			return startName + "001"
		}
		return startName + fmt.Sprintf("%03d", number+1)
	}

	return startName + "001"
}

func (w *WAL) needRotate() bool {
	name := w.writer.Name()
	stat, _ := os.Stat(name)
	size := stat.Size()
	return w.maxSize <= size
}
