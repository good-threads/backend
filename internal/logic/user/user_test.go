package user

import (
	"testing"

	"github.com/good-threads/backend/internal/client/user"
	e "github.com/good-threads/backend/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateFailsIfEmptyUsername(t *testing.T) {
	_, errorIsBadUsername := Setup(nil).Create("", "asd").(*e.BadUsername)
	assert.True(t, errorIsBadUsername)
}

func TestCreateFailsIfEmptyPassword(t *testing.T) {
	_, errorIsBadPassword := Setup(nil).Create("asd", "").(*e.BadPassword)
	assert.True(t, errorIsBadPassword)
}

func TestCreateFailsIfLongPassword(t *testing.T) {
	err := Setup(nil).
		Create("asd", "12345678901234567890123456789012345678901234567890123456789012345678901234567890")
	assert.NotNil(t, err)
}

func TestCreateOk(t *testing.T) {
	client := &user.MockClient{}
	client.On("Persist", "asd", mock.Anything).Return(nil)
	assert.Nil(t, Setup(client).Create("asd", "asd"))
}
