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

type NoCommandFound struct{}

func (e *NoCommandFound) Error() string { return "" }

type NoThreadsFound struct{}

func (e *NoThreadsFound) Error() string { return "" }

type ReceivedCommandsWouldRewriteHistory struct{}

func (e *ReceivedCommandsWouldRewriteHistory) Error() string { return "" }

type GeneratedIDClashed struct{}

func (e *GeneratedIDClashed) Error() string { return "" }

type ThreadNotFound struct{}

func (e *ThreadNotFound) Error() string { return "" }

type KnotNotFound struct{}

func (e *KnotNotFound) Error() string { return "" }

type BadPayload struct{}

func (e *BadPayload) Error() string { return "" }
