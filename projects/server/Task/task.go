package task

type Task struct {
	Title       string   `json:"Title"`
	Description string   `json:"Description"`
	Priority    Priority `json:"Priority"`
	Status      Status   `json:"Status"`
}

type Priority string

const (
	Low      Priority = "Low"
	Normal   Priority = "Normal"
	High     Priority = "High"
	Critical Priority = "Critical"
)

type Status string

const (
	New        Status = "New"
	Done       Status = "Done"
	InProgress Status = "InProgress"
)
