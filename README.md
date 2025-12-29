# Pushpad - Web Push Notifications

[![Go Reference](https://pkg.go.dev/badge/github.com/pushpad/pushpad-go)](https://pkg.go.dev/github.com/pushpad/pushpad-go)
![Build Status](https://github.com/pushpad/pushpad-go/workflows/CI/badge.svg)
 
[Pushpad](https://pushpad.xyz) is a service for sending push notifications from websites and web apps. It uses the **Push API**, which is a standard supported by all major browsers (Chrome, Firefox, Opera, Edge, Safari).

The notifications are delivered in real time even when the users are not on your website and you can target specific users or send bulk notifications.

## Installation

You can get the Go module:

```bash
go get github.com/pushpad/pushpad-go
```

Then import the packages:

```go
import (
  "github.com/pushpad/pushpad-go"
  "github.com/pushpad/pushpad-go/notification"
  "github.com/pushpad/pushpad-go/project"
  "github.com/pushpad/pushpad-go/sender"
  "github.com/pushpad/pushpad-go/subscription"
)
```

## Getting started

First you need to sign up to Pushpad and create a project there.

Then set your authentication credentials and project:

```go
pushpad.Configure("AUTH_TOKEN", 123)
```

- `AUTH_TOKEN` can be found in the user account settings.
- `PROJECT_ID` can be found in the project settings. If your application uses multiple projects, you can pass the `ProjectID` as a param to functions.

```go
res, err := notification.Create(&notification.NotificationCreateParams{
  ProjectID: pushpad.Int64(123),
  Body: pushpad.String("Your message"),
})

notifications, err := notification.List(&notification.NotificationListParams{
  ProjectID: pushpad.Int64(123),
  Page: pushpad.Int64(1),
})

// ...
```

## Collecting user subscriptions to push notifications

You can subscribe the users to your notifications using the Javascript SDK, as described in the [getting started guide](https://pushpad.xyz/docs/pushpad_pro_getting_started).

If you need to generate the HMAC signature for the `uid` you can use this helper:

```go
s := pushpad.SignatureFor("CURRENT_USER_ID")
fmt.Printf("User ID Signature: %s", s)
```

## Sending push notifications

Use `notification.Create()` (or the `Send()` alias) to create and send a notification:

```go
n := notification.NotificationCreateParams{
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
  TTL: pushpad.Int64(604800),

  // optional, prevent Chrome on desktop from automatically closing the notification after a few seconds
  RequireInteraction: pushpad.Bool(true),

  // optional, enable this option if you want a mute notification without any sound
  Silent: pushpad.Bool(false),

  // optional, enable this option only for time-sensitive alerts (e.g. incoming phone call)
  Urgent: pushpad.Bool(false),

  // optional, a string that is passed as an argument to action button callbacks
  CustomData: pushpad.String("123"),

  // optional, add some action buttons to the notification
  // see https://pushpad.xyz/docs/action_buttons
  Actions: &[]notification.NotificationActionParams{
    {
      Title: pushpad.String("My Button 1"),
      TargetURL: pushpad.String("https://example.com/button-link"), // optional
      Icon: pushpad.String("https://example.com/assets/button-icon.png"), // optional
      Action: pushpad.String("myActionName"), // optional
    },
  },

  // optional, bookmark the notification in the Pushpad dashboard (e.g. to highlight manual notifications)
  Starred: pushpad.Bool(true),

  // optional, use this option only if you need to create scheduled notifications (max 5 days)
  // see https://pushpad.xyz/docs/schedule_notifications
  SendAt: pushpad.Time(sendAtTime), // sendAtTime := time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC)

  // optional, add the notification to custom categories for stats aggregation
  // see https://pushpad.xyz/docs/monitoring
  CustomMetrics: pushpad.StringSlice([]string{"examples", "another_metric"}), // up to 3 metrics per notification
}

res, err := notification.Create(&n)

// TARGETING:
// You can use UIDs and Tags for sending the notification only to a specific audience...

// deliver to a user
n := notification.NotificationCreateParams{
  Body: pushpad.String("Hi user1"),
  UIDs: pushpad.StringSlice([]string{"user1"}),
}
res, err := notification.Create(&n)

// deliver to a group of users
n := notification.NotificationCreateParams{
  Body: pushpad.String("Hi users"),
  UIDs: pushpad.StringSlice([]string{"user1","user2","user3"}),
}
res, err := notification.Create(&n)

// deliver to some users only if they have a given preference
// e.g. only "users" who have a interested in "events" will be reached
n := notification.NotificationCreateParams{
  Body: pushpad.String("New event"),
  UIDs: pushpad.StringSlice([]string{"user1","user2"}),
  Tags: pushpad.StringSlice([]string{"events"}),
}
res, err := notification.Create(&n)

// deliver to segments
// e.g. any subscriber that has the tag "segment1" OR "segment2"
n := notification.NotificationCreateParams{
  Body: pushpad.String("Example"),
  Tags: pushpad.StringSlice([]string{"segment1", "segment2"}),
}
res, err := notification.Create(&n)

// you can use boolean expressions
// they can include parentheses and the operators !, &&, || (from highest to lowest precedence)
// https://pushpad.xyz/docs/tags
n := notification.NotificationCreateParams{
  Body: pushpad.String("Example"),
  Tags: pushpad.StringSlice([]string{"zip_code:28865 && !optout:local_events || friend_of:Organizer123"}),
}
res, err := notification.Create(&n)

n := notification.NotificationCreateParams{
  Body: pushpad.String("Example"),
  Tags: pushpad.StringSlice([]string{"tag1 && tag2", "tag3"}), // equal to 'tag1 && tag2 || tag3'
}
res, err := notification.Create(&n)

// deliver to everyone
n := notification.NotificationCreateParams{
  Body: pushpad.String("Hello everybody"),
}
res, err := notification.Create(&n)
```

You can set the default values for most fields in the project settings. See also [the docs](https://pushpad.xyz/docs/rest_api#notifications_api_docs) for more information about notification fields.

If you try to send a notification to a user ID, but that user is not subscribed, that ID is simply ignored.

These fields are returned by the API:

```go
res, err := notification.Create(&n)

// Notification ID
fmt.Println(res.ID) // => 1000

// Estimated number of devices that will receive the notification
// Not available for notifications that use SendAt
fmt.Println(res.Scheduled) // => 5

// Available only if you specify some user IDs (UIDs) in the request:
// it indicates which of those users are subscribed to notifications.
// Not available for notifications that use SendAt
fmt.Println(res.UIDs) // => []string{"user1", "user2"}

// The time when the notification will be sent.
// Available for notifications that use SendAt
fmt.Println(res.SendAt) // => 2025-10-30 10:09:00 +0000 UTC

// Note:
// when a field is not available in the response, it is set to its zero value
```

## Getting push notification data

You can retrieve data for past notifications:

```go
n, err := notification.Get(42, nil)

// get basic attributes
fmt.Println(n.ID) // => 42
fmt.Println(n.Title) // => Foo Bar
fmt.Println(n.Body) // => Lorem ipsum dolor sit amet, consectetur adipiscing elit.
fmt.Println(n.TargetURL) // => https://example.com
fmt.Println(n.TTL) // => 604800
fmt.Println(n.RequireInteraction) // => false
fmt.Println(n.Silent) // => false
fmt.Println(n.Urgent) // => false
fmt.Println(n.IconURL) // => https://example.com/assets/icon.png
fmt.Println(n.BadgeURL) // => https://example.com/assets/badge.png
fmt.Println(n.CreatedAt) // => 2025-07-06 10:09:14 +0000 UTC

// get statistics
fmt.Println(n.ScheduledCount) // => 1
fmt.Println(n.SuccessfullySent) // => 4
fmt.Println(n.OpenedCount) // => 2
```

Or for multiple notifications of a project at once:

```go
notifications, err := notification.List(&notification.NotificationListParams{
  Page: pushpad.Int64(1),
})

// same attributes as for single notification in example above
fmt.Println(notifications[0].ID) // => 42
fmt.Println(notifications[0].Title) // => Foo Bar
```

The REST API paginates the result set. You can pass a `Page` parameter to get the full list in multiple requests.

```go
notifications, err := notification.List(&notification.NotificationListParams{
  Page: pushpad.Int64(2),
})
```

## Scheduled notifications

You can create scheduled notifications that will be sent in the future:

```go
sendAt := time.Now().UTC().Add(60 * time.Second)

scheduled, err := notification.Create(&notification.NotificationCreateParams{
  Body: pushpad.String("This notification will be sent after 60 seconds"),
  SendAt: pushpad.Time(sendAt),
})
```

You can also cancel a scheduled notification:

```go
err := notification.Cancel(scheduled.ID, nil)
```

## Getting subscription count

You can retrieve the number of subscriptions for a given project, optionally filtered by `Tags` or `UIDs`:

```go
totalCount, err := subscription.Count(&subscription.SubscriptionCountParams{})
fmt.Println(totalCount) // => 100

totalCount, err = subscription.Count(&subscription.SubscriptionCountParams{
  UIDs: pushpad.StringSlice([]string{"user1"}),
})
fmt.Println(totalCount) // => 2

totalCount, err = subscription.Count(&subscription.SubscriptionCountParams{
  Tags: pushpad.StringSlice([]string{"sports"}),
})
fmt.Println(totalCount) // => 10

totalCount, err = subscription.Count(&subscription.SubscriptionCountParams{
  Tags: pushpad.StringSlice([]string{"sports && travel"}),
})
fmt.Println(totalCount) // => 5

totalCount, err = subscription.Count(&subscription.SubscriptionCountParams{
  UIDs: pushpad.StringSlice([]string{"user1"}),
  Tags: pushpad.StringSlice([]string{"sports && travel"}),
})
fmt.Println(totalCount) // => 1
```

## Getting push subscription data

You can retrieve the subscriptions for a given project, optionally filtered by `Tags` or `UIDs`:

```go
subscriptions, err := subscription.List(&subscription.SubscriptionListParams{})

subscriptions, err = subscription.List(&subscription.SubscriptionListParams{
  UIDs: pushpad.StringSlice([]string{"user1"}),
})

subscriptions, err = subscription.List(&subscription.SubscriptionListParams{
  Tags: pushpad.StringSlice([]string{"sports"}),
})

subscriptions, err = subscription.List(&subscription.SubscriptionListParams{
  Tags: pushpad.StringSlice([]string{"sports && travel"}),
})

subscriptions, err = subscription.List(&subscription.SubscriptionListParams{
  UIDs: pushpad.StringSlice([]string{"user1"}),
  Tags: pushpad.StringSlice([]string{"sports && travel"}),
})
```

The REST API paginates the result set. You can pass `Page` and `PerPage` parameters to get the full list in multiple requests.

```go
subscriptions, err := subscription.List(&subscription.SubscriptionListParams{
  Page: pushpad.Int64(2),
})
```

You can also retrieve the data of a specific subscription if you already know its id:

```go
subscription.Get(123, nil)
```

## Updating push subscription data

Usually you add data, like user IDs and tags, to the push subscriptions using the [JavaScript SDK](https://pushpad.xyz/docs/javascript_sdk_reference) in the frontend.

However you can also update the subscription data from your server:

```go
subscriptions, err := subscription.List(&subscription.SubscriptionListParams{
  UIDs: pushpad.StringSlice([]string{"user1"}),
})

for _, s := range subscriptions {
  // update the user ID associated to the push subscription
  _, err = subscription.Update(s.ID, &subscription.SubscriptionUpdateParams{
    UID: pushpad.String("myuser1"),
  })

  // update the tags associated to the push subscription
  tags := append([]string{}, s.Tags...)
  tags = append(tags, "another_tag")
  _, err = subscription.Update(s.ID, &subscription.SubscriptionUpdateParams{
    Tags: pushpad.StringSlice(tags),
  })
}
```

## Importing push subscriptions

If you need to [import](https://pushpad.xyz/docs/import) some existing push subscriptions (from another service to Pushpad, or from your backups) or if you simply need to create some test data, you can use this method:

```go
createdSubscription, err := subscription.Create(&subscription.SubscriptionCreateParams{
  Endpoint: pushpad.String("https://example.com/push/f7Q1Eyf7EyfAb1"),
  P256DH: pushpad.String("BCQVDTlYWdl05lal3lG5SKr3VxTrEWpZErbkxWrzknHrIKFwihDoZpc_2sH6Sh08h-CacUYI-H8gW4jH-uMYZQ4="),
  Auth: pushpad.String("cdKMlhgVeSPzCXZ3V7FtgQ=="),
  UID: pushpad.String("exampleUid"),
  Tags: pushpad.StringSlice([]string{"exampleTag1", "exampleTag2"}),
})
```

Please note that this is not the standard way to collect subscriptions on Pushpad: usually you subscribe the users to the notifications using the [JavaScript SDK](https://pushpad.xyz/docs/javascript_sdk_reference) in the frontend.

## Deleting push subscriptions

Usually you unsubscribe a user from push notifications using the [JavaScript SDK](https://pushpad.xyz/docs/javascript_sdk_reference) in the frontend (recommended).

However you can also delete the subscriptions using this library. Be careful, the subscriptions are permanently deleted!

```go
err := subscription.Delete(id, nil)
```

## Managing projects

Projects are usually created manually from the Pushpad dashboard. However you can also create projects from code if you need advanced automation or if you manage [many different domains](https://pushpad.xyz/docs/multiple_domains).

```go
createdProject, err := project.Create(&project.ProjectCreateParams{
  // required attributes
  SenderID: pushpad.Int64(123),
  Name: pushpad.String("My project"),
  Website: pushpad.String("https://example.com"),

  // optional configurations
  IconURL: pushpad.String("https://example.com/icon.png"),
  BadgeURL: pushpad.String("https://example.com/badge.png"),
  NotificationsTTL: pushpad.Int64(604800),
  NotificationsRequireInteract: pushpad.Bool(false),
  NotificationsSilent: pushpad.Bool(false),
})
```

You can also find, update and delete projects:

```go
projects, err := project.List(nil)
for _, p := range projects {
  fmt.Printf("Project %d: %s\n", p.ID, p.Name)
}

existingProject, err := project.Get(123, nil)

updatedProject, err := project.Update(existingProject.ID, &project.ProjectUpdateParams{
  Name: pushpad.String("The New Project Name"),
})

err = project.Delete(existingProject.ID, nil)
```

## Managing senders

Senders are usually created manually from the Pushpad dashboard. However you can also create senders from code.

```go
createdSender, err := sender.Create(&sender.SenderCreateParams{
  // required attributes
  Name: pushpad.String("My sender"),

  // optional configurations
  // do not include these fields if you want to generate them automatically
  VAPIDPrivateKey: pushpad.String("-----BEGIN EC PRIVATE KEY----- ..."),
  VAPIDPublicKey: pushpad.String("-----BEGIN PUBLIC KEY----- ..."),
})
```

You can also find, update and delete senders:

```go
senders, err := sender.List(nil)
for _, s := range senders {
  fmt.Printf("Sender %d: %s\n", s.ID, s.Name)
}

existingSender, err := sender.Get(987, nil)

updatedSender, err := sender.Update(existingSender.ID, &sender.SenderUpdateParams{
  Name: pushpad.String("The New Sender Name"),
})

err = sender.Delete(existingSender.ID, nil)
```

## Error handling

API requests can return errors, described by a `pushpad.APIError` that exposes the HTTP status code and response body. Network issues and other errors return a generic error.

```go
n := notification.NotificationCreateParams{
  Body: pushpad.String("Hello"),
}
_, err := notification.Create(&n)
if err != nil {
  var apiErr *pushpad.APIError
  if errors.As(err, &apiErr) { // HTTP error from the API
    fmt.Println(apiErr.StatusCode, apiErr.Body)
  } else { // network error or other errors
    fmt.Println(err)
  }
}
```

## Documentation

- Pushpad REST API reference: https://pushpad.xyz/docs/rest_api
- Getting started guide (for collecting subscriptions): https://pushpad.xyz/docs/pushpad_pro_getting_started
- JavaScript SDK reference (frontend): https://pushpad.xyz/docs/javascript_sdk_reference

## License

The library is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).
