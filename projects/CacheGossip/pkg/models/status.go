package models

type Status int

const (
	Alive   Status = iota // 0 - узел работает
	Suspect               // 1 - узел под подозрением
	Dead                  // 2 - узел мертв
)

func (s Status) String() string {
	switch s {
	case Alive:
		return "alive"
	case Suspect:
		return "suspect"
	case Dead:
		return "dead"
	default:
		return "unknown"
	}
}
