package models

import "time"

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

func NewStatus(t time.Time) Status {
	elapsed := time.Now().UTC().Sub(t.UTC())
	if elapsed.Seconds() <= 20 {
		return Alive
	}
	if elapsed.Minutes() <= 1 {
		return Suspect
	}
	return Dead
}
