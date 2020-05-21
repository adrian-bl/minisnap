package fs

import (
	"fmt"
	"syscall"

	"github.com/adrian-bl/minisnap/lib/fs/btrfs"
	"github.com/adrian-bl/minisnap/lib/fs/exec"
	"github.com/adrian-bl/minisnap/lib/fs/zfs"
	"github.com/adrian-bl/minisnap/lib/opts"
	"github.com/adrian-bl/minisnap/lib/snapobj"
)

const (
	fsBtrfs = 0x9123683E
	fsZFS   = 0x2FC12FC1
)

type FsSnap interface {
	Description() string
	Gather() ([]*snapobj.SnapObj, error)
	Create(*snapobj.SnapObj) error
	Delete(*snapobj.SnapObj) error
}

func ForVolume(path string, vopts opts.VolOptions, dryRun, verbose bool) (FsSnap, error) {
	buf := &syscall.Statfs_t{}
	if err := syscall.Statfs(path, buf); err != nil {
		return nil, err
	}

	e := &exec.Exec{
		DryRun:  dryRun,
		Verbose: verbose,
	}
	switch buf.Type {
	case fsBtrfs:
		if vopts.Recursive {
			return nil, fmt.Errorf("btrfs does not support recursive snapshots")
		}
		return btrfs.New(path, ".snapshots", e), nil
	case fsZFS:
		return zfs.New(path, "msnap_", e, vopts.Recursive)
	}
	return nil, fmt.Errorf("Unknown fstype: %X", buf.Type)
}
