# Semantic Version Validator

This tiny application and go library enable you to test if a semantic version is valid or not.

## Usage

There are two ways to use this codebase. 1) as an application you can install and run and 2) as a Go library you can use in your own application.

### Console Application

The console application provides a simple tool that you can use locally or as part of a CI system to check that the version you're using is valid. It provides details about the version or the error to point you in the right direction if the version is not valid.

The use is pretty simple. You pass in a single version as an argument and it will tell you about it.

For example, you can run it like so:

```sh
$ semver-isvalid 1.2.3
Found major version of 1
Found minor version of 2
Found patch version of 3
Semantic Version is valid
```

An example with an invalid version would be:

```sh
$ semver-isvalid 1.2.03
Illegal leading 0 found in "patch" part
Invalid Semantic Version: Version segment starts with 0. For more information see https://semver.org
```

For those who look at exit codes, each type of error has a unique exit code.
The codes include:

- 1: Invalid number of arguments passed to application
- 2: A general invalid semantic version
- 3: The version passed in evaluates to an empty string
- 4: There are an invalid number of version parts. 3 are required for Semantic Versions
- 5: Invalid characters were found in a part of a Semantic Version
- 6: A numeric segment starts with 0

### Go Library

The business logic the console application uses is provided in a Go package that other applications can import and use.

Here is a simple example:

```go
package main

import (
    "fmt"

    "github.com/mattfarina/semver-isvalid/pkg/semver"
)

func main() {
    err, msgs := semver.Validate("1.2.3")
    for _, v := range msgs {
        fmt.Println(v)
    }

    if err != nil {
        fmt.Printf("Error in version: %s\n", err)
    } else {
        fmt.Println("Version is valid")
    }
}
```

## Inspiration

It is not uncommon for people or tooling to inadvertently create semantic versions that are invalid. This can lead to consequences when working with tools that depend on valid semantic versions.

[SUSE](https://www.suse.com/) Hackweek opened up an opportunity to hack on a tool to check for validity.

Having written a [semver parsing library](https://github.com/masterminds/semver) (and spending way to much time reading the spec) this seemed like an easy and possibly useful tool to write.

## License

This licensed under the MIT license just as [the inspiration codebase](https://github.com/masterminds/semver) is.
