package seed

import (
	"context"
	"database/sql"
	"fmt"

	createaccount "github.com/brinobruno/ms-wallet-core/internal/usecase/create_account"
	createclient "github.com/brinobruno/ms-wallet-core/internal/usecase/create_client"
	createtransaction "github.com/brinobruno/ms-wallet-core/internal/usecase/create_transaction"
)

func Seed(ctx context.Context,
	createClientUseCase createclient.CreateClientUseCase,
	createAccountUseCase createaccount.CreateAccountUseCase,
	createTransactionUseCase createtransaction.CreateTransactionUseCase,
) {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root", "root", "mysql", "3306", "wallet",
	))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	client1, err := createClientUseCase.Execute(createclient.CreateClientInputDTO{
		Name:  "Seed Client 1",
		Email: "seed@client1.com",
	})
	if err != nil {
		panic(err)
	}
	account1, err := createAccountUseCase.Execute(createaccount.CreateAccountInputDTO{
		ClientID: client1.ID,
	})
	if err != nil {
		panic(err)
	}
	account2, err := createAccountUseCase.Execute(createaccount.CreateAccountInputDTO{
		ClientID: client1.ID,
	})
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare(
		"UPDATE accounts SET balance = ? WHERE id = ?",
	)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(1000, account1.ID)
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(1000, account2.ID)
	if err != nil {
		panic(err)
	}

	output, err := createTransactionUseCase.Execute(ctx, createtransaction.CreateTransactionInputDTO{
		AccountIDFrom: account1.ID,
		AccountIDTo:   account2.ID,
		Amount:        1,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("db seeded", output)
}
