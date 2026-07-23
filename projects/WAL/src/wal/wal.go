package wal

import "os"

var (
	path = "D:\\tmp\\file"
)

func Write(entry Entry) (int, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return -1, err
	}
	_, err = file.WriteString(entry.Command)
	if err != nil {
		return -1, err
	}
	err = file.Sync() // вызывает fsync
	if err != nil {
		return -1, err
	}
	return 0, nil
}

func Read(index int) (*Entry, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	return &Entry{}, nil
}
