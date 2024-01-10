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

func (m *MockClient) AddThread(username string, id string) error {
	args := m.Called(username, id)
	return args.Error(0)
}

func (m *MockClient) RemoveThread(username string, id string) error {
	args := m.Called(username, id)
	return args.Error(0)
}

func (m *MockClient) RelocateThread(username string, id string, newIndex uint) error {
	args := m.Called(username, id, newIndex)
	return args.Error(0)
}
