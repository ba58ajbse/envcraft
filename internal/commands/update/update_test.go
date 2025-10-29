package update

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_makeNewLines(t *testing.T) {
	tests := map[string]struct {
		orgLines []string
		key      string
		value    string
		want     []string
		wantErr  bool
	}{
		"update existing key": {
			orgLines: []string{"FOO=\"bar\"\n", "BAR=\"baz\"\n"},
			key:      "FOO",
			value:    "newval",
			want:     []string{"FOO=\"newval\"\n", "BAR=\"baz\"\n"},
			wantErr:  false,
		},
		"no matching key": {
			orgLines: []string{"FOO=\"bar\"\n"},
			key:      "BAZ",
			value:    "val",
			want:     []string{"FOO=\"bar\"\n"},
			wantErr:  true,
		},
		"empty orgLines": {
			orgLines: []string{},
			key:      "FOO",
			value:    "bar",
			want:     []string{},
			wantErr:  true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cmd := &UpdateCmd{
				Options: UpdateOptions{
					Key:   tt.key,
					Value: tt.value,
				},
				OrgLines: tt.orgLines,
			}
			got, err := cmd.makeNewLines()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestParseUpdateOptions(t *testing.T) {
	tests := map[string]struct {
		opts    []string
		want    *UpdateOptions
		wantErr bool
	}{
		"key value before flags": {
			opts:    []string{"KEY", "VALUE", "-f", "test.env"},
			want:    &UpdateOptions{Key: "KEY", Value: "VALUE", FilePath: "test.env"},
			wantErr: false,
		},
		"flags before key value": {
			opts:    []string{"-f", "test.env", "KEY", "VALUE"},
			want:    &UpdateOptions{Key: "KEY", Value: "VALUE", FilePath: "test.env"},
			wantErr: false,
		},
		"missing value": {
			opts:    []string{"KEY", "-f", "test.env"},
			want:    nil,
			wantErr: true,
		},
		"missing file": {
			opts:    []string{"KEY", "VALUE"},
			want:    nil,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ParseUpdateOptions(tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
