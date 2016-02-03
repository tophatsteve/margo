package margo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildChainSetSize(t *testing.T) {

	lines := []string {"And the Golden Grouse And the Pobble who"}
	chainSet := BuildChainSet(lines, 2)

	assert.Equal(t, 5, len(chainSet.Chains), "Number of chains should be 5")
}

func TestBuildChainSetChainCount(t *testing.T) {

	lines := []string {"And the Golden Grouse And the Pobble who"}
	chainSet := BuildChainSet(lines, 2)

	assert.Equal(t, 2, len(chainSet.Chains["And the"]), "Number of chains with key 'And the' should be 2")
}

func TestLookupChainWithMatchingKey(t *testing.T) {

	lines := []string {"And the Golden Grouse And the Pobble who"}
	chainSet := BuildChainSet(lines, 2)

	assert.Equal(t, 2, len(chainSet.lookupChains("And the")), "Find chains with matching key")
}

func TestLookupChainWithoutMatchingKey(t *testing.T) {

	lines := []string {"And the Golden Grouse And the Pobble who"}
	chainSet := BuildChainSet(lines, 2)

	assert.Equal(t, 0, len(chainSet.lookupChains("Grouse who")), "Do not find chains when no matching key")
}

func TestBuildLookupKey(t *testing.T) {
	chain := Chain{}
	chain.Prefix = append(chain.Prefix, "And")
	chain.Prefix = append(chain.Prefix, "the")
	chain.Suffix = "Golden"

	assert.Equal(t, "the Golden", chain.buildNextLookupKey(), "Next key should be all all but first prefix word plus suffix word")	
}

func TestChainToStringPrefix(t *testing.T) {

	chain := Chain{}
	chain.Prefix = append(chain.Prefix, "And")
	chain.Prefix = append(chain.Prefix, "the")
	chain.Suffix = "Golden"

	assert.Equal(t, "And the", chain.ToStringPrefix(), "Prefix should be all prefix words joined by a space")
}

func TestChainToString(t *testing.T) {

	chain := Chain{}
	chain.Prefix = append(chain.Prefix, "And")
	chain.Prefix = append(chain.Prefix, "the")
	chain.Suffix = "Golden"

	assert.Equal(t, "And the Golden", chain.ToString(), "ToString should return all prefix words and suffix word joined by a space")
}
