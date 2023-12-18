package welcome

type Logic interface {
	Behavior() string
}

type logic struct{}

func New() Logic {
	return &logic{}
}

func (l *logic) Behavior() string {
	return "welcome\n"
}
