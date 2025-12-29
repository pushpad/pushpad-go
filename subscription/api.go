package subscription

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/pushpad/pushpad-go/v1"
)

func List(params *SubscriptionListParams) ([]Subscription, error) {
	if params == nil {
		params = &SubscriptionListParams{}
	}
	projectID, err := pushpad.ResolveProjectID(params.ProjectID)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	if params.Page != nil && *params.Page > 0 {
		query.Set("page", strconv.FormatInt(*params.Page, 10))
	}
	if params.PerPage != nil && *params.PerPage > 0 {
		query.Set("per_page", strconv.FormatInt(*params.PerPage, 10))
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
	_, err = pushpad.DoRequest("GET", fmt.Sprintf("/projects/%d/subscriptions", projectID), query, nil, []int{200}, &subscriptions)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func Count(params *SubscriptionCountParams) (int64, error) {
	if params == nil {
		params = &SubscriptionCountParams{}
	}
	projectID, err := pushpad.ResolveProjectID(params.ProjectID)
	if err != nil {
		return 0, err
	}

	query := url.Values{}
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
	res, err := pushpad.DoRequest("HEAD", fmt.Sprintf("/projects/%d/subscriptions", projectID), query, nil, []int{200}, nil)
	if err != nil {
		return 0, err
	}

	var totalCount int64
	if header := res.Header.Get("X-Total-Count"); header != "" {
		if parsed, parseErr := strconv.ParseInt(header, 10, 64); parseErr == nil {
			totalCount = parsed
		}
	}

	return totalCount, nil
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

func Get(subscriptionID int64, params *SubscriptionGetParams) (*Subscription, error) {
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

func Update(subscriptionID int64, params *SubscriptionUpdateParams) (*Subscription, error) {
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

func Delete(subscriptionID int64, params *SubscriptionDeleteParams) error {
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
