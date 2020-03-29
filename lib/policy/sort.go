package policy

import (
	"github.com/adrian-bl/minisnap/lib/snapobj"
)

type bySnap []*snapobj.SnapObj

func (b bySnap) Len() int {
	return len(b)
}

func (b bySnap) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b bySnap) Less(i, j int) bool {
	return b[i].Epoch.Before(b[j].Epoch)
}
