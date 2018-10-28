// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package hamming_test

import (
	"github.com/tvastar/hamming"
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

type keyer string

func (k keyer) Key() interface{} {
	return string(k)
}

func makeItems(args ...string) []hamming.Item {
	result := make([]hamming.Item, len(args))
	for idx, arg := range args {
		result[idx] = keyer(arg)
	}
	return result
}

func randomInputOutput() [2][]string {
	rand.Seed(42)
	removed := strings.Split("abcdefghijk", "")
	inserted := strings.Split("ABCDEFGHIJK", "")
	common := strings.Split("0123456789", "")
	inputx := append(append([]string(nil), removed...), common...)
	outputx := append(append([]string(nil), inserted...), common...)

	input := make([]string, len(inputx))
	output := make([]string, len(outputx))
	for kk, idx := range rand.Perm(len(input)) {
		input[kk] = inputx[idx]
	}

	for kk, idx := range rand.Perm(len(output)) {
		output[kk] = outputx[idx]
	}

	return [2][]string{input, output}
}

func TestCases(t *testing.T) {
	cases := map[string][2][]string{
		"DeleteAll":    {{"a", "b", "c"}, {}},
		"InsertAll":    {{}, {"a", "b", "c"}},
		"Replace":      {{"a", "b", "c"}, {"d"}},
		"DeleteSome":   {{"a", "b", "c"}, {"a"}},
		"DeleteMiddle": {{"a", "b", "c"}, {"a", "c"}},
		"Shuffle":      {{"a", "b", "c"}, {"c", "b", "a"}},
		"InsertSome":   {{"a", "b", "c"}, {"a", "d", "b", "c"}},
		"Complex":      {{"a", "b", "c"}, {"d", "c", "b"}},

		"Random": randomInputOutput(),
	}

	for testName, v := range cases {
		t.Run(testName, func(t *testing.T) {
			input, output := v[0], v[1]
			in := makeItems(input...)
			out := makeItems(output...)
			hamming.Edits(
				makeItems(input...),
				makeItems(output...),
				func(offset int, before, after []hamming.Item) {
					rest := append([]hamming.Item(nil), in[offset+len(before):]...)
					in = append(append(in[:offset], after...), rest...)
				},
				func(offset, count, distance int) {
					o1 := in[:offset:offset]
					o2 := in[offset : offset+count]
					o3 := in[offset+count : offset+count+distance]
					o4 := in[offset+count+distance:]
					x := append(append(o1, o3...), o2...)
					in = append(x, o4...)
				},
			)
			if !reflect.DeepEqual(in, out) {
				t.Error("Mismatched", input, output, in)
			}
		})
	}
}
