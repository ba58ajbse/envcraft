package delete

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_makeNewLines(t *testing.T) {
	tests := map[string]struct {
		orgLines []string
		key      string
		want     []string
		wantErr  bool
	}{
		"delete existing key": {
			orgLines: []string{"FOO=\"bar\"\n", "BAR=\"baz\""},
			key:      "FOO",
			want:     []string{"BAR=\"baz\""},
			wantErr:  false,
		},
		"delete existing key 2": {
			orgLines: []string{"FOO=\"bar\"\n", "BAR=\"baz\"\n", "FIZZ=\"bazz\""},
			key:      "BAR",
			want:     []string{"FOO=\"bar\"\n", "FIZZ=\"bazz\""},
			wantErr:  false,
		},
		"delete existing include empty line": {
			orgLines: []string{"FOO=\"bar\"\n", "BAR=\"baz\"\n", "", "FIZZ=\"bazz\""},
			key:      "BAR",
			want:     []string{"FOO=\"bar\"\n", "", "FIZZ=\"bazz\""},
			wantErr:  false,
		},
		"delete last line": {
			orgLines: []string{"FOO=\"bar\"\n", "BAR=\"baz\"\n", "FIZZ=\"bazz\""},
			key:      "FIZZ",
			want:     []string{"FOO=\"bar\"\n", "BAR=\"baz\""},
			wantErr:  false,
		},
		"no matching key": {
			orgLines: []string{"FOO=\"bar\"\n"},
			key:      "BAZ",
			want:     []string{"FOO=\"bar\"\n"},
			wantErr:  true,
		},
		"empty orgLines": {
			orgLines: []string{},
			key:      "FOO",
			want:     []string{},
			wantErr:  true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cmd := &DeleteCmd{
				Options: DeleteOptions{
					Key: tt.key,
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

func TestParseDeleteOptions(t *testing.T) {
	tests := map[string]struct {
		opts    []string
		want    *DeleteOptions
		wantErr bool
	}{
		"key before flags": {
			opts:    []string{"KEY", "-f", "test.env"},
			want:    &DeleteOptions{Key: "KEY", FilePath: "test.env"},
			wantErr: false,
		},
		"flags before key": {
			opts:    []string{"-f", "test.env", "KEY"},
			want:    &DeleteOptions{Key: "KEY", FilePath: "test.env"},
			wantErr: false,
		},
		"missing key": {
			opts:    []string{"-f", "test.env"},
			want:    nil,
			wantErr: true,
		},
		"missing file": {
			opts:    []string{"KEY"},
			want:    nil,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ParseDeleteOptions(tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
