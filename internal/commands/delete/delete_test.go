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
