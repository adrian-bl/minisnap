package fs

import (
	"fmt"
	"syscall"

	"github.com/adrian-bl/minisnap/lib/fs/btrfs"
	"github.com/adrian-bl/minisnap/lib/fs/exec"
	"github.com/adrian-bl/minisnap/lib/snapobj"
)

const (
	fsBtrfs = 0x9123683E
)

type FsSnap interface {
	Description() string
	Gather() ([]*snapobj.SnapObj, error)
	Create(*snapobj.SnapObj) error
	Delete(*snapobj.SnapObj) error
}

func ForVolume(path string) (FsSnap, error) {
	buf := &syscall.Statfs_t{}
	if err := syscall.Statfs(path, buf); err != nil {
		return nil, err
	}

	e := &exec.Exec{DryRun: false}
	switch buf.Type {
	case fsBtrfs:
		return btrfs.New(path, ".snapshots", e), nil
	}
	return nil, fmt.Errorf("Unknown fstype: %X", buf.Type)
}
