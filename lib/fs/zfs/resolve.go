package zfs

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
)

var reSplitMounts = regexp.MustCompile(`^(.+)\s+(.+)$`)

// resolveZfsMount finds the volume name of a given mountpoint.
func resolveZfsMount(mp string) (string, error) {
	cmd := exec.Command("zfs", "list", "-H", "-o", "name,mountpoint")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	for _, l := range bytes.Split(out, []byte{'\n'}) {
		if len(l) == 0 {
			continue
		}

		m := reSplitMounts.FindSubmatch(l)
		if len(m) != 3 {
			return "", fmt.Errorf("error parsing line %s", string(l))
		}
		if string(m[2]) == mp {
			return string(m[1]), nil
		}
	}
	return "", fmt.Errorf("failed to find name of mountpoint %s", mp)
}
