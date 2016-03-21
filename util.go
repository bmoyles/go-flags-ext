package flagtypes

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// tilde expansion in path names
func expandUser(path string) string {
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
