/*Package margo is a Markov Chain generator.

 */
package margo

import (
	"fmt"
	"log"
	"strings"
)

// chain is a set of prefixes followed by a suffix.
type chain struct {
	prefix []string
	suffix string
}

// Margo is a collection of chains which also defines how long each prefix should be.
type Margo struct {
	name         string
	prefixLength int
	chains       map[string][]chain
}

// NewMargo creates a new Margo instance from the passed in lines and builds
// a set of Markov Chains with a prefix length of prefixSize.
func NewMargo(lines []string, prefixSize int) Margo {
	m := Margo{}
	m.prefixLength = prefixSize
	m.chains = make(map[string][]chain)
	msg := make(chan []chain)

	defer close(msg)

	for _, v := range lines {
		go buildChainsFromLine(msg, v, prefixSize)
	}

	for x := 0; x < len(lines); x++ {
		chains := <-msg
		for _, v := range chains {
			if _, ok := m.chains[v.toStringPrefix()]; !ok {
				m.chains[v.toStringPrefix()] = []chain{}
			}
			m.chains[v.toStringPrefix()] = append(m.chains[v.toStringPrefix()], v)
		}
	}
	return m
}

// GenerateSentence generates a random sentence using the sentences in lines to build a set of Markov Chains.
// The sentence is generated based on the size of the chain prefix, prefixSize, the maximum length of
// the returned sentence, maxLength, and whether the sentence should start with a capital letter, capitalStart.
func GenerateSentence(lines []string, prefixSize, maxLength int, capitalStart bool) string {
	m := NewMargo(lines, prefixSize)
	return m.buildSentence(maxLength, capitalStart)
}

func (c chain) toStringPrefix() string {
	return strings.Join(c.prefix, " ")
}

func (c chain) toString() string {
	return fmt.Sprint(c.toStringPrefix(), " ", c.suffix)
}

func (c chain) buildNextLookupKey() string {
	words := make([]string, len(c.prefix))
	copy(words, c.prefix)
	if len(words) > 0 {
		words = words[1:len(c.prefix)]
	}
	words = append(words, c.suffix)

	return strings.Join(words, " ")
}

func (m Margo) lookupChains(prefix string) []chain {
	if val, ok := m.chains[prefix]; ok {
		return val
	}

	return []chain{}
}

func buildChainsFromLine(msg chan []chain, line string, prefixSize int) {
	chains := []chain{}
	words := strings.Split(line, " ")

	for i := 0; i < len(words)-prefixSize; i++ {
		c := chain{}
		for p := 0; p < prefixSize; p++ {
			c.prefix = append(c.prefix, words[i+p])
		}
		c.suffix = words[i+prefixSize]
		chains = append(chains, c)
	}

	msg <- chains
}

func (m Margo) pickFirstChain(capitalStart bool) chain {
	keys := make([]string, 0, len(m.chains))
	for k := range m.chains {
		keys = append(keys, k)
	}

	for {
		firstChainValue := m.chains[keys[randomNumber(len(keys))]]
		randomChain := firstChainValue[randomNumber(len(firstChainValue))]
		if capitalStart == false {
			return randomChain
		}

		if isUppercase(randomChain.prefix[0], 0) {
			return randomChain
		}
	}
}

func (m Margo) pickNextChain(c chain) chain {
	chains := m.lookupChains(c.buildNextLookupKey())

	if len(chains) == 0 {
		return chain{}
	}

	return chains[randomNumber(len(chains))]
}

func dumpChains(chains []chain) {
	for _, v := range chains {
		log.Printf("%s", v.toString())
	}
}

func (m Margo) buildSentence(maxLength int, capitalStart bool) string {
	c1 := m.pickFirstChain(capitalStart)
	result := c1.toString()
	for len(c1.suffix) > 0 {
		c1 = m.pickNextChain(c1)

		if string(result[len(result)-1]) == "." || len(result)+len(c1.suffix) > maxLength {
			break
		}

		result = fmt.Sprint(result, " ", c1.suffix)
	}

	return result
}
