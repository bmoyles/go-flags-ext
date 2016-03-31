package flagext

import (
	"os"
)

const (
	// DefaultRWFileFlag is the default set of flags used for ReadWriteFile
	DefaultRWFileFlag int = os.O_CREATE | os.O_RDWR

	// DefaultRWFilePerm is the default FileMode used for ReadWriteFile
	DefaultRWFilePerm os.FileMode = 0600
)

// file is an internal struct that actually contains the file object after parsing.
// Methods that are common to other types that wrap file and only operate on file are included here.
type file struct {
	*os.File
	nopClose bool
}

// Close, when SetNopClose was called with true, or, SetDefaultFile was called with nopClose=true,
// makes the Close operation on the file a no-op similar to io/ioutil.NopCloser. If SetNopClose was not
// called or is false otherwise, the underlying file's close method is called.
func (f *file) Close() error {
	if f.nopClose {
		return nil
	}
	return f.File.Close()
}

// MarshalFlag satisfies the go-flags Marshaler interface and returns the file object's file name
// to serve as a string representation of the file.
func (f *file) MarshalFlag() (string, error) {
	return f.Name(), nil
}

// InternalFile allows one to extract the actual *os.File should it be needed directly (eg to pass to
// functions that ask for an *os.File directly rather than an io interface). Using this directly
// is not generally recommended as it breaks features like NopClose which are handled in the wrapper
// and cannot be applied directly to *os.File.
func (f *file) InternalFile() *os.File {
	return f.File
}

// An InputFile, after parsing, contains an embedded *os.File that is either the user-supplied
// file from command line args, or, the default (if any) if the user did not pass in an explicit
// filename. The file is opened for reading.
type InputFile struct {
	file
	nopClose bool
}

// UnmarshalFlag satisfies the go-flags Unmarshaler interface and is called at parse
// time to convert the user-provided file name into a file opened for reading.
func (f *InputFile) UnmarshalFlag(value string) error {
	inFile, err := os.Open(ExpandUser(value))
	if err != nil {
		return err
	}
	f.file = file{File: inFile, nopClose: f.NopClose()}
	return nil
}

// SetDefaultFile should be called on the empty struct passed to go-flags' parser
// to set a default file to be used if a user did not present a file name on the command line.
// If nopClose is true, the file's Close() method will be wrapped so it returns without actually
// closing the file similar to io/ioutil.NopCloser (handy for files like os.Stdin or os.Stdout that
// one may not want to close immediately).
func (f *InputFile) SetDefaultFile(defaultFile *os.File, nopClose bool) {
	f.file = file{File: defaultFile, nopClose: nopClose}
}

// DefaultToStdin is a convenience method that sets the default file to os.Stdin and nopClose to true
// allowing default input to be taken from stdin rather than a file.
func (f *InputFile) DefaultToStdin() {
	f.SetDefaultFile(os.Stdin, true)
}

// SetNopClose should be called on the empty struct passed to go-flags' parser and
// causes user-provided files from the command line to have their Close() method overridden
// and avoid closing the file if Close() is called. This is independent of the default file
// which is only used if the user did not supply a file on the command line.
func (f *InputFile) SetNopClose(value bool) {
	f.nopClose = value
}

// NopClose returns the current value of nopClose.
func (f *InputFile) NopClose() bool {
	return f.nopClose
}

// An OutputFile, after parsing, contains an embedded *os.File that is either the user-supplied
// file from command line args, or, the default (if any) if the user did not pass in an explicit
// filename. The file is opened for writing using flags os.O_WRONLY|os.O_CREATE|os.O_TRUNC and
// FileMode 0600.
type OutputFile struct {
	file
	nopClose bool
}

