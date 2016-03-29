package flagext

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
)

// MultiSourceString, when unmarshaled, whose contents vary depending on prefixes
// on the input value.
// Provided value file:<filename> will cause the contents of <filename> to be
// read as the actual value.
// Provided value env:<environment_variable> will cause the contents of the
// environment variable specified by <environment_variable> to be read as the actual value.
// Provided values with no prefix are returned as strings, unmanipulated.
type MultiSourceString struct {
	s string
}

// String returns the internal string value after parsing.
func (s *MultiSourceString) String() string {
	return s.s
}

// Bytes returns the internal string value as a byte slice after parsing.
func (s *MultiSourceString) Bytes() []byte {
	return []byte(s.s)
}

// MarshalFlag is provided to satisfy the flags.Marshaler interface and simply returns
// the internal string value.
func (s *MultiSourceString) MarshalFlag() (string, error) {
	return s.s, nil
}

// UnmarshalFlag is provided to satisfy flags.Unmarshaler and takes the string value
// from go-flags parsing and, depending on prefix, sets the internal string value to
// the appropriate contents.
func (s *MultiSourceString) UnmarshalFlag(value string) error {
	switch {
	case value == "":
		return fmt.Errorf("Received empty string, value must be provided")
	case strings.HasPrefix(value, "file:"):
		filename := ExpandUser(strings.TrimPrefix(value, "file:"))
		val, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		s.s = string(val)
		return nil
	case strings.HasPrefix(value, "env:"):
		envVar := strings.TrimPrefix(value, "env:")
		envVal, ok := os.LookupEnv(envVar)
		if (!ok) || (envVal == "") {
			return fmt.Errorf("Environment variable %s empty or non-existent", envVar)
		}
		s.s = envVal
		return nil
	default:
		s.s = value
		return nil
	}
}

// Complete satisfies flags.Completer and returns potential matches based on current input.
// If file: or env: prefixes are detected, files or environment variables are completed.
// Non-prefixed input will return no completions.
func (s *MultiSourceString) Complete(match string) []flags.Completion {
	switch {
	case strings.HasPrefix(match, "file:"):
		filename := ExpandUser(strings.TrimPrefix(match, "file:"))
		items, _ := filepath.Glob(filename + "*")
		completions := make([]flags.Completion, len(items))
		for i, v := range items {
			completions[i].Item = v
		}
		return completions
	case strings.HasPrefix(match, "env:"):
		envVar := strings.TrimPrefix(match, "env:")
		items := os.Environ()
		var completions []flags.Completion
		for _, item := range items {
			parts := strings.Split(item, "=")
			if strings.HasPrefix(parts[0], envVar) {
				completions = append(completions, flags.Completion{Item: parts[0]})
			}
		}
		return completions
	default:
		return []flags.Completion{}
	}
}
