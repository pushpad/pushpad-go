package pushpad

var pushpadAuthToken string
var pushpadProjectID int64

// Configure sets the global credentials for API calls.
func Configure(authToken string, projectID int64) {
	pushpadAuthToken = authToken
	pushpadProjectID = projectID
}
