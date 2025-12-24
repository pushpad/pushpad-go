package project

import (
	"fmt"

	"github.com/pushpad/pushpad-go"
)

func List(params *ProjectListParams) ([]Project, error) {
	var projects []Project
	_, err := pushpad.DoRequest("GET", "/projects", nil, nil, []int{200}, &projects)
	return projects, err
}

func Create(params *ProjectCreateParams) (*Project, error) {
	if params == nil {
		return nil, fmt.Errorf("pushpad: params are required")
	}
	if params.SenderID == 0 {
		return nil, fmt.Errorf("pushpad: sender ID is required")
	}
	if params.Name == "" {
		return nil, fmt.Errorf("pushpad: params.Name is required")
	}
	if params.Website == "" {
		return nil, fmt.Errorf("pushpad: params.Website is required")
	}

	var created Project
	_, err := pushpad.DoRequest("POST", "/projects", nil, params, []int{201}, &created)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func Get(projectID int, params *ProjectGetParams) (*Project, error) {
	if projectID == 0 {
		return nil, fmt.Errorf("pushpad: project ID is required")
	}

	var project Project
	_, err := pushpad.DoRequest("GET", fmt.Sprintf("/projects/%d", projectID), nil, nil, []int{200}, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func Update(projectID int, params *ProjectUpdateParams) (*Project, error) {
	if params == nil {
		return nil, fmt.Errorf("pushpad: params are required")
	}
	if projectID == 0 {
		return nil, fmt.Errorf("pushpad: project ID is required")
	}

	var project Project
	_, err := pushpad.DoRequest("PATCH", fmt.Sprintf("/projects/%d", projectID), nil, params, []int{200}, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func Delete(projectID int, params *ProjectDeleteParams) error {
	if projectID == 0 {
		return fmt.Errorf("pushpad: project ID is required")
	}
	_, err := pushpad.DoRequest("DELETE", fmt.Sprintf("/projects/%d", projectID), nil, nil, []int{202}, nil)
	return err
}
