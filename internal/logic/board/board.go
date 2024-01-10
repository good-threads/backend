package board

import (
	"errors"

	"github.com/good-threads/backend/internal/client/command"
	"github.com/good-threads/backend/internal/client/thread"
	"github.com/good-threads/backend/internal/client/user"
	e "github.com/good-threads/backend/internal/errors"
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
		switch command.Type {
		case "createThread":
			payload := command.Payload.(PayloadCreateThread)
			if err := l.threadClient.Create(username, payload.ID, payload.Name); err != nil {
				return serversideLastProcessedCommandID, err
			}
			if err := l.userClient.AddThread(username, payload.ID); err != nil {
				return serversideLastProcessedCommandID, err
			}
		case "editThreadName":
			payload := command.Payload.(PayloadEditThreadName)
			if err := l.threadClient.EditName(username, payload.ID, payload.Name); err != nil {
				return serversideLastProcessedCommandID, err
			}
		case "hideThread":
			payload := command.Payload.(PayloadHideThread)
			if err := l.userClient.RemoveThread(username, payload.ID); err != nil {
				return serversideLastProcessedCommandID, err
			}
		case "repositionThread":
			payload := command.Payload.(PayloadRelocateThread)
			if err := l.userClient.RelocateThread(username, payload.ID, payload.NewIndex); err != nil {
				return serversideLastProcessedCommandID, err
			}
		case "createKnot":
			payload := command.Payload.(PayloadCreateKnot)
			if err := l.threadClient.AddKnot(username, payload.ThreadID, payload.KnotID, payload.KnotBody); err != nil {
				return serversideLastProcessedCommandID, err
			}
		case "editKnot":
			payload := command.Payload.(PayloadEditKnot)
			if err := l.threadClient.EditKnot(username, payload.ThreadID, payload.KnotID, payload.KnotBody); err != nil {
				return serversideLastProcessedCommandID, err
			}
		case "deleteKnot":
			payload := command.Payload.(PayloadDeleteKnot)
			if err := l.threadClient.DeleteKnot(username, payload.ThreadID, payload.KnotID); err != nil {
				return serversideLastProcessedCommandID, err
			}
		default:
			return serversideLastProcessedCommandID, errors.New("shouldn't happen (TODO(thomasmarlow))")
		}
		if err := l.commandClient.RegisterProcessed(username, command.ID); err != nil {
			return serversideLastProcessedCommandID, err
		}
		serversideLastProcessedCommandID = &command.ID
	}

	return serversideLastProcessedCommandID, nil
}
