package user

import (
	mongoClient "github.com/good-threads/backend/internal/client/mongo"
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

func (m *MockClient) AddThread(transaction mongoClient.Transaction, username string, id string) error {
	args := m.Called(transaction, username, id)
	return args.Error(0)
}

func (m *MockClient) RemoveThread(transaction mongoClient.Transaction, username string, id string) error {
	args := m.Called(transaction, username, id)
	return args.Error(0)
}

func (m *MockClient) RelocateThread(transaction mongoClient.Transaction, username string, id string, newIndex uint) error {
	args := m.Called(transaction, username, id, newIndex)
	return args.Error(0)
}
