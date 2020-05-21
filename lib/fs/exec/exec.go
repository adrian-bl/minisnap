package exec

import (
	"fmt"
	"os"
	"os/exec"
)

type Exec struct {
	// Do not execute commands, just print what would be done.
	DryRun bool
	// Print command which was executed.
	Verbose bool
}

func (e *Exec) Execute(name string, args ...string) error {
	if e.DryRun {
		fmt.Printf("Would execute: %s %q\n", name, args)
		return nil
	}
	if e.Verbose {
		fmt.Printf("Executing %s %q\n", name, args)
	}

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
