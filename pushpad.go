package pushpad

var pushpadAuthToken string
var pushpadProjectID string

func Configure (authToken string, projectID string) {
  pushpadAuthToken = authToken
  pushpadProjectID = projectID
}
