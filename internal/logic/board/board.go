package board

import (
	"errors"

	"github.com/good-threads/backend/internal/client/command"
	"github.com/good-threads/backend/internal/client/metric"
	"github.com/good-threads/backend/internal/client/mongo"
	"github.com/good-threads/backend/internal/client/thread"
	"github.com/good-threads/backend/internal/client/user"
	e "github.com/good-threads/backend/internal/errors"
	"github.com/mitchellh/mapstructure"
)

type Logic interface {
	Get(username string) ([]thread.Thread, []string, *string, error)
	Update(username string, lastProcessedCommandID *string, commands []Command) (*string, error)
}

type logic struct {
	userClient    user.Client
	commandClient command.Client
	threadClient  thread.Client
	metricClient  metric.Client
	mongoClient   mongo.Client
}

func Setup(userClient user.Client, commandClient command.Client, threadClient thread.Client, metricClient metric.Client, mongoClient mongo.Client) Logic {
	return &logic{userClient: userClient, commandClient: commandClient, threadClient: threadClient, metricClient: metricClient, mongoClient: mongoClient}
}

func (l *logic) Get(username string) ([]thread.Thread, []string, *string, error) {

	user, err := l.userClient.Fetch(username)
	if err != nil {
		return nil, nil, nil, err
	}

	lastProcessedCommandID, err := l.commandClient.FetchLastID(username)
	if err != nil {
		switch err.(type) {
		case *e.NoCommandFound:
			break
		default:
			return nil, nil, nil, err
		}
	}

	// TODO(thomasmarlow): refactor so that the aggregation is not executed
	//					   if there are no active threads
	activeThreads, hiddenThreads, err := l.threadClient.FetchAll(username, user.ActiveThreads)
	if err != nil {
		return nil, nil, nil, err
	}

	l.metricClient.RegisterBoardRead()

	return activeThreads, hiddenThreads, lastProcessedCommandID, nil
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
		decodeAndProcess, validCommandType := map[string]func(mongo.Transaction, *logic, string, any) error{
			"createThread":   getDecodeAndProcessFunction[PayloadCreateThread],
			"editThreadName": getDecodeAndProcessFunction[PayloadEditThreadName],
			"hideThread":     getDecodeAndProcessFunction[PayloadHideThread],
			"showThread":     getDecodeAndProcessFunction[PayloadShowThread],
			"relocateThread": getDecodeAndProcessFunction[PayloadRelocateThread],
			"createKnot":     getDecodeAndProcessFunction[PayloadCreateKnot],
			"editKnotBody":   getDecodeAndProcessFunction[PayloadEditKnotBody],
			"deleteKnot":     getDecodeAndProcessFunction[PayloadDeleteKnot],
		}[command.Type]
		if !validCommandType {
			return serversideLastProcessedCommandID, errors.New("shouldn't happen (TODO(thomasmarlow))")
		}
		if err := l.mongoClient.Transactionally(func(transaction mongo.Transaction) error {
			if err := decodeAndProcess(transaction, l, username, command.Payload); err != nil {
				return err
			}
			return l.commandClient.RegisterProcessed(transaction, username, command.ID)
		}); err != nil {
			return serversideLastProcessedCommandID, err
		}
		serversideLastProcessedCommandID = &command.ID
	}

	return serversideLastProcessedCommandID, nil
}

func getDecodeAndProcessFunction[Payload Processable](transaction mongo.Transaction, l *logic, username string, undecodedPayload any) error {
	var payload Payload
	if err := mapstructure.Decode(undecodedPayload, &payload); err != nil {
		return &e.BadPayload{}
	}
	return payload.Process(transaction, l, username)
}

type Processable interface {
	Process(transaction mongo.Transaction, l *logic, username string) error
}

func (p PayloadCreateThread) Process(transaction mongo.Transaction, l *logic, username string) error {
	if err := l.threadClient.Create(transaction, username, p.ID, p.Name); err != nil {
		return err
	}
	return l.userClient.AddThread(transaction, username, p.ID)
}

func (p PayloadEditThreadName) Process(transaction mongo.Transaction, l *logic, username string) error {
	return l.threadClient.EditName(transaction, username, p.ID, p.Name)
}

func (p PayloadHideThread) Process(transaction mongo.Transaction, l *logic, username string) error {
	return l.userClient.RemoveThread(transaction, username, p.ID)
}

func (p PayloadShowThread) Process(transaction mongo.Transaction, l *logic, username string) error {
	return l.userClient.AddThread(transaction, username, p.ID)
}

func (p PayloadRelocateThread) Process(transaction mongo.Transaction, l *logic, username string) error {
	return l.userClient.RelocateThread(transaction, username, p.ID, p.NewIndex)
}

func (p PayloadCreateKnot) Process(transaction mongo.Transaction, l *logic, username string) error {
	return l.threadClient.AddKnot(transaction, username, p.ThreadID, p.KnotID, p.KnotBody)
}

func (p PayloadEditKnotBody) Process(transaction mongo.Transaction, l *logic, username string) error {
	return l.threadClient.EditKnotBody(transaction, username, p.ThreadID, p.KnotID, p.KnotBody)
}

func (p PayloadDeleteKnot) Process(transaction mongo.Transaction, l *logic, username string) error {
	return l.threadClient.DeleteKnot(transaction, username, p.ThreadID, p.KnotID)
}
