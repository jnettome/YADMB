package manager

import (
	"os"
	"sync"
)

// Default cookies path next to the YADMB working directory (production: /opt/yadmb/cookies.txt).
const defaultCookiesPath = "cookies.txt"

var (
	cookieArgsOnce sync.Once
	cookieArgs     []string
)

// ytDlpCookieArgs returns ["--cookies", path] when a Netscape cookies file is present.
// Override with YADMB_COOKIES=/path/to/cookies.txt.
func ytDlpCookieArgs() []string {
	cookieArgsOnce.Do(func() {
		path := os.Getenv("YADMB_COOKIES")
		if path == "" {
			path = defaultCookiesPath
		}
		if st, err := os.Stat(path); err == nil && !st.IsDir() {
			cookieArgs = []string{"--cookies", path}
		}
	})
	return cookieArgs
}
