package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	err := FlagSet().Run(os.Args)

	if err != nil {
		return err
	}

	return nil
}
