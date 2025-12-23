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

func Create(project *ProjectCreateParams) (*Project, error) {
	if project == nil {
		return nil, fmt.Errorf("pushpad: project is required")
	}
	if project.SenderID == 0 {
		return nil, fmt.Errorf("pushpad: sender ID is required")
	}
	if project.Name == "" {
		return nil, fmt.Errorf("pushpad: project name is required")
	}
	if project.Website == "" {
		return nil, fmt.Errorf("pushpad: project website is required")
	}

	var created Project
	_, err := pushpad.DoRequest("POST", "/projects", nil, project, []int{201}, &created)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func Get(projectID int) (*Project, error) {
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

func Update(projectID int, update *ProjectUpdateParams) (*Project, error) {
	if update == nil {
		return nil, fmt.Errorf("pushpad: project update is required")
	}
	if projectID == 0 {
		return nil, fmt.Errorf("pushpad: project ID is required")
	}

	var project Project
	_, err := pushpad.DoRequest("PATCH", fmt.Sprintf("/projects/%d", projectID), nil, update, []int{200}, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func Delete(projectID int) error {
	if projectID == 0 {
		return fmt.Errorf("pushpad: project ID is required")
	}
	_, err := pushpad.DoRequest("DELETE", fmt.Sprintf("/projects/%d", projectID), nil, nil, []int{202}, nil)
	return err
}
