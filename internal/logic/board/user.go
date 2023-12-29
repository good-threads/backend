package board

type Logic interface {
	Get(username string) (string, error)
}

type logic struct{}

func Setup() Logic {
	return &logic{}
}

func (l *logic) Get(username string) (string, error) {
	return "TODO", nil
}
