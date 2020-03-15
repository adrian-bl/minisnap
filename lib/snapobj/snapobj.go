package snapobj

import (
	"fmt"
	"os"
	"strings"
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
	parts := strings.Split(fi.Name(), "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format")
	}

	xtype, err := toType(parts[0])
	if err != nil {
		return nil, err
	}
	xtime, err := time.Parse(time.RFC3339, parts[1])
	if err != nil {
		return nil, err
	}
	return &SnapObj{
		Type:  xtype,
		Epoch: xtime.UTC(),
	}, nil
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

func toType(s string) (Type, error) {
	switch s {
	case "hourly":
		return Hourly, nil
	case "daily":
		return Daily, nil
	case "weekly":
		return Weekly, nil
	case "monthly":
		return Monthly, nil
	case "yearly":
		return Yearly, nil
	default:
		return 0, fmt.Errorf("unknown type string")
	}
}
