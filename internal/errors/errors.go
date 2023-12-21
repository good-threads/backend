package errors

type UsernameAlreadyTaken struct{}

func (e *UsernameAlreadyTaken) Error() string { return "" }

type BadPassword struct{}

func (e *BadPassword) Error() string { return "" }

type BadUsername struct{}

func (e *BadUsername) Error() string { return "" }
