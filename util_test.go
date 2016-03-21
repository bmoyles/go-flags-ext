package flagtypes

import (
	"os/user"
	"path/filepath"
	"testing"
)

func TestExpandUser(t *testing.T) {
	u, err := user.Current()
	if err != nil {
		t.Error("Unable to get current user: ", err)
	}
	testData := []struct{ provided, expected string }{
		{provided: "~/foo/bar", expected: filepath.Join(u.HomeDir, "foo/bar")},
		{provided: "~" + u.Username + "/foo/bar", expected: filepath.Join(u.HomeDir, "foo/bar")},
		{provided: "~somefakeuser/foo/bar", expected: "~somefakeuser/foo/bar"},
	}
	for _, item := range testData {
		if expandUser(item.provided) != item.expected {
			t.Errorf("Provided string (%s) does not match expected string (%s)", item.provided, item.expected)
		}
	}

}
