package flagext

import (
	//"log"
	"os"
	//"path/filepath"
	//"github.com/jessevdk/go-flags"
)

const DefaultRWFileFlag int = os.O_CREATE | os.O_RDWR
const DefaultRWFilePerm os.FileMode = 0600

type file struct {
	*os.File
	nopClose bool
}

func (f *file) Close() error {
	if f.nopClose {
		return nil
	}
	return f.File.Close()
}

func (f *file) MarshalFlag() (string, error) {
	return f.Name(), nil
}

type InputFile struct {
	file
	nopClose bool
}

func (i *InputFile) UnmarshalFlag(value string) error {
	inFile, err := os.Open(ExpandUser(value))
	if err != nil {
		return err
	}
	i.file = file{File: inFile, nopClose: i.NopClose()}
	return nil
}

func (i *InputFile) SetDefaultFile(defaultFile *os.File, nopClose bool) {
	i.file = file{File: defaultFile, nopClose: nopClose}
}

func (i *InputFile) SetNopClose(value bool) {
	i.nopClose = value
}

func (i *InputFile) NopClose() bool {
	return i.nopClose
}

type OutputFile struct {
	file
	nopClose bool
}

func (o *OutputFile) UnmarshalFlag(value string) error {
	outFile, err := os.OpenFile(ExpandUser(value), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	o.file = file{File: outFile, nopClose: o.NopClose()}
	return nil
}

func (o *OutputFile) SetDefaultFile(defaultFile *os.File, nopClose bool) {
	o.file = file{File: defaultFile, nopClose: nopClose}
}

func (o *OutputFile) SetNopClose(value bool) {
	o.nopClose = value
}

func (o *OutputFile) NopClose() bool {
	return o.nopClose
}

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

func (rw *ReadWriteFile) UnmarshalFlag(value string) error {
	rwFile, err := os.OpenFile(value, rw.Flag(), rw.Perm())
	if err != nil {
		return err
	}
	rw.file = file{File: rwFile, nopClose: rw.NopClose()}
	rw.parsed = true
	return nil
}

func (rw *ReadWriteFile) SetDefaultFile(defaultFile *os.File, nopClose bool) {
	rw.file = file{File: defaultFile, nopClose: nopClose}
}

func (rw *ReadWriteFile) SetNopClose(value bool) {
	rw.nopClose = value
}

func (rw *ReadWriteFile) NopClose() bool {
	return rw.nopClose
}

func (rw *ReadWriteFile) SetFlag(value int) {
	rw.flag.flag = value
	rw.flag.overridden = true
}

func (rw *ReadWriteFile) Flag() int {
	if (!rw.flag.overridden) && rw.flag.flag == 0 {
		rw.flag.flag = DefaultRWFileFlag
	}
	return rw.flag.flag
}

func (rw *ReadWriteFile) SetPerm(value os.FileMode) {
	rw.perm.perm = value
	rw.perm.overridden = true
}

func (rw *ReadWriteFile) Perm() os.FileMode {
	if (!rw.perm.overridden) && rw.perm.perm == 0 {
		rw.perm.perm = DefaultRWFilePerm
	}
	return rw.perm.perm
}
