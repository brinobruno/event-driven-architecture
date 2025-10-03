package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	client, _ := NewClient(mockName, mockEmail)
	account := NewAccount(client)
	assert.NotNil(t, account)
	assert.Equal(t, client.ID, account.Client.ID)
}

func TestCreateAccountWithNilClient(t *testing.T) {
	account := NewAccount(nil)
	assert.Nil(t, account)
}

func TestCreditAccount(t *testing.T) {
	client, _ := NewClient(mockName, mockEmail)
	account := NewAccount(client)
	account.Credit(10)
	account.Credit(15)
	assert.Equal(t, account.Balance, float64(25))
}

func TestDebitAccount(t *testing.T) {
	client, _ := NewClient(mockName, mockEmail)
	account := NewAccount(client)
	account.Credit(100)
	account.Debit(50)
	assert.Equal(t, account.Balance, float64(50))
}
