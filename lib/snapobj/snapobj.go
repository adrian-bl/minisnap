package snapobj

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Type int

const (
	Minutely = 60
	Hourly   = 3600
	Daily    = 86400
	Weekly   = 86400 * 7
	Monthly  = 2592000
	Yearly   = 86400 * 360
)

type SnapObj struct {
	Epoch time.Time
	Type  Type
}

// FromFileInfo returns a snap object from a os.FileInfo.
func FromFileInfo(fi os.FileInfo) (*SnapObj, error) {
	if !fi.IsDir() {
		return nil, fmt.Errorf("not a directory")
	}
	return FromString(fi.Name())
}

// FromString returns a snap object from a bare string.
func FromString(s string) (*SnapObj, error) {
	parts := strings.Split(s, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format")
	}

	xtype, err := ToType(parts[0])
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

// FileName returns the file basename to use.
func (so SnapObj) FileName() string {
	return fmt.Sprintf("%s@%s", so.Type.String(), so.Epoch.UTC().Format(time.RFC3339))
}

// SameType returns true if the compared snap objects are of the same snapshot type.
func (so SnapObj) SameType(oo SnapObj) bool {
	return so.Type == oo.Type
}

func (so SnapObj) IsCurrent(t time.Time) bool {
	exp := so.Epoch.Add(time.Duration(so.Type) * time.Second)
	return exp.After(t)
}

// String returns the string version of this type.
func (t Type) String() string {
	switch t {
	case Minutely:
		return "minutely"
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

// ToType takes a string and returns a matching Type.
func ToType(s string) (Type, error) {
	switch s {
	case "minutely":
		return Minutely, nil
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
