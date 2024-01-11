package board

import (
	"errors"
	"log"

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
		switch command.Type {
		case "createThread":
			payload, err := decode[PayloadCreateThread](command.Payload)
			if err != nil {
				return serversideLastProcessedCommandID, &e.BadPayload{}
			}
			if err := l.threadClient.Create(username, payload.ID, payload.Name); err != nil {
				return serversideLastProcessedCommandID, err
			}
			if err := l.userClient.AddThread(username, payload.ID); err != nil {
				return serversideLastProcessedCommandID, err
			}
		case "editThreadName":
			payload, err := decode[PayloadEditThreadName](command.Payload)
			if err != nil {
				return serversideLastProcessedCommandID, &e.BadPayload{}
			}
			if err := l.threadClient.EditName(username, payload.ID, payload.Name); err != nil {
				return serversideLastProcessedCommandID, err
			}
		case "hideThread":
			payload, err := decode[PayloadHideThread](command.Payload)
			if err != nil {
				return serversideLastProcessedCommandID, &e.BadPayload{}
			}
			if err := l.userClient.RemoveThread(username, payload.ID); err != nil {
				return serversideLastProcessedCommandID, err
			}
		case "relocateThread":
			payload, err := decode[PayloadRelocateThread](command.Payload)
			if err != nil {
				return serversideLastProcessedCommandID, &e.BadPayload{}
			}
			if err := l.userClient.RelocateThread(username, payload.ID, payload.NewIndex); err != nil {
				return serversideLastProcessedCommandID, err
			}
		case "createKnot":
			payload, err := decode[PayloadCreateKnot](command.Payload)
			if err != nil {
				return serversideLastProcessedCommandID, &e.BadPayload{}
			}
			log.Println(payload)
			if err := l.threadClient.AddKnot(username, payload.ThreadID, payload.KnotID, payload.KnotBody); err != nil {
				return serversideLastProcessedCommandID, err
			}
		case "editKnot":
			payload, err := decode[PayloadEditKnot](command.Payload)
			if err != nil {
				return serversideLastProcessedCommandID, &e.BadPayload{}
			}
			if err := l.threadClient.EditKnot(username, payload.ThreadID, payload.KnotID, payload.KnotBody); err != nil {
				return serversideLastProcessedCommandID, err
			}
		case "deleteKnot":
			payload, err := decode[PayloadDeleteKnot](command.Payload)
			if err != nil {
				return serversideLastProcessedCommandID, &e.BadPayload{}
			}
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

func decode[Payload any](undecodedPayload any) (Payload, error) {
	var payload Payload
	err := mapstructure.Decode(undecodedPayload, &payload)
	return payload, err
}
