package gateway

import "github.com/brinobruno/ms-wallet-core/internal/entity"

type AccountGateway interface {
	Save(*entity.Account) error
	FindByID(id string) (*entity.Account, error)
	UpdateBalance(account *entity.Account) error
}
