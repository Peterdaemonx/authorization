package main

import (
	"fmt"
	timingwrappers "gitlab.cmpayments.local/creditcard/authorization/internal/timing/wrappers"

	cardinfoclient "gitlab.cmpayments.local/creditcard/card-info-api/pkg/cardinfoapi"
	platformclient "gitlab.cmpayments.local/creditcard/platform/http/client"
	tokenizationclient "gitlab.cmpayments.local/creditcard/tokenization-api/pkg/tokenization"

	"gitlab.cmpayments.local/creditcard/authorization/internal/authorization/app"
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/cardInfo"
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/mock"
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/tokenization"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing"
)

func (app *application) TokenizationService() (app.Tokenizer, error) {
	if app.conf.Development.MockTokenization {
		return mock.Tokenization{}, nil
	}

	tokenizationClient, err := platformclient.New(app.conf.Tokenization.BaseURL, app.conf.JWT, app.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create tokenization client: %w", err)
	}

	detokenizationClient, err := platformclient.New(app.conf.Detokenization.BaseURL, app.conf.JWT, app.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create tokenization client: %w", err)
	}

	tokenizer := tokenization.NewService(
		tokenizationclient.NewClient(tokenizationClient),
		tokenizationclient.NewClient(detokenizationClient),
		app.logger,
	)

	return timingwrappers.AuthorizationTokenizer{Base: tokenizer}, nil
}

func (app *application) CardInfoService() (processing.CardInfoService, error) {
	if app.conf.Development.MockCardInfo {
		return mock.CardInfoApi{}, nil
	}

	// uses fake jwt token. The cardinfoapi does not use a token.
	pc, err := platformclient.New(app.conf.CardInfoApi.BaseURL, app.conf.JWT, app.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create card info client: %w", err)
	}

	return cardInfo.NewService(
		cardinfoclient.NewClient(pc),
	), nil
}
