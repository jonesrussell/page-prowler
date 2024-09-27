package cmd

import (
	"bytes"

	"github.com/spf13/cobra"
)

import "fmt"

// ExecuteCommand is a helper function intended for use in tests.
// It executes a Cobra command with the provided arguments and returns the captured output and error.
func ExecuteCommand(root *cobra.Command, args ...string) (output string, err error) {
	// Execute the command and capture output and error.
	_, output, err = ExecuteCommandC(root, args...)

	if err != nil {
		// Wrap error with additional context.
		err = fmt.Errorf("failed to execute command: %w", err)
	}

	// Returning the output and any error encountered.
	return output, err
}
func ExecuteCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}
