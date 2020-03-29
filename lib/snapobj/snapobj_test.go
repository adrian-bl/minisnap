package snapobj

import (
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type tinfo struct {
	name  string
	isdir bool
}

func (x tinfo) Name() string {
	return x.name
}

func (x tinfo) IsDir() bool {
	return x.isdir
}

func (x tinfo) ModTime() time.Time {
	return time.Unix(3, 0)
}

func (x tinfo) Mode() os.FileMode {
	return 0
}

func (x tinfo) Size() int64 {
	return 0
}

func (x tinfo) Sys() interface{} {
	return nil
}

func TestFileName(t *testing.T) {
	input := []struct {
		input SnapObj
		want  string
	}{
		{
			input: SnapObj{
				Epoch: time.Unix(0, 0),
				Type:  Hourly,
			},
			want: "hourly@1970-01-01T00:00:00Z",
		},
		{
			input: SnapObj{
				Epoch: time.Unix(86400, 0),
				Type:  Daily,
			},
			want: "daily@1970-01-02T00:00:00Z",
		},
		{
			input: SnapObj{
				Epoch: time.Unix(86400*3, 0),
				Type:  Weekly,
			},
			want: "weekly@1970-01-04T00:00:00Z",
		},
		{
			input: SnapObj{
				Epoch: time.Unix(0, 0),
				Type:  Monthly,
			},
			want: "monthly@1970-01-01T00:00:00Z",
		},
		{
			input: SnapObj{
				Epoch: time.Unix(853520053, 0),
				Type:  Yearly,
			},
			want: "yearly@1997-01-17T16:54:13Z",
		},
	}

	for _, tt := range input {
		got := tt.input.FileName()
		if got != tt.want {
			t.Errorf("FileName(%+v) = %s, want %s", tt.input, got, tt.want)
		}
	}
}

func TestFromFileInfo(t *testing.T) {
	input := []struct {
		name    string
		input   tinfo
		wantErr bool
		want    *SnapObj
	}{
		{
			name:    "invalid file",
			input:   tinfo{name: "foo-"},
			wantErr: true,
		},
		{
			name:    "nodir",
			input:   tinfo{name: "yearly@1997-01-17T16:54:13Z", isdir: false},
			wantErr: true,
		},
		{
			name:  "yearly@1997-01-17T16:54:13Z",
			input: tinfo{name: "yearly@1997-01-17T16:54:13Z", isdir: true},
			want: &SnapObj{
				Type:  Yearly,
				Epoch: time.Unix(853520053, 0),
			},
		},
		{
			name:  "monthly@1997-01-17T16:54:14Z",
			input: tinfo{name: "monthly@1997-01-17T16:54:14Z", isdir: true},
			want: &SnapObj{
				Type:  Monthly,
				Epoch: time.Unix(853520054, 0),
			},
		},
		{
			name:  "weekly@1997-01-17T16:54:13Z",
			input: tinfo{name: "weekly@1997-01-17T16:54:13Z", isdir: true},
			want: &SnapObj{
				Type:  Weekly,
				Epoch: time.Unix(853520053, 0),
			},
		},
		{
			name:  "daily@1997-01-17T16:54:13Z",
			input: tinfo{name: "daily@1997-01-17T16:54:13Z", isdir: true},
			want: &SnapObj{
				Type:  Daily,
				Epoch: time.Unix(853520053, 0),
			},
		},
		{
			name:  "hourly@1997-01-17T16:54:21Z",
			input: tinfo{name: "hourly@1997-01-17T16:54:21Z", isdir: true},
			want: &SnapObj{
				Type:  Hourly,
				Epoch: time.Unix(853520061, 0),
			},
		},
		{
			name:    "taeglich@1997-01-17T16:54:21Z",
			input:   tinfo{name: "taeglich@1997-01-17T16:54:21Z", isdir: true},
			wantErr: true,
		},
	}

	for _, tt := range input {
		got, err := FromFileInfo(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("FromFileInput(%s) = _, nil, wanted err", tt.name)
			}
			continue
		}
		if !tt.wantErr && err != nil {
			t.Errorf("FromFileInput(%s) = _, %v, wanted nil", tt.name, err)
		}
		if diff := cmp.Diff(got, tt.want); diff != "" {
			t.Errorf("FromFileInput(%s) mismatch (-want +got)\n%s", tt.name, diff)
		}
	}
}

func TestIsCurrent(t *testing.T) {
	now := time.Unix(80061, 0)
	input := []struct {
		snap *SnapObj
		want bool
	}{
		{
			snap: &SnapObj{
				Type:  Weekly,
				Epoch: time.Unix(80000, 0),
			},
			want: true,
		},
		{
			snap: &SnapObj{
				Type:  Hourly,
				Epoch: time.Unix(90000, 0),
			},
			want: true,
		},
		{
			snap: &SnapObj{
				Type:  Minutely,
				Epoch: time.Unix(80002, 0),
			},
			want: true,
		},
		{
			snap: &SnapObj{
				Type:  Minutely,
				Epoch: time.Unix(80001, 0),
			},
			want: false,
		},
	}

	for _, tt := range input {
		got := tt.snap.IsCurrent(now)
		if got != tt.want {
			t.Errorf("IsCurrent(%v) = %v, want %v", tt.snap, got, tt.want)
		}
	}
}
