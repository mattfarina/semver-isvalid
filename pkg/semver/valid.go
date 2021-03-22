// Package semver provides a means to validate semantic versions
package semver

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	// ErrEmptyString is returned when an empty string is passed in for parsing.
	ErrEmptyString = errors.New("Version string empty")

	// ErrInvalidNumberParts is returned when a number of parts, other than 3,
	// is found
	ErrInvalidNumberParts = errors.New("Version does not have 3 parts")

	// ErrInvalidCharacters is returned when invalid characters are found as
	// part of a version
	ErrInvalidCharacters = errors.New("Invalid characters in version")

	// ErrSegmentStartsZero is returned when a version segment starts with 0.
	// This is invalid in SemVer.
	ErrSegmentStartsZero = errors.New("Version segment starts with 0")
)

// Validate accepts one argument (the version as a string) and returns 2 values
// which are:
// - An error message if the version is not a semantic version
// - A slice of messages with details about the version
func Validate(ver string) (error, []string) {

	// Check if an empty string was passed in
	if len(ver) == 0 {
		return ErrEmptyString, []string{}
	}

	// Split the parts into [0]major, [1]minor, and [2]patch,prerelease,build
	// Semantic Versions are required to have 3 parts
	parts := strings.SplitN(ver, ".", 3)
	if len(parts) != 3 {
		num := len(parts)
		return ErrInvalidNumberParts, []string{fmt.Sprintf("Found %d number of parts", num)}
	}

	v := &version{}

	var tmp []string
	// Trim the patch release right to left to find any metadata or prerelease
	if strings.ContainsAny(parts[2], "-+") {
		tmp = strings.SplitN(parts[2], "+", 2)
		if len(tmp) > 1 {
			v.metadata = tmp[1]
			parts[2] = tmp[0]
		}

		tmp = strings.SplitN(parts[2], "-", 2)
		if len(tmp) > 1 {
			v.pre = tmp[1]
			parts[2] = tmp[0]
		}
	}

	var messages []string

	// Validate each of the major, minor, patch release segments
	for i, p := range parts {
		if !containsOnly(p, num) {
			messages = append(messages, fmt.Sprintf("Illegal non-numeric characters found in %q part", numToName(i)))
			return ErrInvalidCharacters, messages
		}

		if len(p) > 1 && p[0] == '0' {
			messages = append(messages, fmt.Sprintf("Illegal leading 0 found in %q part", numToName(i)))
			return ErrSegmentStartsZero, messages
		}
	}

	// Parse to check the major, minor, and patch versions
	var err error
	v.major, err = strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		messages = append(messages, fmt.Sprint("Unable to parse major part. Must be valid numeric characters [0-9]"))
		return err, messages
	}
	messages = append(messages, fmt.Sprintf("Found major version of %d", v.major))

	v.minor, err = strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		messages = append(messages, fmt.Sprint("Unable to parse minor part. Must be valid numeric characters [0-9]"))
		return err, messages
	}
	messages = append(messages, fmt.Sprintf("Found minor version of %d", v.minor))

	v.patch, err = strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		messages = append(messages, fmt.Sprint("Unable to parse patch part. Must be valid numeric characters [0-9]"))
		return err, messages
	}
	messages = append(messages, fmt.Sprintf("Found patch version of %d", v.patch))

	if v.pre != "" {
		tmp = strings.Split(v.pre, ".")
		for _, p := range tmp {
			if containsOnly(p, num) {
				if len(p) > 1 && p[0] == '0' {
					messages = append(messages, fmt.Sprintf("Illegal leading 0 found in pre-release numeric part %q", p))
					return ErrSegmentStartsZero, messages
				}
			} else if !containsOnly(p, allowed) {
				messages = append(messages, fmt.Sprintf("Illegal characters found in pre-release non-numeric part %q. Must be [0-9A-Za-z-]", p))
				return ErrInvalidCharacters, messages
			}
		}
		messages = append(messages, fmt.Sprintf("Version is a pre-release version rather than a stable release version with a pre-release identifier of %q", v.pre))
		messages = append(messages, fmt.Sprint("NOTICE: A pre-release version indicates that the version is unstable and might not satisfy the intended compatibility requirements as denoted by its associated normal version."))
	}

	if v.metadata != "" {
		tmp = strings.Split(v.metadata, ".")
		for _, p := range tmp {
			if !containsOnly(p, allowed) {
				messages = append(messages, fmt.Sprintf("Illegal characters found in metadata part %q. Must be [0-9A-Za-z-]", p))
				return ErrInvalidCharacters, messages
			}
		}
		messages = append(messages, fmt.Sprintf("Found build metadate on version of %q", v.metadata))
		messages = append(messages, fmt.Sprint("NOTICE: Build metadata MUST be ignored when determining version precedence. Thus two versions that differ only in the build metadata, have the same precedence."))
	}

	return nil, messages
}

const num string = "0123456789"
const allowed string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-" + num

type version struct {
	major, minor, patch uint64
	pre                 string
	metadata            string
}

func numToName(i int) string {
	switch i {
	case 0:
		return "major"
	case 1:
		return "minor"
	case 2:
		return "patch"
	}

	panic("Invalid part number")
}

// Like strings.ContainsAny but does an only instead of any.
func containsOnly(s string, comp string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return !strings.ContainsRune(comp, r)
	}) == -1
}
