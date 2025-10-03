package createaccount

import (
	"github.com/brinobruno/ms-wallet-core/internal/entity"
	"github.com/brinobruno/ms-wallet-core/internal/gateway"
)

type CreateAccountInputDTO struct {
	ClientID string `json:"client_id"`
}

type CreateAccountOutputDTO struct {
	ID string
}

type CreateAccountUseCase struct {
	AccountGateway gateway.AccountGateway
	ClientGateway  gateway.ClientGateway
}

func NewCreateAccountUseCase(
	accountGateway gateway.AccountGateway,
	clientGateway gateway.ClientGateway,
) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		AccountGateway: accountGateway,
		ClientGateway:  clientGateway,
	}
}

func (usecase *CreateAccountUseCase) Execute(
	input CreateAccountInputDTO,
) (*CreateAccountOutputDTO, error) {
	client, err := usecase.ClientGateway.Get(input.ClientID)
	if err != nil {
		return nil, err
	}
	account := entity.NewAccount(client)
	err = usecase.AccountGateway.Save(account)
	if err != nil {
		return nil, err
	}
	output := &CreateAccountOutputDTO{
		ID: account.ID,
	}
	return output, nil
}
