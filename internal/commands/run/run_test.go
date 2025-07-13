package run

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	tests := map[string]struct {
		args    []string
		wantErr bool
	}{
		"success: env command with .env": {
			args:    []string{"-f", "testdata/test.env", "--", "env"},
			wantErr: false,
		},
		"error: missing file": {
			args:    []string{"-f", "testdata/notfound.env", "--", "env"},
			wantErr: true,
		},
		"error: missing command": {
			args:    []string{"-f", "testdata/test.env"},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := Run(tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
