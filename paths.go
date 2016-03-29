package flagext

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// ExpandUser receives a path as a string and attempts to expand a leading tilde
// or leading tilde + username to the actual user home directory in case it is not
// expanded by the shell first. If the home directory cannot be expanded, the path
// is returned as-is.
//
// Example
//
// If user foo does exist, ExpandUser("~foo/bar") should return "<path to foo home directory>/bar"
// If user foo does NOT exist, ExpandUser("~foo/bar") will return "~foo/bar" as-is
func ExpandUser(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}
	parts := strings.Split(path, string(os.PathSeparator))
	userPart := parts[0]

	if userPart == "~" {
		u, err := user.Current()
		if err != nil {
			return path
		}
		parts[0] = u.HomeDir
	} else {
		userPart = strings.TrimPrefix(userPart, "~")
		u, err := user.Lookup(userPart)
		if err != nil {
			return path
		}
		parts[0] = u.HomeDir
	}
	return filepath.Join(parts...)
}
