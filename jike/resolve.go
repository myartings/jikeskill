package jike

import (
	"fmt"
	"net/http"
	"strings"
)

// ResolveShortURL resolves a Jike short URL (e.g., https://okjk.co/xxx) to the final URL
// and extracts the username. Returns the username found in the redirect chain.
// Also accepts bare short codes like "rAgUmv" (auto-prefixed with https://okjk.co/).
func ResolveShortURL(rawURL string) (username string, err error) {
	// Handle bare short codes (4-10 alphanumeric, mixed case)
	if len(rawURL) >= 4 && len(rawURL) <= 10 && !strings.Contains(rawURL, "/") && !strings.Contains(rawURL, ".") {
		hasUpper, hasLower := false, false
		allAlnum := true
		for _, c := range rawURL {
			if c >= 'A' && c <= 'Z' {
				hasUpper = true
			} else if c >= 'a' && c <= 'z' {
				hasLower = true
			} else if c < '0' || c > '9' {
				allAlnum = false
			}
		}
		if allAlnum && hasUpper && hasLower {
			rawURL = "https://okjk.co/" + rawURL
		}
	}
	// Add https:// if missing
	if !strings.Contains(rawURL, "://") {
		rawURL = "https://" + rawURL
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // don't follow redirects automatically
		},
	}

	currentURL := rawURL
	for i := 0; i < 10; i++ { // max 10 redirects
		resp, err := client.Get(currentURL)
		if err != nil {
			return "", fmt.Errorf("resolve URL: %w", err)
		}
		resp.Body.Close()

		// Check if current URL contains a username pattern
		if u := extractUsername(currentURL); u != "" {
			return u, nil
		}

		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			location := resp.Header.Get("Location")
			if location == "" {
				break
			}
			currentURL = location
			continue
		}

		// Final destination, check it
		if u := extractUsername(currentURL); u != "" {
			return u, nil
		}
		break
	}

	return "", fmt.Errorf("could not extract username from URL: %s", rawURL)
}

// extractUsername tries to extract a Jike username from a URL.
// Known patterns:
//   - https://web.okjike.com/u/{username}
//   - https://m.okjike.com/users/{username}
//   - https://www.okjike.com/users/{username}
func extractUsername(url string) string {
	patterns := []string{
		"okjike.com/u/",
		"okjike.com/users/",
	}
	for _, p := range patterns {
		idx := strings.Index(url, p)
		if idx < 0 {
			continue
		}
		rest := url[idx+len(p):]
		// Take until next / or ? or end
		for i, c := range rest {
			if c == '/' || c == '?' || c == '#' {
				rest = rest[:i]
				break
			}
		}
		if rest != "" {
			return rest
		}
	}
	return ""
}
