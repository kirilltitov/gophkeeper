package utils

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion_Print(t *testing.T) {
	type fields struct {
		BuildVersion string
		BuildDate    string
		BuildCommit  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Filled",
			fields: fields{
				BuildVersion: "13.3.7",
				BuildDate:    "09.03.1989",
				BuildCommit:  "deadbeef",
			},
			want: "Build version: 13.3.7\nBuild date: 09.03.1989\nBuild commit: deadbeef\n",
		},
		{
			name: "Empty",
			fields: fields{
				BuildVersion: "",
				BuildDate:    "",
				BuildCommit:  "",
			},
			want: "Build version: N/A\nBuild date: N/A\nBuild commit: N/A\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Version{
				BuildVersion: tt.fields.BuildVersion,
				BuildDate:    tt.fields.BuildDate,
				BuildCommit:  tt.fields.BuildCommit,
			}
			w := &bytes.Buffer{}
			v.Print(w)
			require.Equal(t, tt.want, w.String())
		})
	}
}
