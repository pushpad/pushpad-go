package subscription

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/pushpad/pushpad-go"
)

func List(params *SubscriptionListParams) ([]Subscription, int, error) {
	if params == nil {
		params = &SubscriptionListParams{}
	}
	projectID, err := pushpad.ResolveProjectID(params.ProjectID)
	if err != nil {
		return nil, 0, err
	}

	query := url.Values{}
	if params.Page != nil && *params.Page > 0 {
		query.Set("page", strconv.Itoa(*params.Page))
	}
	if params.PerPage != nil && *params.PerPage > 0 {
		query.Set("per_page", strconv.Itoa(*params.PerPage))
	}
	if params.UIDs != nil {
		for _, uid := range *params.UIDs {
			query.Add("uids[]", uid)
		}
	}
	if params.Tags != nil {
		for _, tag := range *params.Tags {
			query.Add("tags[]", tag)
		}
	}

	var subscriptions []Subscription
	res, err := pushpad.DoRequest("GET", fmt.Sprintf("/projects/%d/subscriptions", projectID), query, nil, []int{200}, &subscriptions)
	if err != nil {
		return nil, 0, err
	}

	totalCount := 0
	if header := res.Header.Get("X-Total-Count"); header != "" {
		if parsed, parseErr := strconv.Atoi(header); parseErr == nil {
			totalCount = parsed
		}
	}

	return subscriptions, totalCount, nil
}

func Create(params *SubscriptionCreateParams) (*Subscription, error) {
	if params == nil {
		return nil, fmt.Errorf("pushpad: params are required")
	}
	projectID, err := pushpad.ResolveProjectID(params.ProjectID)
	if err != nil {
		return nil, err
	}

	var created Subscription
	_, err = pushpad.DoRequest("POST", fmt.Sprintf("/projects/%d/subscriptions", projectID), nil, params, []int{201}, &created)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func Get(subscriptionID int, params *SubscriptionGetParams) (*Subscription, error) {
	if params == nil {
		params = &SubscriptionGetParams{}
	}
	if subscriptionID == 0 {
		return nil, fmt.Errorf("pushpad: subscription ID is required")
	}
	projectID, err := pushpad.ResolveProjectID(params.ProjectID)
	if err != nil {
		return nil, err
	}

	var subscription Subscription
	_, err = pushpad.DoRequest("GET", fmt.Sprintf("/projects/%d/subscriptions/%d", projectID, subscriptionID), nil, nil, []int{200}, &subscription)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func Update(subscriptionID int, params *SubscriptionUpdateParams) (*Subscription, error) {
	if params == nil {
		return nil, fmt.Errorf("pushpad: params are required")
	}
	if subscriptionID == 0 {
		return nil, fmt.Errorf("pushpad: subscription ID is required")
	}
	projectID, err := pushpad.ResolveProjectID(params.ProjectID)
	if err != nil {
		return nil, err
	}

	var subscription Subscription
	_, err = pushpad.DoRequest("PATCH", fmt.Sprintf("/projects/%d/subscriptions/%d", projectID, subscriptionID), nil, params, []int{200}, &subscription)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func Delete(subscriptionID int, params *SubscriptionDeleteParams) error {
	if params == nil {
		params = &SubscriptionDeleteParams{}
	}
	if subscriptionID == 0 {
		return fmt.Errorf("pushpad: subscription ID is required")
	}
	projectID, err := pushpad.ResolveProjectID(params.ProjectID)
	if err != nil {
		return err
	}
	_, err = pushpad.DoRequest("DELETE", fmt.Sprintf("/projects/%d/subscriptions/%d", projectID, subscriptionID), nil, nil, []int{204}, nil)
	return err
}
