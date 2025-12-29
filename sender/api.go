package sender

import (
	"fmt"

	"github.com/pushpad/pushpad-go"
)

func List(params *SenderListParams) ([]Sender, error) {
	var senders []Sender
	_, err := pushpad.DoRequest("GET", "/senders", nil, nil, []int{200}, &senders)
	return senders, err
}

func Create(params *SenderCreateParams) (*Sender, error) {
	if params == nil {
		return nil, fmt.Errorf("pushpad: params are required")
	}

	var created Sender
	_, err := pushpad.DoRequest("POST", "/senders", nil, params, []int{201}, &created)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func Get(senderID int64, params *SenderGetParams) (*Sender, error) {
	if senderID == 0 {
		return nil, fmt.Errorf("pushpad: sender ID is required")
	}

	var sender Sender
	_, err := pushpad.DoRequest("GET", fmt.Sprintf("/senders/%d", senderID), nil, nil, []int{200}, &sender)
	if err != nil {
		return nil, err
	}
	return &sender, nil
}

func Update(senderID int64, params *SenderUpdateParams) (*Sender, error) {
	if params == nil {
		return nil, fmt.Errorf("pushpad: params are required")
	}
	if senderID == 0 {
		return nil, fmt.Errorf("pushpad: sender ID is required")
	}

	var sender Sender
	_, err := pushpad.DoRequest("PATCH", fmt.Sprintf("/senders/%d", senderID), nil, params, []int{200}, &sender)
	if err != nil {
		return nil, err
	}
	return &sender, nil
}

func Delete(senderID int64, params *SenderDeleteParams) error {
	if senderID == 0 {
		return fmt.Errorf("pushpad: sender ID is required")
	}
	_, err := pushpad.DoRequest("DELETE", fmt.Sprintf("/senders/%d", senderID), nil, nil, []int{204}, nil)
	return err
}
