package logic

type test interface {
	Test() string
}

type Test struct{}

func NewTest() *Test {
	return &Test{}
}

func (h *Test) Handler() string {
	return "welcome\n"
}
