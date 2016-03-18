package flagtypes_test

import (
	"fmt"
	"github.com/bmoyles/go-flags-types"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"os"
)

func ExampleFileOrStdin_stdin() {
	var err error

	type opts struct {
		// Example of FileOrStdin with a default value of "-" which should return os.Stdin if the flag is omitted
		Input flagtypes.FileOrStdin `short:"i" long:"input" description:"Input file" default:"-"`
	}

	// Instance of opts that will hold the results from parsing an empty args slice
	var emptyOpts opts

	// Fake args that omits -s which should result in emptyArgs.Input == os.Stdin
	emptyArgs := []string{}

	_, err = flags.ParseArgs(&emptyOpts, emptyArgs)
	if err != nil {
		panic(err)
	}

	inputStat, err := emptyOpts.Input.Stat()
	if err != nil {
		panic(err)
	}

	stdinStat, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Input file name: %s\n", emptyOpts.Input.Name())
	fmt.Printf("Input file same as os.Stdin? %t\n", os.SameFile(inputStat, stdinStat))
	// Output: Input file name: /dev/stdin
	// Input file same as os.Stdin? true
}

func ExampleFileOrStdin_file() {
	var err error

	type opts struct {
		// Example of FileOrStdin with a default value of "-" which should return os.Stdin if the flag is omitted
		Input flagtypes.FileOrStdin `short:"i" long:"input" description:"Input file" default:"-"`
	}

	// Instance of opts that will hold the results of parsing args specifying a file
	var fileOpts opts

	// a dummy file for the sake of example
	dummyFile, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	defer os.Remove(dummyFile.Name())

	// Fake args that provides the dummyFile name
	fileArgs := []string{"-i", dummyFile.Name()}

	_, err = flags.ParseArgs(&fileOpts, fileArgs)
	if err != nil {
		panic(err)
	}

	inputStat, err := fileOpts.Input.Stat()
	if err != nil {
		panic(err)
	}
	dummyStat, err := dummyFile.Stat()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Input file same as dummyFile?: %t\n", os.SameFile(inputStat, dummyStat))
	// Output: Input file same as dummyFile?: true
}

func ExampleFileOrStdout_stdout() {
	var err error

	type opts struct {
		// Example of FileOrStdout with a default value of "-" which should return os.Stdout if the flag is omitted
		Output flagtypes.FileOrStdout `short:"o" long:"output" description:"Output file" default:"-"`
	}

	// Instance of opts that will hold the results from parsing an empty args slice
	var emptyOpts opts

	// Fake args that omits -s which should result in emptyArgs.Input == os.Stdout
	emptyArgs := []string{}

	_, err = flags.ParseArgs(&emptyOpts, emptyArgs)
	if err != nil {
		panic(err)
	}

	outputStat, err := emptyOpts.Output.Stat()
	if err != nil {
		panic(err)
	}

	stdoutStat, err := os.Stdout.Stat()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Output file same as os.Stdout? %t\n", os.SameFile(outputStat, stdoutStat))
	// Output: Output file same as os.Stdout? true
}

func ExampleFileOrStdout_file() {
	var err error

	type opts struct {
		// Example of FileOrStdout with a default value of "-" which should return os.Stdout if the flag is omitted
		Output flagtypes.FileOrStdout `short:"o" long:"output" description:"Output file" default:"-"`
	}

	// Instance of opts that will hold the results of parsing args specifying a file
	var fileOpts opts

	// a dummy file for the sake of example
	dummyFile, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	defer os.Remove(dummyFile.Name())

	// Fake args that provides the dummyFile name
	fileArgs := []string{"-o", dummyFile.Name()}

	_, err = flags.ParseArgs(&fileOpts, fileArgs)
	if err != nil {
		panic(err)
	}

	outputStat, err := fileOpts.Output.Stat()
	if err != nil {
		panic(err)
	}
	dummyStat, err := dummyFile.Stat()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Output file same as dummyFile?: %t\n", os.SameFile(outputStat, dummyStat))
	fmt.Printf("Output file fileMode: %04o", outputStat.Mode())
	// Output: Output file same as dummyFile?: true
	// Output file fileMode: 0600
}
