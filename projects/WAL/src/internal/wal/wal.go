package wal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type position struct {
	offset int64
	count  int
	file   string
}

type WAL struct {
	index map[uint64]position // индекс → смещение в файле

	mu sync.RWMutex

	direction, file string
	writer, reader  *os.File

	maxSize       int64
	lastIdex      uint64
	snapshotIndex uint64
}

func NewWal(dirPath, fileName string) *WAL {
	return &WAL{make(map[uint64]position), sync.RWMutex{}, dirPath, fileName, nil, nil, 1024 * 1024, 0, 0}
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

func (w *WAL) StartPeriodicSnapshot(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop() // важно освободить ресурсы

	for range ticker.C {
		if err := w.saveSnapshot(); err != nil {
			log.Printf("snapshot creation failed: %v", err)
		}
	}
}

func (w *WAL) Close() error {
	// получается что мы можем получить ошибку при закрытие writer и не попытаться закрыть reader
	if err := w.writer.Close(); err != nil {
		return err
	}
	return w.reader.Close()
}

func (w *WAL) Write(command Command) (uint64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.needRotate() {
		err := w.rotate()
		if err != nil {
			return 0, err
		}
	}
	entry, err := NewEntry(w.lastIdex+1, command)
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
	stat, _ := os.Stat(w.writer.Name())
	w.index[entry.Index] = position{offset, n, stat.Name()}
	w.lastIdex++
	return entry.Index, nil
}

func (w *WAL) Read(index uint64) (Command, error) {
	if index >= w.lastIdex {
		return Command{}, fmt.Errorf("invalid index, max index %d", w.lastIdex)
	}
	position := w.index[index]
	buf := make([]byte, position.count)
	stat, _ := os.Stat(w.reader.Name())
	if stat.Name() != position.file {
		filePath := filepath.Join(w.direction, w.file)
		reader, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
		if err != nil {
			return Command{}, errors.Wrap(err, "error in reading")
		}
		w.reader = reader
	}
	_, err := w.reader.ReadAt(buf, position.offset)

	if err != nil {
		return Command{}, errors.Wrap(err, "error in reading")
	}

	entry, err := FromLog(strings.TrimSpace(string(buf)))
	if err != nil {
		return Command{}, errors.Wrap(err, "error in reading")
	}
	return entry.Command, nil
}

func (w *WAL) saveSnapshot() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	path := filepath.Join(w.direction, "snapshot")

	state := make(map[string]string)

	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		f, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error in reading snapshot, %w", err)
		}
		err = json.Unmarshal(f, &state)
		if err != nil {
			return fmt.Errorf("error in unmarshal snapshot, %w", err)
		}
	}

	for i := w.snapshotIndex + 1; i <= w.lastIdex; i++ {
		entity, _ := w.Read(i)

		switch entity.Command {
		case Insert:
			state[entity.Property] = entity.Value
		case Delete:
			delete(state, entity.Property)
		case Update:
			state[entity.Property] = entity.Value
		}
	}
	j, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("error in serializing, %w", err)
	}
	err = os.WriteFile(path, j, 0644)
	if err != nil {
		return fmt.Errorf("error in saving snapshot file, %w", err)
	}
	w.snapshotIndex = w.lastIdex
	return nil
}

func (w *WAL) rotate() error {
	stat, _ := os.Stat(w.writer.Name())

	newFileName := newFileName(stat.Name(), w.file)
	filePath := filepath.Join(w.direction, newFileName)
	writer, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return errors.Wrap(err, "cant rotate")
	}
	err = w.writer.Sync() // медленно, лучше делать раз в n мс
	if err != nil {
		return errors.Wrap(err, "cant rotate")
	}
	w.writer = writer
	return nil
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
