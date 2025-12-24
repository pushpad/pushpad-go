package notification

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/pushpad/pushpad-go"
)

func List(params *NotificationListParams) ([]Notification, error) {
	projectID := 0
	if params != nil {
		projectID = params.ProjectID
	}
	projectID, err := pushpad.ResolveProjectID(projectID)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	if params != nil && params.Page > 0 {
		query.Set("page", strconv.Itoa(params.Page))
	}

	var notifications []Notification
	_, err = pushpad.DoRequest("GET", fmt.Sprintf("/projects/%d/notifications", projectID), query, nil, []int{200}, &notifications)
	return notifications, err
}

func Create(params *NotificationCreateParams) (*NotificationCreateResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("pushpad: params are required")
	}
	if params.Body == "" {
		return nil, fmt.Errorf("pushpad: params.Body is required")
	}
	projectID, err := pushpad.ResolveProjectID(params.ProjectID)
	if err != nil {
		return nil, err
	}

	var response NotificationCreateResponse
	_, err = pushpad.DoRequest("POST", fmt.Sprintf("/projects/%d/notifications", projectID), nil, params, []int{201}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func Get(notificationID int, params *NotificationGetParams) (*Notification, error) {
	if notificationID == 0 {
		return nil, fmt.Errorf("pushpad: notification ID is required")
	}
	var notification Notification
	_, err := pushpad.DoRequest("GET", fmt.Sprintf("/notifications/%d", notificationID), nil, nil, []int{200}, &notification)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func Cancel(notificationID int, params *NotificationCancelParams) error {
	if notificationID == 0 {
		return fmt.Errorf("pushpad: notification ID is required")
	}
	_, err := pushpad.DoRequest("DELETE", fmt.Sprintf("/notifications/%d/cancel", notificationID), nil, nil, []int{204}, nil)
	return err
}
