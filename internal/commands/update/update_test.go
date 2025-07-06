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
