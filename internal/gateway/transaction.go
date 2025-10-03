package gateway

import "github.com/brinobruno/ms-wallet-core/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
