package flagtypes

import (
	"os"
)

// FileOrStdin is a go-flags option type that, after argument parsing, should
// contain the requested filename opened for reading, or, os.Stdin if the
// requested filename is "-"
type FileOrStdin struct {
	*os.File
}

func (f FileOrStdin) MarshalFlag() (string, error) {
	if _, err := f.Stat(); err != nil {
		return "", err
	}
	return f.Name(), nil
}

func (f *FileOrStdin) UnmarshalFlag(value string) error {
	var err error
	switch value {
	case "":
		return os.ErrInvalid
	case "-":
		f.File = os.Stdin
		return nil
	default:
		f.File, err = os.Open(value)
		return err
	}
}

func (f *FileOrStdin) Close() error {
	if f.File == os.Stdin {
		return nil
	}
	return f.File.Close()
}

// FileOrStdout is a go-flags option type that, after argument parsing, should
// contain the requested filename opened for writing, or, os.Stdout if the
// requested filename is "-". The file, when not os.Stdout, is opened with
// os.O_WRONLY|os.O_CREATE|os.O_TRUNC and FileMode 0600.
type FileOrStdout struct {
	*os.File
}

func (f FileOrStdout) MarshalFlag() (string, error) {
	if _, err := f.Stat(); err != nil {
		return "", err
	}
	return f.Name(), nil
}

func (f *FileOrStdout) UnmarshalFlag(value string) error {
	var err error
	switch value {
	case "":
		return os.ErrInvalid
	case "-":
		f.File = os.Stdout
		return nil
	default:
		f.File, err = os.OpenFile(value, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		return err
	}
}

func (f *FileOrStdout) Close() error {
	if f.File == os.Stdout {
		return nil
	}
	return f.File.Close()
}
