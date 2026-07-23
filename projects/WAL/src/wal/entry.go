package wal

type Entry struct {
	Command string
	Values  []string
}
