package main

import (
	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/mock"
	timingwrappers "gitlab.cmpayments.local/creditcard/authorization/internal/timing/wrappers"
	"gitlab.cmpayments.local/creditcard/platform/events/pubsub"
)

func (app *application) MessagePublisher() (pubsub.Publisher, error) {
	if app.conf.Development.MockPublisher {
		//os.Setenv("PUBSUB_EMULATOR_HOST", fmt.Sprintf("0.0.0.0:%v", 8085))
		return mock.NewMockPublisher(), nil
	}

	connection, err := pubsub.NewConnection(
		app.ctx, app.conf.GCP.ProjectID, app.logger, app.conf.GCP.PubSub.Timeout)
	if err != nil {
		return nil, err
	}

	return timingwrappers.Publisher{Base: connection}, nil
}
