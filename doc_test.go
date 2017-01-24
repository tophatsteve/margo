package margo_test

import (
	"fmt"
	"github.com/tophatsteve/margo"
)

// Generate a sentence that is no longer than 140 characters long, and starts with
// a capital letter. The set of Markov Chains is generated from the lines slice with
// a prefix size of 2.
// The text in the lines slice won't actually generate any variations because there
// are no matching prefixes, but it serves as an example.
func ExampleGenerateSentence() {
	lines := []string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		"Aenean nec mattis sapien, nec scelerisque mauris.",
		"Sed et tortor sit amet lectus laoreet finibus vitae et erat.",
		"Integer imperdiet sodales urna a vehicula.",
		"Donec aliquam finibus dignissim.",
	}

	chainPrefixSize := 2

	m := margo.NewMargo(lines, chainPrefixSize)
	s := m.GenerateSentence(140, true)

	fmt.Println(s)
}
