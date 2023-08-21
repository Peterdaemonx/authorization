package mock

import (
	"context"

	gcloudPubsub "cloud.google.com/go/pubsub"
	"gitlab.cmpayments.local/creditcard/platform/events/pubsub"
)

type Publisher struct{}

func (p Publisher) CreateTopic(ctx context.Context, topicID string) *gcloudPubsub.Topic {
	return nil
}

func (p Publisher) CreateSubscription(ctx context.Context, subscriptionID string, topic *gcloudPubsub.Topic) *gcloudPubsub.Subscription {
	return nil
}

func NewMockPublisher() Publisher {
	return Publisher{}
}

func (p Publisher) Publish(ctx context.Context, topic string, _ pubsub.Publishable) error {
	return nil
}
