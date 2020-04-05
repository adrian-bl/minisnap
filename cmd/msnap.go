package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/adrian-bl/minisnap/lib/fs"
	"github.com/adrian-bl/minisnap/lib/policy"
	"github.com/adrian-bl/minisnap/lib/snapobj"
)

var (
	dryRun = flag.Bool("dry_run", false, "do not execute, just print what would be done")
)

func main() {
	flag.Parse()
	fakePol := policy.Policy{
		Now: time.Now(),
		Keep: map[snapobj.Type]int{
			snapobj.Minutely: 3,
			snapobj.Daily:    4,
		},
	}

	vol := "/"
	fss, err := fs.ForVolume(vol, *dryRun)
	if err != nil {
		xfail("failed to open volume %s: %v", vol, err)
	}
	fmt.Printf("Working on %s\n", fss.Description())

	cur, err := fss.Gather()
	if err != nil {
		xfail("failed to gather current snapshots: %v", err)
	}

	plan, err := fakePol.Plan(cur)
	if err != nil {
		panic(err)
	}

	var failed bool
	for _, o := range plan {
		if !o.Delete {
			if err := fss.Create(o.Target); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating %s: %v\n", o.Target, err)
				failed = true
			}
		}
	}
	if failed {
		xfail("Create phase had errors, refusing to delete anything")
	}

	for _, o := range plan {
		if o.Delete {
			if err := fss.Delete(o.Target); err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting %s: %v\n", o.Target, err)
				failed = true
			}
		}
	}
	if failed {
		xfail("Delete phase had errors")
	}
}

func xfail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
