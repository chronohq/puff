// Copyright Chrono Technologies LLC
// SPDX-License-Identifier: MIT

package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

const (
	uuidV4 = 4
	uuidV7 = 7

	defaultBlobBytes   = 1024 * 1024
	defaultDataBytes   = 16
	defaultUUIDVersion = uuidV4
	defaultDelimiter   = "\n"

	bytesParam       = "bytes"
	compactParam     = "compact"
	delimiterParam   = "delimiter"
	numParam         = "num"
	outputParam      = "output"
	suffixParam      = "suffix"
	timeParam        = "time"
	urlSafeParam     = "url-safe"
	uuidVersionParam = "version"

	red      = "\033[31m"
	resetRed = "\033[0m"
)

var (
	errInvalidIterations = errors.New("--num must be greater than 0")
	errBlankDelimiter    = errors.New("--delimiter cannot be blank")
	errInvalidUUIDV7Time = errors.New("UUIDv7 does not support pre unix epoch")
)

// version holds the application version number. This value is set at build
// time using the -ldflags build flag. The default value here is a placehodler.
var version = "placeholder"

// randomBytes generates a slice of random bytes of the given length.
// Returns an error if the length is non-positive or if crypto/rand fails.
func randomBytes(len int) ([]byte, error) {
	if len <= 0 {
		return nil, errors.New("length must be greater than 0")
	}

	ret := make([]byte, len)

	_, err := rand.Read(ret)

	return ret, err
}

// resolveDelimiter converts specific escape sequences into control characters.
//
// Supported escape sequences:
//
// \n - Newline
// \t - Tab
//
// Any other input is returned as-is, for example:
//
// resolveDelimiter(`\n`) -> "\n" (newline character)
// resolveDelimiter(`\t`) -> "\t" (tab character)
// resolveDelimiter(`,`) -> "," (literal comma)
func resolveDelimiter(delim string) string {
	switch delim {
	case `\n`:
		return "\n"
	case `\t`:
		return "\t"
	default:
		return delim
	}
}

// openFileExclusive opens a file for writing, failing if the file exists.
func openFileExclusive(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
}

// paintError wraps the error message in ANSI red color codes and returns a
// new error. It returns nil if the provided error is nil.
func paintError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(red+"%v"+resetRed, err)
}

// printError writes an error message to the provided io.Writer in red text.
// It formats the error message with a preceding "error:" label and follows
// the output with two newline characters for clear separation.
func printError(w io.Writer, err error) {
	redErr := paintError(fmt.Errorf("error: %v", err))
	fmt.Fprintf(w, "%v\n\n", redErr)
}

// generateHex generates one or more random hexadecimal strings and outputs
// them. The length of each string is determined by the specified byte size.
// The output is printed to stdout, with each value separated by the specified
// delimiter, except the final value.
func generateHex(c *cli.Context) error {
	iterations := c.Int(numParam)
	suffix := c.String(suffixParam)

	if iterations <= 0 {
		printError(c.App.ErrWriter, errInvalidIterations)
		return cli.ShowSubcommandHelp(c)
	}

	delimiter := resolveDelimiter(c.String(delimiterParam))

	if len(delimiter) == 0 {
		return errBlankDelimiter
	}

	for i := range iterations {
		src, err := randomBytes(c.Int(bytesParam))

		if err != nil {
			return paintError(err)
		}

		if i == iterations-1 {
			delimiter = defaultDelimiter
		}

		fmt.Printf("%s%s%s", hex.EncodeToString(src), suffix, delimiter)
	}

	return nil
}

// generateUUID generates one or more UUID strings in hexadeicmal. It defaults
// to generating version 7 UUIDs but also supports the widely used version 4.
func generateUUID(c *cli.Context) error {
	version := c.Int(uuidVersionParam)
	compact := c.Bool(compactParam)
	iterations := c.Int(numParam)
	suffix := c.String(suffixParam)

	if version != uuidV4 && version != uuidV7 {
		return paintError(errors.New("invalid uuid version (supported: 4, 7)"))
	}

	if iterations <= 0 {
		printError(c.App.ErrWriter, errInvalidIterations)
		return cli.ShowSubcommandHelp(c)
	}

	delimiter := resolveDelimiter(c.String(delimiterParam))

	if len(delimiter) == 0 {
		return paintError(errBlankDelimiter)
	}

	for i := range iterations {
		var err error
		var id uuid.UUID

		switch version {
		case uuidV4:
			id, err = generateUUIDV4()
		case uuidV7:
			id, err = generateUUIDV7(c.String(timeParam))
		}

		if err != nil {
			return paintError(err)
		}

		line := id.String()

		if compact {
			line = strings.ReplaceAll(line, "-", "")
		}

		if i == iterations-1 {
			delimiter = defaultDelimiter
		}

		fmt.Printf("%s%s%s", line, suffix, delimiter)
	}

	return nil
}

