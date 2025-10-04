package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/brinobruno/ms-wallet-core/internal/database"
	"github.com/brinobruno/ms-wallet-core/internal/event"
	"github.com/brinobruno/ms-wallet-core/internal/event/handler"
	"github.com/brinobruno/ms-wallet-core/internal/seed"
	createaccount "github.com/brinobruno/ms-wallet-core/internal/usecase/create_account"
	createclient "github.com/brinobruno/ms-wallet-core/internal/usecase/create_client"
	createtransaction "github.com/brinobruno/ms-wallet-core/internal/usecase/create_transaction"
	"github.com/brinobruno/ms-wallet-core/internal/web"
	"github.com/brinobruno/ms-wallet-core/internal/web/webserver"
	"github.com/brinobruno/ms-wallet-core/pkg/events"
	"github.com/brinobruno/ms-wallet-core/pkg/kafka"
	"github.com/brinobruno/ms-wallet-core/pkg/uow"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root", "root", "mysql", "3306", "wallet",
	))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "walletcore",
	}
	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()

	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()
	eventDispatcher.Register("BalanceUpdated", handler.NewUpdateBalanceKafkaHandler(kafkaProducer))
	balanceUpdatedEvent := event.NewBalanceUpdated()

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})
	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createClientUseCase := createclient.NewCreateClientUseCase(clientDb)
	createAccountUseCase := createaccount.NewCreateAccountUseCase(accountDb, clientDb)
	createTransactionUseCase := createtransaction.NewCreateTransactionUseCase(
		uow,
		eventDispatcher,
		transactionCreatedEvent,
		balanceUpdatedEvent,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	webServer := webserver.NewWebServer(":" + port)
	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webServer.AddHandler("/clients", clientHandler.CreateClient)
	webServer.AddHandler("/accounts", accountHandler.CreateAccount)
	webServer.AddHandler("/transactions", transactionHandler.CreateTransaction)

	seed.Seed(ctx, *createClientUseCase, *createAccountUseCase, *createTransactionUseCase)

	fmt.Println("Server is running on port", port)
	webServer.Start()
}
