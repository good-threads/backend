package common

type Logic interface {
	Ping() string
}

type logic struct{}

func New() Logic {
	return &logic{}
}

func (l *logic) Ping() string {
	return "pong\n"
}
