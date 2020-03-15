package policy

import (
	"fmt"
	"os"

	"github.com/adrian-bl/minisnap/lib/snapobj"
)

func Enforce(input []os.FileInfo) error {
	batch := make(map[string][]snapobj.SnapObj)
	for _, fi := range input {
		so, err := snapobj.FromFileInfo(fi)
		if err != nil {
			return err
		}
		k := so.String()
		batch[k] = append(batch[k], so)
	}
	fmt.Printf(">> %+v\n", batch)
	return nil
}
