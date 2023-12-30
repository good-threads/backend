package board

import (
	"github.com/good-threads/backend/internal/client/changeset"
	"github.com/good-threads/backend/internal/client/thread"
	"github.com/good-threads/backend/internal/client/user"
	e "github.com/good-threads/backend/internal/errors"
)

type Logic interface {
	Get(username string) ([]thread.Thread, *string, error)
}

type logic struct {
	userClient      user.Client
	changesetClient changeset.Client
	threadClient    thread.Client
}

func Setup(userClient user.Client, changesetClient changeset.Client, threadClient thread.Client) Logic {
	return &logic{userClient: userClient, changesetClient: changesetClient, threadClient: threadClient}
}

func (l *logic) Get(username string) ([]thread.Thread, *string, error) {

	user, err := l.userClient.Fetch(username)
	if err != nil {
		return nil, nil, err
	}

	lastProcessedChangesetID, err := l.changesetClient.FetchLastID(username)
	if err != nil {
		switch err.(type) {
		case *e.NoChangesetFound:
			break
		default:
			return nil, nil, err
		}
	}

	if len(user.Threads) == 0 {
		return []thread.Thread{}, lastProcessedChangesetID, nil
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

	return threads, lastProcessedChangesetID, nil
}
