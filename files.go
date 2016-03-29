package flagext

import (
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
)

// InputFile, when unmarshaled, embeds an os.File opened for reading.
type InputFile struct {
	*os.File
	// NopClose, if true, overrides the file's Close() method to do nothing
	// similar to ioutil.NopCloser.
	NopClose bool
}

// NewInputFile creates an InputFile from an existing opened os.File.
// Useful for setting defaults via struct initialization rather than struct tags.
func NewInputFile(f *os.File) InputFile {
	return InputFile{File: f}
}

// MasrshalFlag returns the input file name
func (f *InputFile) MarshalFlag() (string, error) {
	return f.Name(), nil
}

// UnmarshalFlag attempts to return a file opened for reading with the input
// string value used as the filename
func (f *InputFile) UnmarshalFlag(value string) error {
	file, err := os.Open(ExpandUser(value))
	if err != nil {
		return err
	}
	*f = InputFile{File: file}
	return nil
}

// Close, if InputFile.NopClose is true, simply returns nil, otherwise the embedded os.File's
// Close is called.
func (f *InputFile) Close() error {
	if f.NopClose {
		return nil
	}
	return f.File.Close()
}

// DefaultToStdin is a convenience method that should be used when setting defaults via
// struct initialization rather than tags so the default file is os.Stdin
func (f *InputFile) DefaultToStdin() {
	*f = InputFile{File: os.Stdin}
	f.NopClose = true
}

// Complete satisfies flags.Completer and returns potential file matches based on current input.
func (f *InputFile) Complete(match string) []flags.Completion {
	items, _ := filepath.Glob(ExpandUser(match) + "*")
	completions := make([]flags.Completion, len(items))

	for i, v := range items {
		completions[i].Item = v
	}
	return completions
}

// OutputFile, when unmarshaled, embeds an os.File opened for writing.
// The file is opened with os.O_WRONLY|os.O_CREATE|os.O_TRUNC and FileMode 0600.
type OutputFile struct {
	*os.File
	// NopClose, if true, overrides the file's Close() method to do nothing
	// similar to ioutil.NopCloser.
	NopClose bool
}

// NewOutputFile creates an OutputFile from an existing opened os.File.
// Useful for setting defaults via struct initialization rather than struct tags.
func NewOutputFile(f *os.File) OutputFile {
	return OutputFile{File: f}
}

// UnmasrshalFlag returns the output file name
func (f *OutputFile) MarshalFlag() (string, error) {
	return f.Name(), nil
}

// UnmarshalFlag attempts to return a file opened for writing with the input
// string value used as the filename
func (f *OutputFile) UnmarshalFlag(value string) error {
	file, err := os.OpenFile(ExpandUser(value), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	*f = OutputFile{File: file}
	return nil
}

// Close, if OutputFile.NopClose is true, simply returns nil, otherwise the
// embedded os.File's Close is called.
func (f *OutputFile) Close() error {
	if f.NopClose {
		return nil
	}
	return f.File.Close()
}

// DefaultToStdout is a convenience method that should be used when setting defaults via
// struct initialization rather than tags so the default file is os.Stdout
func (f *OutputFile) DefaultToStdout() {
	*f = OutputFile{File: os.Stdout}
	f.NopClose = true
}

// Complete satisfies flags.Completer and returns potential file matches based on current input.
func (f *OutputFile) Complete(match string) []flags.Completion {
	items, _ := filepath.Glob(ExpandUser(match) + "*")
	completions := make([]flags.Completion, len(items))

	for i, v := range items {
		completions[i].Item = v
	}
	return completions
}
