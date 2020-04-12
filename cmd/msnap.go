package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adrian-bl/minisnap/lib/fs"
	"github.com/adrian-bl/minisnap/lib/policy"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION] vol [vol...]\n\n", os.Args[0])
		flag.PrintDefaults()
	}
}

var (
	dryRun   = flag.Bool("dry_run", false, "do not execute, just print what would be done")
	confFile = flag.String("config", "/etc/minisnap.conf", "path to configuration file")
)

func main() {
	flag.Parse()

	vols := flag.Args()
	if len(vols) == 0 {
		flag.Usage()
		xfail("\nNo volumes given, exiting")
	}

	conf, err := parseConfig(*confFile)
	if err != nil {
		xfail("failed to parse '%s': %v", *confFile, err)
	}

	for _, vol := range vols {
		vol = filepath.Clean(vol)
		vp, ok := conf[vol]
		if !ok {
			xfail(fmt.Sprintf("volume %s: not defined in config", vol))
		}
		p := &policy.Policy{
			Now:  time.Now(),
			Keep: vp,
		}
		if err := snapshot(vol, p, *dryRun); err != nil {
			xfail(fmt.Sprintf("volume %s: %v", vol, err))
		}
	}

}

// snapshot performs the snapshotting operation on the given volume.
func snapshot(vol string, p *policy.Policy, dryRun bool) error {
	fss, err := fs.ForVolume(vol, dryRun)
	if err != nil {
		return fmt.Errorf("failed to open volume: %v", err)
	}

	fmt.Printf("Working on %s\n", fss.Description())
	cur, err := fss.Gather()
	if err != nil {
		return fmt.Errorf("failed to gather current snapshots: %v", err)
	}

	plan, err := p.Plan(cur)
	if err != nil {
		return fmt.Errorf("could not construct a plan: %v", err)
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
		return fmt.Errorf("errors during create phase, refusing to enter delete phase")
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
		return fmt.Errorf("delete phase had errors")
	}

	return nil
}

func xfail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
