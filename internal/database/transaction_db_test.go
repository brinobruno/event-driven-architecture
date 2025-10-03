package database

import (
	"database/sql"
	"testing"

	"github.com/brinobruno/ms-wallet-core/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type TransactionDBTestSuite struct {
	suite.Suite
	db            *sql.DB
	transactionDB *TransactionDB
	client1       *entity.Client
	client2       *entity.Client
	accountFrom   *entity.Account
	accountTo     *entity.Account
}

func (s *TransactionDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec(
		"CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at date)",
	)
	db.Exec(
		"CREATE TABLE accounts (id varchar(255), client_id varchar(255), balance float, created_at date)",
	)
	db.Exec(
		"CREATE TABLE transactions (id varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount float, created_at date)",
	)
	s.transactionDB = NewTransactionDB(db)

	client1, _ := entity.NewClient("Client 1", "j@1.com")
	s.client1 = client1
	client2, _ := entity.NewClient("Client 2", "j@2.com")
	s.client2 = client2

	accountFrom := entity.NewAccount(client1)
	accountFrom.Balance = 1000
	s.accountFrom = accountFrom
	accountTo := entity.NewAccount(client2)
	accountTo.Balance = 1000
	s.accountTo = accountTo
}

func (s *TransactionDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE clients")
	s.db.Exec("DROP TABLE accounts")
	s.db.Exec("DROP TABLE transactions")
}

func TestTransactionDBTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDBTestSuite))
}

func (s *TransactionDBTestSuite) TestCreate() {
	transaction, err := entity.NewTransaction(s.accountFrom, s.accountTo, 100)
	s.Nil(err)
	err = s.transactionDB.Create(transaction)
	s.Nil(err)
}
