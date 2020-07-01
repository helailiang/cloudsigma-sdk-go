package cloudsigma

import (
	"context"
	"fmt"
	"net/http"
)

const subscriptionsBasePath = "subscriptions"

type SubscriptionsService service

type Subscription struct {
	Amount      string `json:"amount,omitempty"`
	AutoRenew   bool   `json:"auto_renew,omitempty"`
	FreeTier    bool   `json:"free_tier,omitempty"`
	ID          string `json:"id,omitempty"`
	Period      string `json:"period,omitempty"`
	Price       string `json:"price,omitempty"`
	Remaining   string `json:"remaining,omitempty"`
	Resource    string `json:"resource,omitempty"`
	ResourceURI string `json:"resource_uri,omitempty"`
	Status      string `json:"status,omitempty"`
	UUID        string `json:"uuid"`
}

type SubscriptionCreateRequest struct {
	Subscriptions []Subscription `json:"objects"`
}

type subscriptionsRoot struct {
	Meta          *Meta          `json:"meta,omitempty"`
	Subscriptions []Subscription `json:"objects,omitempty"`
}

func (s *SubscriptionsService) List(ctx context.Context) ([]Subscription, *Response, error) {
	path := fmt.Sprintf("%v/", subscriptionsBasePath)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(subscriptionsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.Subscriptions, resp, nil
}

func (s *SubscriptionsService) Create(ctx context.Context, subscriptionCreateRequest *SubscriptionCreateRequest) ([]Subscription, *Response, error) {
	if subscriptionCreateRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/", subscriptionsBasePath)

	req, err := s.client.NewRequest(http.MethodPost, path, subscriptionCreateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(SubscriptionCreateRequest)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Subscriptions, resp, nil
}
