package policy

import (
	"testing"
	"time"

	"github.com/adrian-bl/minisnap/lib/snapobj"

	"github.com/google/go-cmp/cmp"
)

func TestPlan(t *testing.T) {
	sof := func(s string) *snapobj.SnapObj {
		v, err := snapobj.FromString(s)
		if err != nil {
			panic(err)
		}
		return v
	}

	now := time.Unix(90000123, 0).UTC()

	input := []struct {
		name   string
		policy *Policy
		input  []*snapobj.SnapObj
		want   []*Plan
	}{
		{
			name:   "empty",
			policy: &Policy{},
			input:  []*snapobj.SnapObj{},
			want:   []*Plan{},
		},
		{
			name: "create",
			policy: &Policy{
				Now: now,
				Keep: map[snapobj.Type]int{
					snapobj.Hourly: 2,
					snapobj.Daily:  1,
				},
			},
			input: []*snapobj.SnapObj{
				sof("hourly@1972-11-07T10:54:13Z"),
				sof("daily@1972-11-07T16:54:13Z"),
			},
			want: []*Plan{
				{
					Target: sof("hourly@1972-11-07T16:02:03Z"),
				},
			},
		},
		{
			name: "create 999",
			policy: &Policy{
				Now: now,
				Keep: map[snapobj.Type]int{
					snapobj.Hourly: 999,
					snapobj.Daily:  1,
				},
			},
			input: []*snapobj.SnapObj{
				sof("hourly@1972-11-07T10:54:13Z"),
				sof("daily@1972-11-07T16:54:13Z"),
			},
			want: []*Plan{
				{
					Target: sof("hourly@1972-11-07T16:02:03Z"),
				},
			},
		},
		{
			name: "wipe all",
			policy: &Policy{
				Now: now,
			},
			input: []*snapobj.SnapObj{
				sof("hourly@1900-01-17T16:54:13Z"),
				sof("daily@1900-01-17T16:54:13Z"),
			},
			want: []*Plan{
				{
					Delete: true,
					Target: sof("hourly@1900-01-17T16:54:13Z"),
				},
				{
					Delete: true,
					Target: sof("daily@1900-01-17T16:54:13Z"),
				},
			},
		},
		{
			name: "wipe two",
			policy: &Policy{
				Now: now,
				Keep: map[snapobj.Type]int{
					snapobj.Hourly: 2,
				},
			},
			input: []*snapobj.SnapObj{
				sof("hourly@1972-11-07T14:50:13Z"),
				sof("hourly@1972-11-07T17:54:13Z"),
				sof("hourly@1972-11-07T13:54:13Z"),
				sof("hourly@1972-11-07T19:54:13Z"),
			},
			want: []*Plan{
				{
					Delete: true,
					Target: sof("hourly@1972-11-07T13:54:13Z"),
				},
				{
					Delete: true,
					Target: sof("hourly@1972-11-07T14:50:13Z"),
				},
			},
		},
		{
			name: "wipe non current",
			policy: &Policy{
				Now: now,
				Keep: map[snapobj.Type]int{
					snapobj.Hourly: 0,
				},
			},
			input: []*snapobj.SnapObj{
				sof("hourly@1972-11-07T00:54:13Z"),
				sof("hourly@1972-11-07T01:54:13Z"),
				sof("hourly@1972-11-07T15:54:13Z"), // stays.
				sof("hourly@1972-11-07T03:54:13Z"),
			},
			want: []*Plan{
				{
					Delete: true,
					Target: sof("hourly@1972-11-07T00:54:13Z"),
				},
				{
					Delete: true,
					Target: sof("hourly@1972-11-07T01:54:13Z"),
				},
				{
					Delete: true,
					Target: sof("hourly@1972-11-07T03:54:13Z"),
				},
			},
		},
		{
			name: "wipe non current #2",
			policy: &Policy{
				Now: now,
				Keep: map[snapobj.Type]int{
					snapobj.Hourly: 1,
				},
			},
			input: []*snapobj.SnapObj{
				sof("hourly@1972-11-07T00:54:13Z"),
				sof("hourly@1972-11-07T01:54:13Z"),
				sof("hourly@1972-11-07T15:54:13Z"), // stays.
				sof("hourly@1972-11-07T03:54:13Z"),
			},
			want: []*Plan{
				{
					Delete: true,
					Target: sof("hourly@1972-11-07T00:54:13Z"),
				},
				{
					Delete: true,
					Target: sof("hourly@1972-11-07T01:54:13Z"),
				},
				{
					Delete: true,
					Target: sof("hourly@1972-11-07T03:54:13Z"),
				},
			},
		},
		{
			name: "wipe non current, keep 2",
			policy: &Policy{
				Now: now,
				Keep: map[snapobj.Type]int{
					snapobj.Hourly: 2,
				},
			},
			input: []*snapobj.SnapObj{
				sof("hourly@1972-11-07T00:54:13Z"),
				sof("hourly@1972-11-07T01:54:13Z"),
				sof("hourly@1972-11-07T15:54:13Z"), // stays.
				sof("hourly@1972-11-07T03:54:13Z"), // expired but kept.
			},
			want: []*Plan{
				{
					Delete: true,
					Target: sof("hourly@1972-11-07T00:54:13Z"),
				},
				{
					Delete: true,
					Target: sof("hourly@1972-11-07T01:54:13Z"),
				},
			},
		},
		{
			name: "drop 2, keep 1, create 1",
			policy: &Policy{
				Now: now,
				Keep: map[snapobj.Type]int{
					snapobj.Hourly: 2,
				},
			},
			input: []*snapobj.SnapObj{
				sof("hourly@1972-11-07T00:54:13Z"),
				sof("hourly@1972-11-07T01:54:13Z"),
				sof("hourly@1972-11-07T03:54:13Z"), // expired but kept.
			},
			want: []*Plan{
				{
					Delete: false,
					Target: sof("hourly@1972-11-07T16:02:03Z"),
				},
				{
					Delete: true,
					Target: sof("hourly@1972-11-07T00:54:13Z"),
				},
				{
					Delete: true,
					Target: sof("hourly@1972-11-07T01:54:13Z"),
				},
			},
		},
		{
			name: "create 1, drop 3",
			policy: &Policy{
				Now: now,
				Keep: map[snapobj.Type]int{
					snapobj.Hourly: 1,
				},
			},
			input: []*snapobj.SnapObj{
				sof("hourly@1972-11-07T00:54:13Z"),
				sof("hourly@1972-11-07T01:54:13Z"),
				sof("hourly@1972-11-07T03:54:13Z"),
			},
			want: []*Plan{
				{
					Delete: false,
					Target: sof("hourly@1972-11-07T16:02:03Z"),
				},
				{
					Delete: true,
					Target: sof("hourly@1972-11-07T00:54:13Z"),
				},
				{
					Delete: true,
					Target: sof("hourly@1972-11-07T01:54:13Z"),
				},
				{
					Delete: true,
					Target: sof("hourly@1972-11-07T03:54:13Z"),
				},
			},
		},
		{
			name: "simple swap",
			policy: &Policy{
				Now: now,
				Keep: map[snapobj.Type]int{
					snapobj.Hourly: 1,
				},
			},
			input: []*snapobj.SnapObj{
				sof("hourly@1972-11-07T00:54:13Z"),
			},
			want: []*Plan{
				{
					Delete: false,
					Target: sof("hourly@1972-11-07T16:02:03Z"),
				},
				{
					Delete: true,
					Target: sof("hourly@1972-11-07T00:54:13Z"),
				},
			},
		},
	}

	for _, tt := range input {
		got, err := tt.policy.Plan(tt.input)
		if err != nil {
			t.Errorf("Plan(%s) = _, %v, want nil err", tt.name, err)
		}
		if diff := cmp.Diff(got, tt.want); diff != "" {
			t.Errorf("Plan(%s) mismatch (-want +got)\n%s", tt.name, diff)
		}
	}
}
