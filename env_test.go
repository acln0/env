// Copyright 2019 Andrei Tudor Călin
//
// Permission to use, copy, modify, and/or distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package env_test

import (
	"fmt"
	"os"
	"testing"

	"acln.ro/env"

	"github.com/google/go-cmp/cmp"
)

func TestString(t *testing.T) {
	tests := []struct {
		m    env.Map
		want string
	}{
		{
			m:    env.Map{},
			want: "",
		},
		{
			m:    env.Map{"k": "v"},
			want: "k=v",
		},
		{
			m:    env.Map{"k": ""},
			want: "k=",
		},
		{
			m:    env.Map{"FOO": "x", "BAR": "y"},
			want: "BAR=y FOO=x",
		},
		{
			m:    env.Map{"FOO": "x", "BAR": "", "BAZ": "z"},
			want: "BAR= BAZ=z FOO=x",
		},
	}
	for _, tt := range tests {
		got := tt.m.String()
		if diff := cmp.Diff(got, tt.want); diff != "" {
			stringErrorf(t, tt.m, got, tt.want, diff)
		}
	}
}

func stringErrorf(t *testing.T, m env.Map, got, want, diff string) {
	t.Helper()
	basicm := map[string]string(m)
	t.Errorf("%#v.Encode() = %q, want %q: %s", basicm, got, want, diff)
}

func TestFormat(t *testing.T) {
	tests := []struct {
		m      env.Map
		format string
		want   string
	}{
		{
			m:      env.Map{},
			format: "%v",
			want:   "",
		},
		{
			m:      env.Map{},
			format: "%+v",
			want:   "",
		},
		{
			m:      env.Map{"k": "v"},
			format: "%v",
			want:   "k=v",
		},
		{
			m:      env.Map{"k": "v"},
			format: "%+v",
			want:   "k=v",
		},
		{
			m:      env.Map{"k": ""},
			format: "%v",
			want:   "k=",
		},
		{
			m:      env.Map{"k": ""},
			format: "%+v",
			want:   "k=",
		},
		{
			m:      env.Map{"FOO": "x", "BAR": "y"},
			format: "%v",
			want:   "BAR=y FOO=x",
		},
		{
			m:      env.Map{"FOO": "x", "BAR": "y"},
			format: "%+v",
			want:   "BAR=y\nFOO=x",
		},
		{
			m:      env.Map{"FOO": "x", "BAR": "", "BAZ": "z"},
			format: "%v",
			want:   "BAR= BAZ=z FOO=x",
		},
		{
			m:      env.Map{"FOO": "x", "BAR": "", "BAZ": "z"},
			format: "%+v",
			want:   "BAR=\nBAZ=z\nFOO=x",
		},
		{
			m:      env.Map{"FOO": "x", "BAR": "y"},
			format: "%d",
			want:   "",
		},
		{
			m:      env.Map{"FOO": "x", "BAR": "y"},
			format: "%f",
			want:   "",
		},
	}
	for _, tt := range tests {
		got := fmt.Sprintf(tt.format, tt.m)
		if diff := cmp.Diff(got, tt.want); diff != "" {
			formatErrorf(t, tt.m, tt.format, got, tt.want, diff)
		}
	}
}

func formatErrorf(t *testing.T, m env.Map, format, got, want, diff string) {
	t.Helper()
	basicm := map[string]string(m)
	t.Errorf("formatting %#v with %v: got %q, want %q: %s", basicm, format, got, want, diff)
}

func TestEncode(t *testing.T) {
	tests := []struct {
		m    env.Map
		want []string
	}{
		{
			m:    env.Map{},
			want: []string{},
		},
		{
			m:    env.Map{"k": "v"},
			want: []string{"k=v"},
		},
		{
			m:    env.Map{"k": ""},
			want: []string{"k="},
		},
		{
			m:    env.Map{"FOO": "x", "BAR": "y"},
			want: []string{"BAR=y", "FOO=x"},
		},
		{
			m:    env.Map{"FOO": "x", "BAR": "", "BAZ": "z"},
			want: []string{"BAR=", "BAZ=z", "FOO=x"},
		},
	}
	for _, tt := range tests {
		got := tt.m.Encode()
		if diff := cmp.Diff(got, tt.want); diff != "" {
			encodeErrorf(t, tt.m, got, tt.want, diff)
		}
	}
}

