package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const mockName = "John Doe"
const mockEmail = "j@j.com"

func TestCreateNewClient(t *testing.T) {
	client, err := NewClient(mockName, mockEmail)

	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, client.Name, mockName)
	assert.Equal(t, client.Email, mockEmail)
}

func TestCreateNewClientWhenArgsAreInvalid(t *testing.T) {
	client, err := NewClient("", "")
	assert.NotNil(t, err)
	assert.Nil(t, client)
}

func TestCreateUpdateClient(t *testing.T) {
	client, _ := NewClient(mockName, mockEmail)
	err := client.Update("John Doe Jr", "jd@jr.com")

	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, client.Name, "John Doe Jr")
	assert.Equal(t, client.Email, "jd@jr.com")
}

func TestCreateUpdateClientWhenArgsAreInvalid(t *testing.T) {
	client, _ := NewClient(mockName, mockEmail)
	err := client.Update("", "j@r.com")
	assert.Error(t, err, "name is required")
}

func TestAddAccountToClient(t *testing.T) {
	client, _ := NewClient(mockName, mockEmail)
	account := NewAccount(client)
	err := client.AddAccount(account)
	assert.Nil(t, err)
	assert.Equal(t, len(client.Accounts), 1)
}
