package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/mattfarina/semver-isvalid/pkg/semver"
	"github.com/spf13/cobra"
)

func main() {
	var cmd = &cobra.Command{
		Use:   "semver-isvalid [version]",
		Short: "semver-isvalid allows you to validate a single semantic version",
		Long:  longdesc,
		Run: func(cmd *cobra.Command, args []string) {
			la := len(args)
			if la == 0 {
				_ = cmd.Help()
				return
			} else if la != 1 {
				red.Fprintf(os.Stderr, "Wrong number of arguments supplied. 1 argument required but found %d\n", la)
				os.Exit(1)
			}
			validate(args[0])
		},
	}

	cmd.PersistentFlags().BoolVar(&withV, "with-v", false, "allow v at start of version")

	cmd.Execute()
}

var red = color.New(color.FgRed)

var withV = false

const longdesc = `semver-isvalid allows you to validate a single semantic version

In addition to validating a semantic version, semver-isvalid will tell you
information about the version it has discovered. This can include specific
details about validation issues, what it found about the version, and notices
that may help in understanding the version.

For example, you can run it like so:

    $ semver-isvalid 1.2.3
    Found major version of 1
    Found minor version of 2
    Found patch version of 3
    Semantic Version is valid

An example with an invalid version would be:

    $ semver-isvalid 1.2.03
    Illegal leading 0 found in "patch" part
    Invalid Semantic Version: Version segment starts with 0. For more
    information see https://semver.org

The "v" at the start of a version is NOT part of Semantic Versioning. If you
want to allow that as part of the version you can use the --with-v flag. For
example:

    $ semver-isvalid v1.2.3 --with-v
    Found major version of 1
    Found minor version of 2
    Found patch version of 3
    Semantic Version is valid

Without the --with-v this would have returned an error as being invalid.

For those who look at exit codes, each type of error has a unique exit code.
The codes include:

- 1: Invalid number of arguments passed to application
- 2: A general invalid semantic version
- 3: The version passed in evaluates to an empty string
- 4: There are an invalid number of version parts. 3 are required for Semantic
     Versions
- 5: Invalid characters were found in a part of a Semantic Version
- 6: A numeric segment starts with 0

For more information on Semantic Versions please visit the specification
at https://semver.org.

You can learn more about this application at
https://github.com/mattfarina/semver-isvalid

`

func validate(ver string) {

	if withV {
		ver = strings.TrimPrefix(ver, "v")
	}

	err, msgs := semver.Validate(ver)
	for _, v := range msgs {
		fmt.Println(v)
	}

	errmsg := "Invalid Semantic Version: %s. For more information see https://semver.org\n"
	switch err {
	case semver.ErrEmptyString:
		red.Fprintf(os.Stderr, errmsg, semver.ErrEmptyString)
		os.Exit(3)
	case semver.ErrInvalidNumberParts:
		red.Fprintf(os.Stderr, errmsg, semver.ErrInvalidNumberParts)
		os.Exit(4)

	case semver.ErrInvalidCharacters:
		red.Fprintf(os.Stderr, errmsg, semver.ErrInvalidCharacters)
		os.Exit(5)

	case semver.ErrSegmentStartsZero:
		red.Fprintf(os.Stderr, errmsg, semver.ErrSegmentStartsZero)
		os.Exit(6)
	case nil:
		fmt.Println("Semantic Version is valid")
	default:
		red.Fprint(os.Stderr, "Invalid Semantic Version. For more information see https://semver.org\n")
		os.Exit(2)
	}
}
