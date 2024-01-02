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
	if err := cmd.Execute(); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
