package device

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mssola/user_agent"
)

// Parse extracts a friendly, human-readable device and browser description from the HTTP request.
func Parse(r *http.Request) string {
	rawUA := r.UserAgent()
	if rawUA == "" {
		return "Unknown Device"
	}

	ua := user_agent.New(rawUA)
	browserName, _ := ua.Browser()
	os := ua.OS()

	// Return `Unidentified`, if failed to identify both browser and os
	if browserName == "" && os == "" {
		return "Unidentified"
	}

	if browserName == "" {
		browserName = "Unknown Browser"
	}
	if os == "" {
		os = "Unknown OS"
	}

	// Add a friendly device type hint if it's mobile or tablet
	deviceHint := ""
	if ua.Mobile() {
		deviceHint = " (Mobile)"
	} else if strings.Contains(strings.ToLower(rawUA), "tablet") || strings.Contains(strings.ToLower(rawUA), "ipad") {
		deviceHint = " (Tablet)"
	}

	return fmt.Sprintf("%s on %s%s", browserName, os, deviceHint)
}