// generateBase64 generates cryptographically secure, randomly generated
// base64-encoded strings. It uses standard or URL-safe encoding based on
// the "--url-safe" command-line option.
func generateBase64(c *cli.Context) error {
	iterations := c.Int(numParam)
	delimiter := resolveDelimiter(c.String(delimiterParam))
	urlSafe := c.Bool(urlSafeParam)
	suffix := c.String(suffixParam)

	if len(delimiter) == 0 {
		return paintError(errBlankDelimiter)
	}

	for i := range iterations {
		src, err := randomBytes(c.Int(bytesParam))

		if err != nil {
			return paintError(err)
		}

		if i == iterations-1 {
			delimiter = defaultDelimiter
		}

		if urlSafe {
			fmt.Printf("%s%s%s", base64.RawURLEncoding.EncodeToString(src), suffix, delimiter)
		} else {
			fmt.Printf("%s%s%s", base64.StdEncoding.EncodeToString(src), suffix, delimiter)
		}
	}

	return nil
}

// generateBinaryBlob generates a binary blob of random bytes and writes it
// to the specified file. The file path must be provided. The number of bytes
// to generate is 1MB unless specified by the command line argument.
//
// Currently, a cryptographically secure RNG is used to produce the data.
// While this is suitable for use-cases like unique ID generation, it may be
// more than necessary for generating test blobs.
func generateBinaryBlob(c *cli.Context) error {
	outputFilePath := c.String(outputParam)

	file, err := openFileExclusive(outputFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	src, err := randomBytes(c.Int(bytesParam))
	if err != nil {
		return err
	}

	if _, err = file.Write(src); err != nil {
		return err
	}

	return nil
}

// parseTimeInput attempts to parse the given string time input. On success,
// it will return the time.Time value computed from the input.
func parseTimeInput(input string) (time.Time, error) {
	// first, check if the input is a unix timestamp
	if parsed, err := strconv.ParseInt(input, 10, 64); err == nil {
		// UUIDv7 does not support pre unix epoch
		if parsed < 0 {
			return time.Time{}, errInvalidUUIDV7Time
		}

		if parsed <= 9999999999 {
			return time.Unix(parsed, 0), nil
		} else {
			return time.UnixMilli(parsed), nil
		}
	}

	// not an integer, try various string formats
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, input); err == nil {
			if t.Before(time.Unix(0, 0)) {
				return time.Time{}, errInvalidUUIDV7Time
			}
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid time value: %s", input)
}

func main() {
	delimiterFlag := &cli.StringFlag{
		Name:    "delimiter",
		Aliases: []string{"d"},
		Usage:   "delimiter between values",
		Value:   defaultDelimiter,
	}

	suffixFlag := &cli.StringFlag{
		Name:    "suffix",
		Aliases: []string{"s"},
		Usage:   "optional value to append to the values",
	}

	hexCommandFlags := []cli.Flag{
		&cli.IntFlag{
			Name:    "bytes",
			Aliases: []string{"b"},
			Usage:   "length of the source data in bytes",
			Value:   defaultDataBytes,
		},
		&cli.IntFlag{
			Name:    "num",
			Aliases: []string{"n"},
			Usage:   "number of hex strings to generate",
			Value:   1,
		},
		delimiterFlag,
		suffixFlag,
	}

	app := &cli.App{
		Name:    "puff",
		Usage:   "Generate random values in different formats",
		Version: version,

		// default to the hex command if no subcommand is provided
		Action: generateHex,
		Flags:  hexCommandFlags,

		Commands: []*cli.Command{
			{
				Name:   "hex",
				Usage:  "Generate random hexadecimal strings",
				Action: generateHex,
				Flags:  hexCommandFlags,
			},
			{
				Name:   "uuid",
				Usage:  "Generate UUID strings (default: UUIDv4)",
				Action: generateUUID,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "num",
						Aliases: []string{"n"},
						Usage:   "number of uuid strings to generate",
						Value:   1,
					},
					&cli.IntFlag{
						Name:    "version",
						Aliases: []string{"v"},
						Usage:   "uuid version to generate",
						Value:   defaultUUIDVersion,
					},
					&cli.BoolFlag{
						Name:  "compact",
						Usage: "print uuid strings without dashes",
					},
					&cli.StringFlag{
						Name:    "time",
						Aliases: []string{"t"},
						Usage:   "timestamp for UUID v7 (iso8601 or unix timestamp)",
					},
					delimiterFlag,
					suffixFlag,
				},
			},
			{
				Name:   "base64",
				Usage:  "Generate base64 strings (default: 16-bytes)",
				Action: generateBase64,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "bytes",
						Aliases: []string{"b"},
						Usage:   "length of the source data in bytes",
						Value:   defaultDataBytes,
					},
					&cli.IntFlag{
						Name:    "num",
						Aliases: []string{"n"},
						Usage:   "number of base64 strings to generate",
						Value:   1,
					},
					&cli.BoolFlag{
						Name:  "url-safe",
						Usage: "use url-safe encoding",
					},
					delimiterFlag,
					suffixFlag,
				},
			},
			{
				Name:   "binary",
				Usage:  "Generate a random binary blob",
				Action: generateBinaryBlob,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "bytes",
						Aliases: []string{"b"},
						Usage:   "number of random bytes to generate",
						Value:   defaultBlobBytes,
					},
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"o"},
						Usage:    "file path to write the generated bytes",
						Required: true,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
