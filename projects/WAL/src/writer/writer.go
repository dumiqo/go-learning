package writer

type Entry struct {
	Command string
	Values  []string
}

func Write(entry Entry) (int, error) {
	return 0, nil
}

func Read(index int) (Entry, error) {

	return Entry{}, nil
}
