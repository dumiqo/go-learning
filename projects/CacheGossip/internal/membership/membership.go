package membership

type Membership struct {
	Members map[string]Member
}

func NewMembership(hosts map[string]string) *Membership {
	return &Membership{}
}

type Member struct {
	Url    string
	Port   int
	Status Status
}

type Status int

const (
	Up      Status = 0
	Down    Status = 1
	Pending Status = 2
)
