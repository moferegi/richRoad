package client

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func normalizeTokenValue(raw string) string {
	return strings.TrimSpace(raw)
}

func buildTokenFingerprint(raw string) string {
	token := normalizeTokenValue(raw)
	if token == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func maskTokenForDisplay(raw string) string {
	token := normalizeTokenValue(raw)
	if token == "" {
		return ""
	}
	if len(token) <= 8 {
		return strings.Repeat("*", len(token))
	}
	return token[:4] + "****" + token[len(token)-4:]
}

func shortTokenFingerprint(value string) string {
	value = strings.TrimSpace(value)
	if len(value) <= 12 {
		return value
	}
	return value[:12]
}
