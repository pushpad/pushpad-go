package pushpad

var pushpadAuthToken string
var pushpadProjectID int

// Configure sets the global credentials for API calls.
func Configure(authToken string, projectID int) {
	pushpadAuthToken = authToken
	pushpadProjectID = projectID
}
