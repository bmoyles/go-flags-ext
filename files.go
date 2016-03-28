package flagext

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// equivalent of ioutil.NopCloser() for writers so we can avoid closing os.Stdout
type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }

// NopWriteCloser wraps an io.Writer in a struct that exposes a no-op Close() method and returns it as an io.WriteCloser
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return nopWriteCloser{w}
}

// InputFile, after argument parsing, should contain the requested filename opened for reading
type InputFile struct {
	*os.File
}

// MasrshalFlag returns the string representation of the InputFile (the file name) or an error if there was
// a problem calling Stat() on the file
func (f *InputFile) MarshalFlag() (string, error) {
	if _, err := f.Stat(); err != nil {
		return "", err
	}
	return f.Name(), nil
}

// UnmarshalFlag attempts to return a file opened for reading with the input string value used as the filename
func (f *InputFile) UnmarshalFlag(value string) error {
	f.File, err = os.Open(ExpandUser(value))
	return err
}

// OutputFile is a go-flags option type that, after argument parsing, should
// contain the requested filename opened for writing. The file is opened with
// os.O_WRONLY|os.O_CREATE|os.O_TRUNC and FileMode 0600.
type OutputFile struct {
	*os.File
}

// UnmasrshalFlag returns the string representation of the OutputFile (the file name) or an error if there was
// a problem calling Stat() on the file
func (f *OutputFile) MarshalFlag() (string, error) {
	if _, err := f.Stat(); err != nil {
		return "", err
	}
	return f.Name(), nil
}

// UnmarshalFlag attempts to return a file opened for writing with the input string value used as the filename
func (f *OutputFile) UnmarshalFlag(value string) error {
	f.File, err = os.OpenFile(ExpandUser(value), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	return err
}

// MultiSourceString is a go-flags option type that, depending on the presence of a
// prefix in the value, may read the actual value from another location.
// Value file:<filename> returns the contents of the file specified by <filename> as a string
// Value env:<environment_variable> returns the contents of the environment variable specified by <environment_variable> as a string
// Values with no prefix are returned as strings, unmanipulated
type MultiSourceString struct {
	s string
}

func (s *MultiSourceString) String() string {
	return s.s
}

func (s *MultiSourceString) MarshalFlag() (string, error) {
	return s.s, nil
}

func (s *MultiSourceString) UnmarshalFlag(value string) error {
	switch {
	case value == "":
		return "", fmt.Errorf("Received empty string, value must be provided")
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
			return "", fmt.Errorf("Environment variable %s empty or non-existent", envVar)
		}
		s.s = envVal
		nil
	default:
		s.s = value
		return nil
	}
}
