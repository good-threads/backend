package board

import (
	"errors"

	"github.com/good-threads/backend/internal/client/command"
	"github.com/good-threads/backend/internal/client/thread"
	"github.com/good-threads/backend/internal/client/user"
	e "github.com/good-threads/backend/internal/errors"
	"github.com/mitchellh/mapstructure"
)

type Logic interface {
	Get(username string) ([]thread.Thread, *string, error)
	Update(username string, lastProcessedCommandID *string, commands []Command) (*string, error)
}

type logic struct {
	userClient    user.Client
	commandClient command.Client
	threadClient  thread.Client
}

func Setup(userClient user.Client, commandClient command.Client, threadClient thread.Client) Logic {
	return &logic{userClient: userClient, commandClient: commandClient, threadClient: threadClient}
}

func (l *logic) Get(username string) ([]thread.Thread, *string, error) {

	user, err := l.userClient.Fetch(username)
	if err != nil {
		return nil, nil, err
	}

	lastProcessedCommandID, err := l.commandClient.FetchLastID(username)
	if err != nil {
		switch err.(type) {
		case *e.NoCommandFound:
			break
		default:
			return nil, nil, err
		}
	}

	if len(user.Threads) == 0 {
		return []thread.Thread{}, lastProcessedCommandID, nil
	}

	threads, err := l.threadClient.Fetch(user.Threads)
	if err != nil {
		switch err.(type) {
		case *e.NoThreadsFound:
			break
		default:
			return nil, nil, err
		}
	}

	return threads, lastProcessedCommandID, nil
}

func (l *logic) Update(username string, clientsideLastProcessedCommandID *string, commands []Command) (*string, error) {

	serversideLastProcessedCommandID, err := l.commandClient.FetchLastID(username)
	if err != nil {
		switch err.(type) {
		case *e.NoCommandFound:
			break
		default:
			return serversideLastProcessedCommandID, err
		}
	}

	if clientsideLastProcessedCommandID == nil && serversideLastProcessedCommandID != nil {
		return serversideLastProcessedCommandID, &e.ReceivedCommandsWouldRewriteHistory{}
	}

	if clientsideLastProcessedCommandID != nil && serversideLastProcessedCommandID == nil {
		return nil, &e.ReceivedCommandsWouldRewriteHistory{}
	}

	if clientsideLastProcessedCommandID != nil && serversideLastProcessedCommandID != nil && *clientsideLastProcessedCommandID != *serversideLastProcessedCommandID {
		return serversideLastProcessedCommandID, &e.ReceivedCommandsWouldRewriteHistory{}
	}

	for _, command := range commands {
		decodeAndProcess, validCommandType := map[string]func(*logic, string, any) error{
			"createThread":   getDecodeAndProcessFunction[PayloadCreateThread],
			"editThreadName": getDecodeAndProcessFunction[PayloadEditThreadName],
			"hideThread":     getDecodeAndProcessFunction[PayloadHideThread],
			"relocateThread": getDecodeAndProcessFunction[PayloadRelocateThread],
			"createKnot":     getDecodeAndProcessFunction[PayloadCreateKnot],
			"editKnot":       getDecodeAndProcessFunction[PayloadEditKnot],
			"deleteKnot":     getDecodeAndProcessFunction[PayloadDeleteKnot],
		}[command.Type]
		if !validCommandType {
			return serversideLastProcessedCommandID, errors.New("shouldn't happen (TODO(thomasmarlow))")
		}
		if err := decodeAndProcess(l, username, command.Payload); err != nil {
			return serversideLastProcessedCommandID, err
		}
		if err := l.commandClient.RegisterProcessed(username, command.ID); err != nil {
			return serversideLastProcessedCommandID, err
		}
		serversideLastProcessedCommandID = &command.ID
	}

	return serversideLastProcessedCommandID, nil
}

func getDecodeAndProcessFunction[Payload Processable](l *logic, username string, undecodedPayload any) error {
	var payload Payload
	if err := mapstructure.Decode(undecodedPayload, &payload); err != nil {
		return &e.BadPayload{}
	}
	return payload.Process(l, username)
}

type Processable interface {
	Process(l *logic, username string) error
}

func (p PayloadCreateThread) Process(l *logic, username string) error {
	if err := l.threadClient.Create(username, p.ID, p.Name); err != nil {
		return err
	}
	return l.userClient.AddThread(username, p.ID)
}

func (p PayloadEditThreadName) Process(l *logic, username string) error {
	return l.threadClient.EditName(username, p.ID, p.Name)
}

func (p PayloadHideThread) Process(l *logic, username string) error {
	return l.userClient.RemoveThread(username, p.ID)
}

func (p PayloadRelocateThread) Process(l *logic, username string) error {
	return l.userClient.RelocateThread(username, p.ID, p.NewIndex)
}

func (p PayloadCreateKnot) Process(l *logic, username string) error {
	return l.threadClient.AddKnot(username, p.ThreadID, p.KnotID, p.KnotBody)
}

func (p PayloadEditKnot) Process(l *logic, username string) error {
	return l.threadClient.EditKnot(username, p.ThreadID, p.KnotID, p.KnotBody)
}

func (p PayloadDeleteKnot) Process(l *logic, username string) error {
	return l.threadClient.DeleteKnot(username, p.ThreadID, p.KnotID)
}
