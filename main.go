/*
Copyright © 2023 RUSSELL JONES
*/
package main

import (
	"fmt"
	"os"

	"github.com/jonesrussell/page-prowler/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
