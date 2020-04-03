package exec

import (
	"fmt"
	"os"
	"os/exec"
)

type Exec struct {
	DryRun bool
}

func (e *Exec) Execute(name string, args ...string) error {
	if e.DryRun {
		fmt.Printf("Would execute: %s %q\n", name, args)
		return nil
	}
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
