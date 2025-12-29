package pushpad

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const DefaultBaseURL = "https://pushpad.xyz/api/v1"

// APIError represents a non-2xx API response.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("pushpad: unexpected status code %d", e.StatusCode)
	}
	return fmt.Sprintf("pushpad: status %d: %s", e.StatusCode, e.Body)
}

// ResolveProjectID returns the provided project ID or the configured default project ID.
func ResolveProjectID(projectID *int64) (int64, error) {
	if projectID != nil && *projectID != 0 {
		return *projectID, nil
	}
	if pushpadProjectID == 0 {
		return 0, fmt.Errorf("pushpad: project ID is required")
	}
	return pushpadProjectID, nil
}

// DoRequest performs an HTTP request against the Pushpad API.
func DoRequest(method, path string, query url.Values, body any, okStatuses []int, out any) (*http.Response, error) {
	ctx := context.Background()
	baseURL := strings.TrimRight(DefaultBaseURL, "/")
	endpoint := baseURL + path
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if pushpadAuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+pushpadAuthToken)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	ok := false
	for _, code := range okStatuses {
		if res.StatusCode == code {
			ok = true
			break
		}
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if !ok {
		return res, &APIError{StatusCode: res.StatusCode, Body: string(bodyBytes)}
	}

	if out != nil && len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, out); err != nil {
			return res, err
		}
	}

	return res, nil
}
