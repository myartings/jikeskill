package jike

import (
	"fmt"
	"net/http"
	"strings"
)

// ResolveShortURL resolves a Jike short URL (e.g., https://okjk.co/xxx) to the final URL
// and extracts the username. Returns the username found in the redirect chain.
func ResolveShortURL(rawURL string) (username string, err error) {
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
