# Pushpad - Web Push Notifications
 
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
)
```

## Getting started

First you need to sign up to Pushpad and create a project there.

Then set your authentication credentials and project:

```go
pushpad.Configure("AUTH_TOKEN", "PROJECT_ID")
```

- `AUTH_TOKEN` can be found in the user account settings.
- `PROJECT_ID` can be found in the project settings. If your application uses multiple projects, you can set the `ProjectID` as an additional field for `Notification`.

## Collecting user subscriptions to push notifications

You can subscribe the users to your notifications using the Javascript SDK, as described in the [getting started guide](https://pushpad.xyz/docs/pushpad_pro_getting_started).

If you need to generate the HMAC signature for the `uid` you can use this helper:

```go
s := pushpad.SignatureFor("CURRENT_USER_ID")
fmt.Printf("User ID Signature: %s", s)
```

## Sending push notifications

```go
n := pushpad.Notification {
  Body: "Hello world!",
  Title: "Website Name", // optional, defaults to your project name
  TargetURL: "https://example.com", // optional, defaults to your project website
  IconURL: "https://example.com/assets/icon.png", // optional, defaults to the project icon
  BadgeURL: "https://example.com/assets/badge.png", // optional, defaults to the project badge
  ImageURL: "https://example.com/assets/image.png", // optional, an image to display in the notification content
  TTL: 604800, // optional, drop the notification after this number of seconds if a device is offline
  RequireInteraction: true, // optional, prevent Chrome on desktop from automatically closing the notification after a few seconds
  Silent: false, // optional, enable this option if you want a mute notification without any sound
  Urgent: false, // optional, enable this option only for time-sensitive alerts (e.g. incoming phone call)
  CustomData: "123", // optional, a string that is passed as an argument to action button callbacks
  Starred: true, // optional, bookmark the notification in the Pushpad dashboard (e.g. to highlight manual notifications)
  // optional, use this option only if you need to create scheduled notifications (max 5 days)
  // see https://pushpad.xyz/docs/schedule_notifications
  SendAt: &sendAtTime, // sendAtTime := time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC)
  // optional, add the notification to custom categories for stats aggregation
  // see https://pushpad.xyz/docs/monitoring
  CustomMetrics: []string{"examples", "another_metric"}, // up to 3 metrics per notification
}

res, err := n.Send()

// TARGETING:
// You can use UIDs and Tags for sending the notification only to a specific audience...

// deliver to a user
n := pushpad.Notification { Body: "Hi user1", UIDs: []string{"user1"} }
res, err := n.Send()

// deliver to a group of users
n := pushpad.Notification { Body: "Hi users", UIDs: []string{"user1","user2","user3"} }
res, err := n.Send()

// deliver to some users only if they have a given preference
// e.g. only "users" who have a interested in "events" will be reached
n := pushpad.Notification { Body: "New event", UIDs: []string{"user1","user2"}, Tags: []string{"events"} }
res, err := n.Send()

// deliver to segments
// e.g. any subscriber that has the tag "segment1" OR "segment2"
n := pushpad.Notification { Body: "Example", Tags: []string{"segment1", "segment2"} }
res, err := n.Send()

// you can use boolean expressions 
// they must be in the disjunctive normal form (without parenthesis)
n := pushpad.Notification { Body: "Example", Tags: []string{"zip_code:28865 && !optout:local_events || friend_of:Organizer123"} }
res, err := n.Send()

// deliver to everyone
n := pushpad.Notification { Body: "Hello everybody" }
res, err := n.Send()
```

You can set the default values for most fields in the project settings. See also [the docs](https://pushpad.xyz/docs/rest_api#notifications_api_docs) for more information about notification fields.

If you try to send a notification to a user ID, but that user is not subscribed, that ID is simply ignored.

The methods above return a `NotificationResponse struct`:

- `ID` is the id of the notification on Pushpad
- `Scheduled` is the estimated reach of the notification (i.e. the number of devices to which the notification will be sent, which can be different from the number of users, since a user may receive notifications on multiple devices)
- `UIDs` (only when the `UIDs` field is set) are the user IDs that will be actually reached by the notification because they are subscribed to your notifications. For example if you send a notification to `{"uid1", "uid2", "uid3"}`, but only `"uid1"` is subscribed, you will get `{"uid1"}` in response. Note that if a user has unsubscribed after the last notification sent to him, he may still be reported for one time as subscribed (this is due to the way the W3C Push API works).
- `SendAt` is present only for scheduled notifications. The fields `Scheduled` and `UIDs` are not available in this case.


## License

The library is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).
