package hypermid

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// VerifyWebhookSignature validates an incoming webhook signature.
//
// When HyperMid sends a webhook, it includes:
//   - X-Hypermid-Signature: HMAC-SHA256 hex digest of the raw body
//   - X-Hypermid-Event: event type (e.g. "swap.completed")
//
// body is the raw request body, signature is the X-Hypermid-Signature header value,
// and secret is the webhook signing secret returned when the webhook was created.
//
// Uses constant-time comparison to prevent timing attacks.
func VerifyWebhookSignature(body []byte, signature string, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
