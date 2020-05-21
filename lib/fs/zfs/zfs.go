package zfs

import (
	"bytes"
	"fmt"
	oe "os/exec"
	"strings"

	"github.com/adrian-bl/minisnap/lib/fs/exec"
	"github.com/adrian-bl/minisnap/lib/snapobj"
)

type Zfs struct {
	name       string
	mountpoint string
	snapprefix string
	recursive  bool
	exec       *exec.Exec
}

func New(path, sprefix string, e *exec.Exec, recursive bool) (*Zfs, error) {
	name, err := resolveZfsMount(path)
	if err != nil {
		return nil, err
	}
	z := &Zfs{
		mountpoint: path,
		name:       name,
		snapprefix: sprefix,
		recursive:  recursive,
		exec:       e,
	}
	return z, nil
}

func (z *Zfs) Description() string {
	return fmt.Sprintf("%s using ZFS vol %s", z.mountpoint, z.name)
}

func (z *Zfs) Gather() ([]*snapobj.SnapObj, error) {
	cmd := oe.Command("zfs", "list", "-H", "-t", "snapshot", "-o", "name", z.name)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	e := []*snapobj.SnapObj{}
	pfx := []byte(fmt.Sprintf("%s@%s", z.name, z.snapprefix))
	for _, l := range bytes.Split(out, []byte{'\n'}) {
		if len(l) == 0 {
			continue
		}
		if bytes.Compare(l, pfx) < 1 {
			// not a snapshot managed by us.
			continue
		}
		// convert back into something snapobj understands: must agree with snapName().
		id := strings.Replace(string(l[len(pfx):]), "::", "@", 1)
		so, err := snapobj.FromString(id)
		if err != nil {
			return nil, err
		}
		e = append(e, so)
	}
	return e, nil
}

func (z *Zfs) Create(s *snapobj.SnapObj) error {
	args := []string{"snapshot"}
	if z.recursive {
		args = append(args, "-r")
	}
	args = append(args, fmt.Sprintf("%s@%s", z.name, z.snapName(s)))
	return z.exec.Execute("zfs", args...)
}

func (z *Zfs) Delete(s *snapobj.SnapObj) error {
	args := []string{"destroy"}
	if z.recursive {
		args = append(args, "-r")
	}
	args = append(args, fmt.Sprintf("%s@%s", z.name, z.snapName(s)))
	return z.exec.Execute("zfs", args...)
}

func (z *Zfs) snapName(s *snapobj.SnapObj) string {
	// zfs can not contain @ signs in snapshot names.
	sname := strings.Replace(s.FileName(), "@", "::", 1)
	return fmt.Sprintf("%s%s", z.snapprefix, sname)
}
