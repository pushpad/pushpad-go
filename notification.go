package pushpad

import (
  "fmt"
  "bytes"
  "io"
  "encoding/json"
  "net/http"
)

type Notification struct {
  ProjectID string `json:"-"`
  Body string `json:"body"`
  Title string `json:"title,omitempty"`
  TargetUrl string `json:"target_url,omitempty"`
  IconUrl string `json:"icon_url,omitempty"`
  BadgeUrl string `json:"badge_url,omitempty"`
  ImageUrl string `json:"image_url,omitempty"`
  Ttl int `json:"ttl,omitempty"`
  RequireInteraction bool `json:"require_interaction,omitempty"`
  Silent bool `json:"silent,omitempty"`
  Urgent bool `json:"urgent,omitempty"`
  CustomData string `json:"custom_data,omitempty"`
  CustomMetrics []string `json:"custom_metrics,omitempty"`
  Starred bool `json:"starred,omitempty"`
}

func (n Notification) Broadcast() (string, error) {
  if n.ProjectID == "" {
    n.ProjectID = pushpadProjectID
  }
  
  notificationJSON, err := json.Marshal(n)
  
  if err != nil {
    return "", err
  }
  
  req, err := http.NewRequest("POST", "https://pushpad.xyz/api/v1/projects/" + n.ProjectID + "/notifications", bytes.NewBuffer(notificationJSON))
  
  if err != nil {
    return "", err
  }
  
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Authorization", "Token token=\"" + pushpadAuthToken + "\"")
  
  client := &http.Client{}
  
  res, err := client.Do(req)
  
  if err != nil {
    return "", err
  }
  
  defer res.Body.Close()
  
  if res.StatusCode != 201 {
    return "", fmt.Errorf("HTTP status code was %d", res.StatusCode)
  }
  
  bodyBytes, err := io.ReadAll(res.Body)
  
  if err != nil {
    return "", err
  }
  
  bodyString := string(bodyBytes)
  
  return bodyString, nil
}
