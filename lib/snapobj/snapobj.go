package snapobj

import (
	"fmt"
	"os"
	"time"
)

type Type int

const (
	Hourly  = 3600
	Daily   = 86400
	Weekly  = 86400 * 7
	Monthly = 2592000
	Yearly  = 86400 * 360
)

type SnapObj struct {
	Epoch time.Time
	Type  Type
}

func FromFileInfo(fi os.FileInfo) (*SnapObj, error) {
	if !fi.IsDir() {
		return nil, fmt.Errorf("not a directory")
	}
	return nil, fmt.Errorf("not implemented yet")
}

func (so SnapObj) FileName() string {
	return fmt.Sprintf("%s@%s", so.Type.String(), so.Epoch.UTC().Format(time.RFC3339))
}

func (t Type) String() string {
	switch t {
	case Hourly:
		return "hourly"
	case Daily:
		return "daily"
	case Weekly:
		return "weekly"
	case Monthly:
		return "monthly"
	case Yearly:
		return "yearly"
	default:
		return fmt.Sprintf("uk-%d", t)
	}
}
