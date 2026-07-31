package wal

type Snapshot struct {
	Index    uint64
	State    map[string]string
	FilePath string
}
