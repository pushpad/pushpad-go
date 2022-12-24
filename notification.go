package pushpad

import (
  "fmt"
  "bytes"
  "io"
  "encoding/json"
  "net/http"
  "time"
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
  SendAt *time.Time `json:"send_at,omitempty"`
  Uids []string `json:"uids"`
  Tags []string `json:"tags"`
}

type NotificationResponse struct {
  ID int `json:"id"`
  Scheduled int `json:"scheduled"`
  Uids []string `json:"uids"`
  SendAt time.Time `json:"send_at"`
}

func (n Notification) Send() (*NotificationResponse, error) {
  if n.ProjectID == "" {
    n.ProjectID = pushpadProjectID
  }
  
  notificationJSON, err := json.Marshal(n)
  
  if err != nil {
    return nil, err
  }
  
  req, err := http.NewRequest("POST", "https://pushpad.xyz/api/v1/projects/" + n.ProjectID + "/notifications", bytes.NewBuffer(notificationJSON))
  
  if err != nil {
    return nil, err
  }
  
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Authorization", "Token token=\"" + pushpadAuthToken + "\"")
  
  client := &http.Client{}
  
  res, err := client.Do(req)
  
  if err != nil {
    return nil, err
  }
  
  defer res.Body.Close()
  
  bodyBytes, err := io.ReadAll(res.Body)
  
  if err != nil {
    return nil, err
  }
  
  bodyString := string(bodyBytes)
  
  if res.StatusCode != 201 {
    return nil, fmt.Errorf("Response was HTTP %d: %s", res.StatusCode, bodyString)
  }
  
  var r *NotificationResponse
  json.Unmarshal(bodyBytes, &r)
    
  return r, nil
}
