package pushpad

import (
  "fmt"
)

type Notification struct {
  Body string
  Title string
  TargetUrl string
  IconUrl string
  BadgeUrl string
  ImageUrl string
  Ttl int
  RequireInteraction bool
  Silent bool
  Urgent bool
  CustomData string
  CustomMetrics []string
  Starred bool
}

func (n Notification) Broadcast() {
  fmt.Printf("%+v\n", n)
}
