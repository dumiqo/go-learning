package wal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	return &WAL{make(map[uint64]Position), sync.RWMutex{}, dirPath, fileName, nil, nil, 1024} //1 kb
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
	// w.mu.RLock()
	// defer w.mu.RUnlock()

	name := w.reader.Name()
	stat, _ := os.Stat(name)
	//у нас есть имя текущего файла и име исходного файла, нужно понять какой индекс выставить и создать новый файл
	//ридер не нужно создавать, так как приходится просчитывать все файлы для поиска нужного
	f := stat.Name()
	size := stat.Size()
	w.file = f + string(size)
	w.direction = w.file
	// w.reader.ReadDir()
}

func (w *WAL) needRotate() bool {
	name := w.reader.Name()
	stat, _ := os.Stat(name)
	return w.maxSize <= stat.Size()
}
