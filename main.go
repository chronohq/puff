package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

const (
	defaultBlobBytes = 1024 * 1024
	defaultHexBytes  = 16
	defaultDelimiter = "\n"

	bytesParam     = "bytes"
	compactParam   = "compact"
	delimiterParam = "delimiter"
	numParam       = "num"
	outputParam    = "output"
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

// generateHex generates one or more random hexadecimal strings and outputs
// them. The length of each string is determined by the specified byte size.
// The output is printed to stdout, with each value separated by the specified
// delimiter, except the final value.
func generateHex(c *cli.Context) error {
	iterations := c.Int(numParam)

	if iterations <= 0 {
		return errors.New("num must be greater than 0")
	}

	delimiter := resolveDelimiter(c.String(delimiterParam))

	if len(delimiter) == 0 {
		return errors.New("delimiter cannot be blank")
	}

	for i := 0; i < iterations; i++ {
		src, err := randomBytes(c.Int(bytesParam))

		if err != nil {
			return err
		}

		if i < iterations-1 {
			fmt.Printf("%s%s", hex.EncodeToString(src), delimiter)
		} else {
			// In streaming mode, the final delimiter is a newline character.
			fmt.Println(hex.EncodeToString(src))
		}
	}

	return nil
}

// generateUUID generates one or more UUID version 7 strings in hexadeicmal.
func generateUUID(c *cli.Context) error {
	compact := c.Bool(compactParam)
	iterations := c.Int(numParam)

	if iterations <= 0 {
		return errors.New("num must be greater than 0")
	}

	delimiter := resolveDelimiter(c.String(delimiterParam))

	if len(delimiter) == 0 {
		return errors.New("delimiter cannot be blank")
	}

	for i := 0; i < iterations; i++ {
		id, err := uuid.NewV7()

		if err != nil {
			return err
		}

		line := id.String()

		if compact {
			line = strings.ReplaceAll(line, "-", "")
		}

		if i < iterations-1 {
			fmt.Printf("%s%s", line, delimiter)
		} else {
			fmt.Println(line)
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

	if len(outputFilePath) == 0 {
		return errors.New("output file path is required")
	}

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

func main() {
	app := &cli.App{
		Name:    "puff",
		Usage:   "Generate random values in different formats",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:   "hex",
				Usage:  "Generate random hexadecimal strings",
				Action: generateHex,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "bytes",
						Aliases: []string{"b"},
						Usage:   "length of the source data in bytes",
						Value:   defaultHexBytes,
					},
					&cli.IntFlag{
						Name:    "num",
						Aliases: []string{"n"},
						Usage:   "number of hex strings to generate",
						Value:   1,
					},
					&cli.StringFlag{
						Name:    "delimiter",
						Aliases: []string{"d", "delimit"},
						Usage:   "delimiter between values",
						Value:   defaultDelimiter,
					},
				},
			},
			{
				Name:   "uuid",
				Usage:  "Generate UUID strings (default: UUIDv7)",
				Action: generateUUID,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "num",
						Aliases: []string{"n"},
						Usage:   "number of uuid strings to generate",
						Value:   1,
					},
					&cli.BoolFlag{
						Name:  "compact",
						Usage: "print uuid strings without dashes",
					},
					&cli.StringFlag{
						Name:    "delimiter",
						Aliases: []string{"d", "delimit"},
						Usage:   "delimiter between values",
						Value:   defaultDelimiter,
					},
				},
			},
			{
				Name:    "binary",
				Aliases: []string{"bin"},
				Usage:   "Generate a random binary blob",
				Action:  generateBinaryBlob,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "bytes",
						Aliases: []string{"b"},
						Usage:   "number of random bytes to generate",
						Value:   defaultBlobBytes,
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "file path to write the generated bytes",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
