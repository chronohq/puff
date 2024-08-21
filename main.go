package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	defaultHexBytes  = 16
	defaultDelimiter = "\n"

	bytesParam     = "bytes"
	delimiterParam = "delimiter"
	numParam       = "num"
)

// randomBytes generates a slice of random bytes of the given length.
// Returns an error if the length is non-positive or if crypto/rand fails.
func randomBytes(len int) ([]byte, error) {
	if len == 0 {
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

func main() {
	app := &cli.App{
		Name:  "puff",
		Usage: "Generate random values in different formats",
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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
