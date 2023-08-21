package cardInfo

import (
	"context"
	"fmt"

	"gitlab.cmpayments.local/creditcard/authorization/internal/processing"
	"gitlab.cmpayments.local/creditcard/card-info-api/pkg/cardinfoapi"
)

func NewService(client cardinfoapi.Client) service {
	return service{client: client}
}

type service struct {
	client cardinfoapi.Client
}

func (s service) FetchCardInfo(ctx context.Context, pan string) (processing.CardInfo, error) {
	infoResponse, err := s.client.Info(ctx, pan)
	if err != nil {
		return processing.CardInfo{}, fmt.Errorf("cannot fetch card info: %w", err)
	}
	return mapCardInfo(infoResponse), nil
}

func mapCardInfo(response cardinfoapi.InfoResponse) processing.CardInfo {
	return processing.CardInfo{Info: processing.Info{
		Low:      response.Info.Low,
		High:     response.Info.High,
		Scheme:   response.Info.Scheme,
		Programs: mapProgram(response.Info.Programs),
	},
	}
}

func mapProgram(programs []cardinfoapi.Program) []processing.Program {
	var mappedProgram processing.Program
	var mappedPrograms []processing.Program
	for _, program := range programs {
		mappedProgram.ID = program.ID
		mappedProgram.Priority = program.Priority
		mappedProgram.Issuer = mapIssuer(program.Issuer)
		mappedPrograms = append(mappedPrograms, mappedProgram)
	}
	return mappedPrograms
}

func mapIssuer(issuer cardinfoapi.Issuer) processing.Issuer {
	return processing.Issuer{
		Name:        issuer.Name,
		Country:     issuer.Country,
		Website:     issuer.Website,
		Phonenumber: issuer.Phonenumber,
	}
}
