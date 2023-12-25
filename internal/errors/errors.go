package errors

type UsernameAlreadyTaken struct{}

func (e *UsernameAlreadyTaken) Error() string { return "" }

type BadPassword struct{}

func (e *BadPassword) Error() string { return "" }

type BadUsername struct{}

func (e *BadUsername) Error() string { return "" }

type WrongCredentials struct{}

func (e *WrongCredentials) Error() string { return "" }

type UserNotFound struct{}

func (e *UserNotFound) Error() string { return "" }
