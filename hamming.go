// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Package hamming computes the edits that convert one sequence
// to another.
//
// It is similar to the hamming  edit sequence but despite the name it
// is not the same.  In particular, it uses move operations as well as
// delete/insert and is a simple rather than optimal version.
package hamming

// Item represents an item with a key to use for equality
// comparisons. The Key value itself should be "comparable".
type Item interface {
	Key() interface{}
}

// Edits finds a sequence of spice/move operations that will convert
// the input seqeunce of items to the output.
//
// There is no requirement that the input items have distinct keys: an
// unique internal key is constructed using a counter.  So, the nth
// occurence of a key in the input is considered  the same as the nth
// occurence of the same key in the output.
//
// The Splice operation is used for deletes as well as inserts. The
// offset specifies the index in the current version where the set of
// "before" items  are replaced with the values from the after.
//
// See the test for example splice and move operations for arrays.
func Edits(in, out []Item, splice func(offset int, before, after []Item), move func(offset, count, distance int)) {
	d := differ{in, out, splice, move}
	d.Diff()
}

type differ struct {
	Input  []Item
	Output []Item
	Splice func(offset int, before, after []Item)
	Move   func(offset, count, distance int)
}

func (d differ) Diff() {
	// counter is used to make keys unique
	counter := map[interface{}]int{}

	// indices tracks unique key to index in input
	indices := map[interface{}]int{}

	// this loop sets input to [0..n) + computes indices
	input := make([]int, len(d.Input))
	for kk, item := range d.Input {
		key := item.Key()
		count := counter[key]
		counter[key]++

		key = [2]interface{}{key, count} // unique!
		indices[key] = kk
		input[kk] = kk
	}

	output := make([]int, len(d.Output))
	counter = map[interface{}]int{}
	for kk, item := range d.Output {
		key := item.Key()
		count := counter[key]
		counter[key]++

		key = [2]interface{}{key, count} // unique!

		// if key matches element in input, use input item index.
		if idx, ok := indices[key]; ok {
			output[kk] = idx
			delete(indices, key)
		} else {
			// use -(kk+1) as index so it guarantees neg number
			// this negative check is used to identify items in
			// the output  that are not there in the input
			output[kk] = -(kk + 1)
		}
	}

	// all items left over in "indices" are meant  to be deleted
	removed := map[int]bool{}
	for _, v := range indices {
		removed[v] = true
	}
	d.diff(input, output, removed)
}

func (d differ) diff(input, output []int, removed map[int]bool) {
	empty := []Item(nil)
	i, o := len(input), len(output)

	for i > 0 && o > 0 { // iterate end to start
		switch {
		case input[i-1] == output[o-1]: // same
			i--
			o--
		case removed[input[i-1]]: // deleted
			d.Splice(i-1, []Item{d.Input[input[i-1]]}, empty)
			i--
		case output[o-1] < 0: // inserted
			item := d.Output[-1-output[o-1]]
			d.Splice(i, empty, []Item{item})
			o--
		default: // moved  from before to this position
			index := -1
			for kk := 0; kk < i; kk++ {
				if input[kk] == output[o-1] {
					index = kk
					break
				}
			}
			d.Move(index, 1, i-1-index)
			input = append(input[:index], input[index+1:i]...)
			i--
			o--
		}
	}

	if o > 0 { // left over output is just inserted
		items := make([]Item, o)
		for kk := 0; kk < o; kk++ {
			items[kk] = d.Output[-1-output[kk]]
		}
		d.Splice(0, empty, items)
	}

	if i > 0 { // left over input is just deleted
		items := make([]Item, i)
		for kk := 0; kk < i; kk++ {
			items[kk] = d.Input[input[kk]]
		}
		d.Splice(0, items, empty)
	}
}
