package notification

// Actions returns a pointer to a slice of actions for optional payload fields.
func Actions(actions ...NotificationAction) *[]NotificationAction {
	return &actions
}
