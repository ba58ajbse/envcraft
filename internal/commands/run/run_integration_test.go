package run

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Run_Integration(t *testing.T) {
	tests := map[string]struct {
		input   string
		want    map[string]string
		wantErr bool
	}{
		"basic": {
			input: "testdata/.env",
			want:  map[string]string{"BASIC": "basic"},
		},
		"empty": {
			input: "testdata/.env",
			want:  map[string]string{"EMPTY": ""},
		},
		"empty simgle quotes": {
			input: "testdata/.env",
			want:  map[string]string{"EMPTY_SINGLE_QUOTES": ""},
		},
		"empty double quotes": {
			input: "testdata/.env",
			want:  map[string]string{"EMPTY_DOUBLE_QUOTES": ""},
		},
		"empty back quotes": {
			input: "testdata/.env",
			want:  map[string]string{"EMPTY_BACK_QUOTES": ""},
		},
		"single quotes": {
			input: "testdata/.env",
			want:  map[string]string{"SINGLE_QUOTES": "single_quotes"},
		},
		"single quotes spaced": {
			input: "testdata/.env",
			want:  map[string]string{"SINGLE_QUOTES_SPACED": "    single quotes    "},
		},
		"double quotes": {
			input: "testdata/.env",
			want:  map[string]string{"DOUBLE_QUOTES": "double_quotes"},
		},
		"double quotes spaced": {
			input: "testdata/.env",
			want:  map[string]string{"DOUBLE_QUOTES_SPACED": "    double quotes    "},
		},
		"double quotes inside single": {
			input:   "testdata/.env",
			want:    map[string]string{"DOUBLE_QUOTES_INSIDE_SINGLE": `double "quotes" work inside single quotes`},
			wantErr: false,
		},
		// "double quotes with no space backet": {
		// 	input:   "testdata/.env",
		// 	want:    map[string]string{"DOUBLE_QUOTES_WITH_NO_SPACE_BRACKET": "{ port: }"},
		// 	wantErr: false,
		// },
		"single quotes inside double": {
			input:   "testdata/.env",
			want:    map[string]string{"SINGLE_QUOTES_INSIDE_DOUBLE": "single 'quotes' work inside double quotes"},
			wantErr: false,
		},
		"backticks inside single": {
			input: "testdata/.env",
			want:  map[string]string{"BACKTICKS_INSIDE_SINGLE": "`backticks` work inside single quotes"},
		},

		"backticks inside double": {
			input: "testdata/.env",
			want:  map[string]string{"BACKTICKS_INSIDE_DOUBLE": "`backticks` work inside double quotes"},
		},

		"backticks": {
			input: "testdata/.env",
			want:  map[string]string{"BACKTICKS": "backticks"},
		},
		"backticks spaced": {
			input: "testdata/.env",
			want:  map[string]string{"BACKTICKS_SPACED": "    backticks    "},
		},

		"double quotes inside backticks": {
			input: "testdata/.env",
			want:  map[string]string{"DOUBLE_QUOTES_INSIDE_BACKTICKS": "double \"quotes\" work inside backticks"},
		},

		"single quotes inside backticks": {
			input: "testdata/.env",
			want:  map[string]string{"SINGLE_QUOTES_INSIDE_BACKTICKS": "single 'quotes' work inside backticks"},
		},

		"double and single quotes inside backticks": {
			input: "testdata/.env",
			want:  map[string]string{"DOUBLE_AND_SINGLE_QUOTES_INSIDE_BACKTICKS": "double \"quotes\" and single 'quotes' work inside backticks"},
		},
		"dont expand unquoted": {
			input: "testdata/.env",
			want:  map[string]string{"DONT_EXPAND_UNQUOTED": `dontexpand\nnewlines`},
		},
		"dont expand squoted": {
			input: "testdata/.env",
			want:  map[string]string{"DONT_EXPAND_SQUOTED": `dontexpand\nnewlines`},
		},
		"inline comments": {
			input: "testdata/.env",
			want:  map[string]string{"INLINE_COMMENTS": "inline comments"},
		},
		"inline comments single quotes": {
			input: "testdata/.env",
			want:  map[string]string{"INLINE_COMMENTS_SINGLE_QUOTES": "inline comments outside of #singlequotes"},
		},
		"inline comments double quotes": {
			input: "testdata/.env",
			want:  map[string]string{"INLINE_COMMENTS_DOUBLE_QUOTES": "inline comments outside of #doublequotes"},
		},
		"inline comments backticks": {
			input: "testdata/.env",
			want:  map[string]string{"INLINE_COMMENTS_BACKTICKS": "inline comments outside of #backticks"},
		},
		"inline comments space": {
			input: "testdata/.env",
			want:  map[string]string{"INLINE_COMMENTS_SPACE": "inline comments start with a"},
		},
		"equal signs": {
			input: "testdata/.env",
			want:  map[string]string{"EQUAL_SIGNS": "equals=="},
		},
		"retain inner quotes": {
			input: "testdata/.env",
			want:  map[string]string{"RETAIN_INNER_QUOTES": `{"foo": "bar"}`},
		},
		"retain inner quotes as string": {
			input: "testdata/.env",
			want:  map[string]string{"RETAIN_INNER_QUOTES_AS_STRING": `{"foo": "bar"}`},
		},
		"retain inner quotes as backticks": {
			input: "testdata/.env",
			want:  map[string]string{"RETAIN_INNER_QUOTES_AS_BACKTICKS": `{"foo": "bar's"}`},
		},
		"trim space from unquoted": {
			input: "testdata/.env",
			want:  map[string]string{"TRIM_SPACE_FROM_UNQUOTED": "some spaced out string"},
		},
		"username": {
			input: "testdata/.env",
			want:  map[string]string{"USERNAME": "therealnerdybeast@example.tld"},
		},
		"spaced key": {
			input: "testdata/.env",
			want:  map[string]string{"SPACED_KEY": "parsed"},
		},
	}

	// NOTE: Environment variables are set only in the subprocess, so no need to unset them after the test.
	// 実行するコマンド: env
	args := []string{"-f", "testdata/.env", "--", "env"}

	// コマンドの標準出力をキャプチャするためにexec.Commandを直接使う
	cmd := exec.Command(os.Args[0], append([]string{"-test.run=TestHelperProcess", "--"}, args...)...)
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		t.Fatalf("Run() error: %v\nOutput: %s", err, out.String())
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			output := out.String()
			for k, v := range tt.want {
				kv := k + "=" + v
				assert.Containsf(t, output, kv, "env %s=%s not found in output", k, v)
			}
		})
	}

	// 期待する環境変数
	// expects := map[string]string{
	// 	"BASIC": "basic",
	// 	// "BAR": "baz",
	// 	// 他のキーも必要に応じて追加
	// }

	// 出力に期待値が含まれているか確認
}

// テスト用のヘルパープロセス
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// Run() を呼び出す
	Run(os.Args[3:])
	os.Exit(0)
}
