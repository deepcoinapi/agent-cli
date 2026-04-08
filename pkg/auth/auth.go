// Package auth provides HMAC-SHA256 signing for the DeepCoin API.
//
// Signature = Base64(HMAC-SHA256(timestamp + method + requestPath + body, secretKey))
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

// Timestamp returns the current UTC time in ISO 8601 format for DC-ACCESS-TIMESTAMP.
func Timestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
}

// Sign computes the HMAC-SHA256 signature for DC-ACCESS-SIGN.
func Sign(timestamp, method, requestPath, body, secretKey string) string {
	message := timestamp + method + requestPath + body
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
