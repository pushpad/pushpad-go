# Upgrading to version 1.x

This version is a major rewrite of the library and adds support for the full REST API, including Notifications, Subscriptions, Projects and Senders.

This version has some breaking changes:

- When you call `pushpad.Configure` the `projectID` argument is now a `int64` instead of a `string`.
- `pushpad.Notification` is now used only for some API responses, but not for API requests. If you want to create / send a notification, use `notification.Create(&notificationCreateParams)`.
- All fields in `notification.NotificationCreateParams` are pointers, and you can use helpers like `pushpad.String`, `pushpad.StringSlice`, `pushpad.Int64`, etc. to create the pointers easily. For example, when you create the params for a notification, use `Body: pushpad.String("Hello")` instead of `Body: "Hello"`.
- The response to the creation of a notification is now a `NotificationCreateResponse struct` instead of `NotificationResponse struct` (the only difference is the `struct` name and the use of `int64` instead of `int` for some fields).
