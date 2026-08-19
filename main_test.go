package main

import "testing"

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name      string
		argument  string
		want      string
		wantError bool
	}{
		{name: "token", argument: "1a2b3c4d5e6f", want: "1a2b3c4d5e6f"},
		{name: "token pasted with whitespace", argument: "  1a2b3c4d5e6f\n", want: "1a2b3c4d5e6f"},
		{name: "empty argument", argument: "", wantError: true},
		{name: "spaces only", argument: "   ", wantError: true},
		{name: "newline only", argument: "\n", wantError: true},
		{name: "tab only", argument: "\t", wantError: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := validateToken(test.argument)

			if test.wantError {
				if err == nil {
					t.Fatalf("got token %q and no error, want an error", got)
				}

				return
			}

			if err != nil {
				t.Fatalf("got error %v, want token %q", err, test.want)
			}

			if got != test.want {
				t.Errorf("got %q, want %q", got, test.want)
			}
		})
	}
}
