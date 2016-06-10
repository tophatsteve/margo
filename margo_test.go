package margo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildChainSetSize(t *testing.T) {

	lines := []string{"And the Golden Grouse And the Pobble who"}
	chainSet := buildChainSet(lines, 2)

	assert.Equal(t, 5, len(chainSet.chains), "Number of chains should be 5")
}

func TestBuildChainSetChainCount(t *testing.T) {

	lines := []string{"And the Golden Grouse And the Pobble who"}
	chainSet := buildChainSet(lines, 2)

	assert.Equal(t, 2, len(chainSet.chains["And the"]), "Number of chains with key 'And the' should be 2")
}

func TestLookupChainWithMatchingKey(t *testing.T) {

	lines := []string{"And the Golden Grouse And the Pobble who"}
	chainSet := buildChainSet(lines, 2)

	assert.Equal(t, 2, len(chainSet.lookupChains("And the")), "Find chains with matching key")
}

func TestLookupChainWithoutMatchingKey(t *testing.T) {

	lines := []string{"And the Golden Grouse And the Pobble who"}
	chainSet := buildChainSet(lines, 2)

	assert.Equal(t, 0, len(chainSet.lookupChains("Grouse who")), "Do not find chains when no matching key")
}

func TestBuildLookupKey(t *testing.T) {
	chain := chain{}
	chain.prefix = append(chain.prefix, "And")
	chain.prefix = append(chain.prefix, "the")
	chain.suffix = "Golden"

	assert.Equal(t, "the Golden", chain.buildNextLookupKey(), "Next key should be all all but first prefix word plus suffix word")
}

func TestChainToStringPrefix(t *testing.T) {

	chain := chain{}
	chain.prefix = append(chain.prefix, "And")
	chain.prefix = append(chain.prefix, "the")
	chain.suffix = "Golden"

	assert.Equal(t, "And the", chain.toStringPrefix(), "Prefix should be all prefix words joined by a space")
}

func TestChainToString(t *testing.T) {

	chain := chain{}
	chain.prefix = append(chain.prefix, "And")
	chain.prefix = append(chain.prefix, "the")
	chain.suffix = "Golden"

	assert.Equal(t, "And the Golden", chain.toString(), "ToString should return all prefix words and suffix word joined by a space")
}
