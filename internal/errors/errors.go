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

type InvalidSession struct{}

func (e *InvalidSession) Error() string { return "" }

type SessionNotFound struct{}

func (e *SessionNotFound) Error() string { return "" }

type NoChangesetFound struct{}

func (e *NoChangesetFound) Error() string { return "" }

type NoThreadsFound struct{}

func (e *NoThreadsFound) Error() string { return "" }
