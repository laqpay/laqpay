/*
cli is a command line client for interacting with a laqpay node and offline wallet management
*/
package main

import (
	"fmt"
	"os"

	"../../src/cli"
)

func main() {

	cfg, err := cli.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	laqCLI, err := cli.NewCLI(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := laqCLI.Execute(); err != nil {
		os.Exit(1)
	}
}
