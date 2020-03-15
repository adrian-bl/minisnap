package snapobj

import (
	"testing"
	"time"
)

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