// UnmarshalFlag satisfies the go-flags Unmarshaler interface and is called at parse
// time to convert the user-provided file name into a file opened for writing.
func (f *OutputFile) UnmarshalFlag(value string) error {
	outFile, err := os.OpenFile(ExpandUser(value), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	f.file = file{File: outFile, nopClose: f.NopClose()}
	return nil
}

// SetDefaultFile should be called on the empty struct passed to go-flags' parser
// to set a default file to be used if a user did not present a file name on the command line.
// If nopClose is true, the file's Close() method will be wrapped so it returns without actually
// closing the file similar to io/ioutil.NopCloser (handy for files like os.Stdin or os.Stdout that
// one may not want to close immediately).
func (f *OutputFile) SetDefaultFile(defaultFile *os.File, nopClose bool) {
	f.file = file{File: defaultFile, nopClose: nopClose}
}

// DefaultToStdout is a convenience method that sets the default file to os.Stdout and nopClose to true
// allowing default input to be taken from stdin rather than a file.
func (f *OutputFile) DefaultToStdout() {
	f.SetDefaultFile(os.Stdout, true)
}

// SetNopClose should be called on the empty struct passed to go-flags' parser and
// causes user-provided files from the command line to have their Close() method overridden
// and avoid closing the file if Close() is called. This is independent of the default file
// which is only used if the user did not supply a file on the command line.
func (f *OutputFile) SetNopClose(value bool) {
	f.nopClose = value
}

// NopClose returns the current value of nopClose.
func (f *OutputFile) NopClose() bool {
	return f.nopClose
}

// An ReadWriteFile, after parsing, contains an embedded *os.File that is either the user-supplied
// file from command line args, or, the default (if any) if the user did not pass in an explicit
// filename. The file is opened for both reading and writing with the default flag contained in
// DefaultRWFileFlag and default FileMode contained in DefaultRWFileMode.
type ReadWriteFile struct {
	file
	flag struct {
		flag       int
		overridden bool
	}
	perm struct {
		perm       os.FileMode
		overridden bool
	}
	parsed bool
}

// UnmarshalFlag satisfies the go-flags Unmarshaler interface and is called at parse
// time to convert the user-provided file name into a file opened for reading and writing.
func (f *ReadWriteFile) UnmarshalFlag(value string) error {
	rwFile, err := os.OpenFile(value, f.Flag(), f.Perm())
	if err != nil {
		return err
	}
	f.file = file{File: rwFile, nopClose: f.NopClose()}
	f.parsed = true
	return nil
}

// SetDefaultFile should be called on the empty struct passed to go-flags' parser
// to set a default file to be used if a user did not present a file name on the command line.
// If nopClose is true, the file's Close() method will be wrapped so it returns without actually
// closing the file similar to io/ioutil.NopCloser (handy for files like os.Stdin or os.Stdout that
// one may not want to close immediately).
func (f *ReadWriteFile) SetDefaultFile(defaultFile *os.File, nopClose bool) {
	f.file = file{File: defaultFile, nopClose: nopClose}
}

// SetNopClose should be called on the empty struct passed to go-flags' parser and
// causes user-provided files from the command line to have their Close() method overridden
// and avoid closing the file if Close() is called. This is independent of the default file
// which is only used if the user did not supply a file on the command line.
func (f *ReadWriteFile) SetNopClose(value bool) {
	f.nopClose = value
}

// NopClose returns the current value of nopClose.
func (f *ReadWriteFile) NopClose() bool {
	return f.nopClose
}

// SetFlag should be called on the empty struct passed to go-flags' parser and causes user-provided
// files on the command line to be opened with the provided flag value rather than DefaultRWFileFlag.
func (f *ReadWriteFile) SetFlag(value int) {
	f.flag.flag = value
	f.flag.overridden = true
}

// Flag returns the current default flag value.
func (f *ReadWriteFile) Flag() int {
	if (!f.flag.overridden) && f.flag.flag == 0 {
		f.flag.flag = DefaultRWFileFlag
	}
	return f.flag.flag
}

// SetPerm should be called on the empty struct passed to go-flags' parser and causes user-provided
// files on the command line to be opened with the provided FileMode rather than DefaultRWFilePerm.
func (f *ReadWriteFile) SetPerm(value os.FileMode) {
	f.perm.perm = value
	f.perm.overridden = true
}

// Perm returns the current default perm value.
func (f *ReadWriteFile) Perm() os.FileMode {
	if (!f.perm.overridden) && f.perm.perm == 0 {
		f.perm.perm = DefaultRWFilePerm
	}
	return f.perm.perm
}
