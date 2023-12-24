package user

import (
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Persist(username string, passwordHash []byte) error {
	args := m.Called(username, passwordHash)
	return args.Error(0)
}
