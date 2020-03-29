package policy

import (
	"sort"
	"time"

	"github.com/adrian-bl/minisnap/lib/snapobj"
)

type Policy struct {
	Now  time.Time
	Keep map[snapobj.Type]int
}

type Plan struct {
	// Whether or not this is a delete operation.
	Delete bool
	// Name of the affected snapshot.
	Target *snapobj.SnapObj
}

func (p Policy) Plan(s []*snapobj.SnapObj) ([]*Plan, error) {
	pl := make([]*Plan, 0)
	catalog := make(map[snapobj.Type][]*snapobj.SnapObj)

	// First, separate all snapshots by type and sort them by time.
	for _, o := range s {
		catalog[o.Type] = append(catalog[o.Type], o)
	}
	for _, o := range catalog {
		sort.Sort(bySnap(o))
	}

	// Second: check all types to see if we need to create a new snapshot.
	for t, o := range catalog {
		if p.Keep[t] < 1 {
			// skip as there should be no snapshots of this type.
			continue
		}

		var current bool
		for _, x := range o {
			if x.IsCurrent(p.Now) {
				current = true
				break
			}
		}
		if current {
			continue
		}
		// no current snapshot? Add it to our plan AND add a fake object to
		// the catalog to make its length match 'the future'.
		pl = append(pl, &Plan{Target: &snapobj.SnapObj{Epoch: p.Now, Type: t}})
		catalog[t] = append(catalog[t], &snapobj.SnapObj{})
	}

	for t, o := range catalog {
		k := p.Keep[t]
		for _, x := range o {
			if len(o) <= k {
				break
			}
			// If not expired and non-zero, add it to kill list.
			if !x.IsCurrent(p.Now) && !x.Epoch.IsZero() {
				k++
				pl = append(pl, &Plan{Delete: true, Target: x})
			}
		}
	}
	return pl, nil
}
