package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

type config struct {
	length   int
	specials bool
	count    int
}

var errPosArgSpecified = errors.New("Positional arguments specified")

func main() {
	c, err := parseArgs(os.Stdout, os.Args[1:])
	if err != nil {
		if errors.Is(err, errPosArgSpecified) {
			fmt.Fprintln(os.Stdout, err)
		}
		os.Exit(1)
	}

	err = ValidateArgs(c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}

	GeneratePass(os.Stdout, c)
}

func parseArgs(w io.Writer, args []string) (config, error) {
	c := config{}

	fs := flag.NewFlagSet("PassGenerator", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.IntVar(&c.length, "l", 8, "length of password")
	fs.BoolVar(&c.specials, "s", false, "add special symbols")
	fs.IntVar(&c.count, "c", 1, "number of passwords to generate")

	err := fs.Parse(args)
	if err != nil {
		return c, err
	}

	if fs.NArg() > 0 {
		return c, errors.New("unexpected positional arguments")
	}

	return c, nil
}

func ValidateArgs(c config) error {
	if c.length < 1 {
		return errors.New("length must be greater than zero")
	}

	if c.count < 1 {
		return errors.New("count must be greater than zero")
	}

	return nil
}

func GeneratePass(w io.Writer, c config) {
	var alph = "abcdefghigklmnopqrstuvwxyz"
	var seedRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	switch c.specials {
	case true:
		var fullAlph = alph + "1234567890!@£$%^&*()-=_+:;|/?><~`][{}±§"

		for i := 0; i < c.count; i++ {
			pass := make([]byte, c.length)
			for j := range pass {
				pass[j] = fullAlph[seedRand.Intn(len(fullAlph))]
			}
			fmt.Fprintln(w, string(pass))
		}

	case false:
		for i := 0; i < c.count; i++ {
			pass := make([]byte, c.length)
			for j := range pass {
				pass[j] = alph[seedRand.Intn(len(alph))]
			}
			fmt.Fprintln(w, string(pass))
		}
	}
}
