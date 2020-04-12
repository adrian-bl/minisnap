package btrfs

import (
	"fmt"
	"io"
	"os"

	"github.com/adrian-bl/minisnap/lib/snapobj"
)

type exec interface {
	Execute(name string, args ...string) error
}

type Btrfs struct {
	path    string
	snapdir string
	exec    exec
}

func New(path, snapdir string, exec exec) *Btrfs {
	return &Btrfs{path: path, snapdir: snapdir, exec: exec}
}

func (b *Btrfs) wdir() string {
	p := b.path
	if len(p) > 0 && p[len(p)-1] != '/' {
		p += "/"
	}
	p += b.snapdir
	return p
}

func (b *Btrfs) Description() string {
	return fmt.Sprintf("%s using btrfs", b.wdir())
}

func (b *Btrfs) Gather() ([]*snapobj.SnapObj, error) {
	g := make([]*snapobj.SnapObj, 0)
	fh, err := os.Open(b.wdir())
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	for {
		fi, err := fh.Readdir(3)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		for _, e := range fi {
			so, err := snapobj.FromFileInfo(e)
			if err == nil {
				g = append(g, so)
			}
		}
	}
	return g, nil
}

func (b *Btrfs) Create(s *snapobj.SnapObj) error {
	return b.exec.Execute("btrfs", "subvol", "snapshot", b.path, fmt.Sprintf("%s/%s", b.wdir(), s.FileName()))
}

func (b *Btrfs) Delete(s *snapobj.SnapObj) error {
	return b.exec.Execute("btrfs", "subvol", "delete", fmt.Sprintf("%s/%s", b.wdir(), s.FileName()))
}
