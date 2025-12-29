package pushpad

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// SignatureFor generates the HMAC signature for a user ID using the configured token.
func SignatureFor(uid string) string {
	h := hmac.New(sha256.New, []byte(pushpadAuthToken))
	h.Write([]byte(uid))
	return hex.EncodeToString(h.Sum(nil))
}
