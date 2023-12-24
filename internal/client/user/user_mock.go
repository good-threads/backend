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

func (m *MockClient) Fetch(username string) (*User, error) {
	args := m.Called(username)
	return args.Get(0).(*User), args.Error(1)
}
