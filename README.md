# Pushpad - Web Push Notifications

[![Go Reference](https://pkg.go.dev/badge/github.com/pushpad/pushpad-go)](https://pkg.go.dev/github.com/pushpad/pushpad-go)
![Build Status](https://github.com/pushpad/pushpad-go/workflows/CI/badge.svg)
 
[Pushpad](https://pushpad.xyz) is a service for sending push notifications from websites and web apps. It uses the **Push API**, which is a standard supported by all major browsers (Chrome, Firefox, Opera, Edge, Safari).

The notifications are delivered in real time even when the users are not on your website and you can target specific users or send bulk notifications.

## Installation

You can get the Go module:

```go
go get github.com/pushpad/pushpad-go
```

Then import the packages:

```go
import (
  "github.com/pushpad/pushpad-go"
  "github.com/pushpad/pushpad-go/notification"
)
```

Other resources are available in the `subscription`, `project`, and `sender` packages.

## Getting started

First you need to sign up to Pushpad and create a project there.

Then set your authentication credentials and project:

```go
pushpad.Configure("AUTH_TOKEN", 123)
```

- `AUTH_TOKEN` can be found in the user account settings.
- `PROJECT_ID` can be found in the project settings. If your application uses multiple projects, you can pass the `ProjectID` as a param to functions.

## Collecting user subscriptions to push notifications

You can subscribe the users to your notifications using the Javascript SDK, as described in the [getting started guide](https://pushpad.xyz/docs/pushpad_pro_getting_started).

If you need to generate the HMAC signature for the `uid` you can use this helper:

```go
s := pushpad.SignatureFor("CURRENT_USER_ID")
fmt.Printf("User ID Signature: %s", s)
```

## Sending push notifications

```go
n := notification.NotificationCreateParams {
  // optional, defaults to the project configured via pushpad.Configure
  ProjectID: pushpad.Int(0),

  // required, the main content of the notification
  Body: pushpad.String("Hello world!"),

  // optional, the title of the notification (defaults to your project name)
  Title: pushpad.String("Website Name"),

  // optional, open this link on notification click (defaults to your project website)
  TargetURL: pushpad.String("https://example.com"),

  // optional, the icon of the notification (defaults to the project icon)
  IconURL: pushpad.String("https://example.com/assets/icon.png"),

  // optional, the small icon displayed in the status bar (defaults to the project badge)
  BadgeURL: pushpad.String("https://example.com/assets/badge.png"),

  // optional, an image to display in the notification content
  // see https://pushpad.xyz/docs/sending_images
  ImageURL: pushpad.String("https://example.com/assets/image.png"),

  // optional, drop the notification after this number of seconds if a device is offline
  TTL: pushpad.Int(604800),

  // optional, prevent Chrome on desktop from automatically closing the notification after a few seconds
  RequireInteraction: pushpad.Bool(true),

  // optional, enable this option if you want a mute notification without any sound
  Silent: pushpad.Bool(false),

  // optional, enable this option only for time-sensitive alerts (e.g. incoming phone call)
  Urgent: pushpad.Bool(false),

  // optional, a string that is passed as an argument to action button callbacks
  CustomData: pushpad.String("123"),

  // optional, bookmark the notification in the Pushpad dashboard (e.g. to highlight manual notifications)
  Starred: pushpad.Bool(true),

  // optional, use this option only if you need to create scheduled notifications (max 5 days)
  // see https://pushpad.xyz/docs/schedule_notifications
  SendAt: pushpad.Time(sendAtTime), // sendAtTime := time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC)

  // optional, add the notification to custom categories for stats aggregation
  // see https://pushpad.xyz/docs/monitoring
  CustomMetrics: pushpad.Strings([]string{"examples", "another_metric"}), // up to 3 metrics per notification
}

res, err := notification.Create(&n)

// TARGETING:
// You can use UIDs and Tags for sending the notification only to a specific audience...

// deliver to a user
n := notification.NotificationCreateParams { Body: pushpad.String("Hi user1"), UIDs: pushpad.Strings([]string{"user1"}) }
res, err := notification.Create(&n)

// deliver to a group of users
n := notification.NotificationCreateParams { Body: pushpad.String("Hi users"), UIDs: pushpad.Strings([]string{"user1","user2","user3"}) }
res, err := notification.Create(&n)

// deliver to some users only if they have a given preference
// e.g. only "users" who have a interested in "events" will be reached
n := notification.NotificationCreateParams { Body: pushpad.String("New event"), UIDs: pushpad.Strings([]string{"user1","user2"}), Tags: pushpad.Strings([]string{"events"}) }
res, err := notification.Create(&n)

// deliver to segments
// e.g. any subscriber that has the tag "segment1" OR "segment2"
n := notification.NotificationCreateParams { Body: pushpad.String("Example"), Tags: pushpad.Strings([]string{"segment1", "segment2"}) }
res, err := notification.Create(&n)

// you can use boolean expressions
// they can include parentheses and the operators !, &&, || (from highest to lowest precedence)
// https://pushpad.xyz/docs/tags
n := notification.NotificationCreateParams { Body: pushpad.String("Example"), Tags: pushpad.Strings([]string{"zip_code:28865 && !optout:local_events || friend_of:Organizer123"}) }
res, err := notification.Create(&n)

// deliver to everyone
n := notification.NotificationCreateParams { Body: pushpad.String("Hello everybody") }
res, err := notification.Create(&n)
```

You can set the default values for most fields in the project settings. See also [the docs](https://pushpad.xyz/docs/rest_api#notifications_api_docs) for more information about notification fields.

If you try to send a notification to a user ID, but that user is not subscribed, that ID is simply ignored.

The methods above return a `NotificationCreateResponse struct`:

- `ID` is the id of the notification on Pushpad
- `Scheduled` is the estimated reach of the notification (i.e. the number of devices to which the notification will be sent, which can be different from the number of users, since a user may receive notifications on multiple devices)
- `UIDs` (only when the `UIDs` field is set) are the user IDs that will be actually reached by the notification because they are subscribed to your notifications. For example if you send a notification to `{"uid1", "uid2", "uid3"}`, but only `"uid1"` is subscribed, you will get `{"uid1"}` in response. Note that if a user has unsubscribed after the last notification sent to him, he may still be reported for one time as subscribed (this is due to the way the W3C Push API works).
- `SendAt` is present only for scheduled notifications. The fields `Scheduled` and `UIDs` are not available in this case.


## License

The library is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).
