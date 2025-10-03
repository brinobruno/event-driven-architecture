package gateway

import "github.com/brinobruno/ms-wallet-core/internal/entity"

type ClientGateway interface {
	Get(id string) (*entity.Client, error)
	Save(*entity.Client) error
}
