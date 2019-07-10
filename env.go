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

// Package env provides conveniences for working with environment variables,
// particularly in the context of executing external commands.
package env

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Map is a convenient representation of a set of environment variables.
type Map map[string]string

// String encodes the Map as space-separated "key=value" pairs, sorted
// lexicographically by key.
func (m Map) String() string {
	sb := new(strings.Builder)
	m.print(sb, ' ')
	return sb.String()
}

// Format implements fmt.Formatter for Map as follows:
//
// If the verb is anything but 'v', Format produces no output.
//
// If the '+' flag is specified, Format emits newline separated "key=value"
// pairs. Otherwise, it emits space-separated "key=value" pairs.
//
// Values are sorted lexicographically by key.
func (m Map) Format(s fmt.State, verb rune) {
	if verb != 'v' {
		return
	}
	if s.Flag('+') {
		m.print(s, '\n')
	} else {
		m.print(s, ' ')
	}
}

func (m Map) print(w io.Writer, sep rune) {
	i := 0
	for _, k := range m.keys() {
		fmt.Fprintf(w, "%s=%s", k, m[k])
		if i < len(m)-1 {
			fmt.Fprintf(w, "%c", sep)
		}
		i++
	}
}

func (m Map) keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Encode encodes the Map as a slice of "key=value" pairs, suitable for use
// with the os/exec package.
func (m Map) Encode() []string {
	kvs := make([]string, 0, len(m))
	for _, k := range m.keys() {
		kv := fmt.Sprintf("%s=%s", k, m[k])
		kvs = append(kvs, kv)
	}
	return kvs
}

// Diff computes differences between m and n.
func (m Map) Diff(n Map) Diff {
	d := Diff{}
	for k, mval := range m {
		nval, ok := n[k]
		switch {
		case !ok:
			if d.OnlyInM == nil {
				d.OnlyInM = make(Map)
			}
			d.OnlyInM[k] = mval
		case mval != nval:
			d.Changes = append(d.Changes, Change{
				Key:    k,
				MValue: mval,
				NValue: nval,
			})
		}
	}
	for k, nval := range n {
		_, ok := m[k]
		if !ok {
			if d.OnlyInN == nil {
				d.OnlyInN = make(Map)
			}
			d.OnlyInN[k] = nval
		}
	}
	return d
}

// Diff describes differences between two environments, "M" and "N".
type Diff struct {
	OnlyInM Map
	Changes []Change
	OnlyInN Map
}

// Change describes a change in a value in the environment.
type Change struct {
	Key    string
	MValue string
	NValue string
}

func (c Change) String() string {
	return fmt.Sprintf("%s: %s -> %s", c.Key, c.MValue, c.NValue)
}

// Variables returns a Map of the process environment.
func Variables() Map {
	return Parse(os.Environ()...)
}

// Parse parses a list of environment variables in "key=value" format.
// Values not in "key=value" format are ignored.
func Parse(kvs ...string) Map {
	m := make(Map)
	for _, kv := range kvs {
		i := strings.IndexRune(kv, '=')
		if i == -1 {
			continue
		}
		k := kv[:i]
		v := kv[i+1:]
		m[k] = v
	}
	return m
}

// Merge merges environment variable maps. In case of key collisions, values
// which appear later in the maps list take precedence.
func Merge(maps ...Map) Map {
	merged := make(Map)
	for _, m := range maps {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}