func encodeErrorf(t *testing.T, m env.Map, got, want []string, diff string) {
	t.Helper()
	basicm := map[string]string(m)
	t.Errorf("%#v.Encode() = %v, want %v: %s", basicm, got, want, diff)
}

func TestParse(t *testing.T) {
	tests := []struct {
		kvs  []string
		want env.Map
	}{
		{
			kvs:  []string{},
			want: env.Map{},
		},
		{
			kvs:  []string{"FOO=x"},
			want: env.Map{"FOO": "x"},
		},
		{
			kvs:  []string{"FOO=x", "blah"},
			want: env.Map{"FOO": "x"},
		},
		{
			kvs:  []string{"FOO="},
			want: env.Map{"FOO": ""},
		},
		{
			kvs:  []string{"FOO=x", "BAR=y"},
			want: env.Map{"FOO": "x", "BAR": "y"},
		},
	}
	for _, tt := range tests {
		got := env.Parse(tt.kvs...)
		if diff := cmp.Diff(got, tt.want); diff != "" {
			parseErrorf(t, tt.kvs, got, tt.want, diff)
		}
	}
}

func parseErrorf(t *testing.T, kvs []string, got, want env.Map, diff string) {
	t.Helper()
	basicgot := map[string]string(got)
	basicwant := map[string]string(want)
	t.Errorf("Parse(%v) = %#v, want %#v: %s", kvs, basicgot, basicwant, diff)
}

func TestMerge(t *testing.T) {
	tests := []struct {
		maps []env.Map
		want env.Map
	}{
		{
			maps: nil,
			want: env.Map{},
		},
		{
			maps: []env.Map{},
			want: env.Map{},
		},
		{
			maps: []env.Map{
				{"FOO": "x"},
			},
			want: env.Map{"FOO": "x"},
		},
		{
			maps: []env.Map{
				{"FOO": "x"},
				{"BAR": "y"},
			},
			want: env.Map{
				"FOO": "x",
				"BAR": "y",
			},
		},
		{
			maps: []env.Map{
				{
					"FOO": "x",
				},
				{
					"BAR": "y",
				},
				{
					"BAR": "z",
				},
			},
			want: env.Map{
				"FOO": "x",
				"BAR": "z",
			},
		},
	}
	for _, tt := range tests {
		got := env.Merge(tt.maps...)
		if diff := cmp.Diff(got, tt.want); diff != "" {
			mergeErrorf(t, tt.maps, got, tt.want, diff)
		}
	}
}

func mergeErrorf(t *testing.T, maps []env.Map, got, want env.Map, diff string) {
	t.Helper()
	t.Errorf("Merge(%v) = %v, want %v: %s", maps, got, want, diff)
}

func TestVariables(t *testing.T) {
	got := env.Variables()
	want := env.Parse(os.Environ()...)
	if diff := cmp.Diff(got, want); diff != "" {
		variablesErrorf(t, got, want)
	}
}

func variablesErrorf(t *testing.T, got, want env.Map) {
	t.Helper()
	t.Errorf("Variables() = %v, want %v", got, want)
}

func TestDiff(t *testing.T) {
	tests := []struct {
		x, y env.Map
		want env.Diff
	}{
		{
			x: env.Map{"FOO": "x"},
			y: env.Map{"BAR": "y"},
			want: env.Diff{
				OnlyInM: env.Map{"FOO": "x"},
				Changes: nil,
				OnlyInN: env.Map{"BAR": "y"},
			},
		},
		{
			x: env.Map{"FOO": "x", "BAR": "a"},
			y: env.Map{"FOO": "x", "BAR": "b"},
			want: env.Diff{
				Changes: []env.Change{
					{
						Key:    "BAR",
						MValue: "a",
						NValue: "b",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		got := tt.x.Diff(tt.y)
		if diff := cmp.Diff(got, tt.want); diff != "" {
			diffErrorf(t, tt.x, tt.y, got, tt.want, diff)
		}
	}
}

func diffErrorf(t *testing.T, x, y env.Map, got, want env.Diff, diff string) {
	t.Helper()
	t.Errorf("Diff(%v, %v) = %#v, want %#v: %s", x, y, got, want, diff)
}

func TestChangeString(t *testing.T) {
	ch := env.Change{
		Key:    "FOO",
		MValue: "a",
		NValue: "b",
	}
	want := "FOO: a -> b"
	got := ch.String()
	if got != want {
		t.Errorf("%#v.String() = %q, want %q", ch, got, want)
	}
}
