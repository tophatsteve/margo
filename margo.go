//Package margo is a Markov Chain generator.
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

// GenerateSentence generates a random sentence no bigger than maxLength. If maxLength is 0,
// then there is no limit on the length of sentence. If capitalStart is set to true then
// the sentence will begin with a word that starts with a capital letter in the set of lines used
// to when this Margo instance was created by a call to NewMargo.
func (m Margo) GenerateSentence(maxLength int, capitalStart bool) string {
	return m.buildSentence(maxLength, capitalStart)
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

func (m Margo) buildSentence(maxLength int, capitalStart bool) string {
	c1 := m.pickFirstChain(capitalStart)
	result := c1.toString()
	for len(c1.suffix) > 0 {
		c1 = m.pickNextChain(c1)

		// if the last character is a period, the sentence is complete
		if string(result[len(result)-1]) == "." {
			break
		}

		// if we have set a max length, and adding the latest suffix will make
		// the sentence bigger than the max length, the sentence is complete
		if maxLength > 0 && len(result)+len(c1.suffix) > maxLength {
			break
		}

		result = fmt.Sprint(result, " ", c1.suffix)
	}

	return result
}

func (m Margo) lookupChains(prefix string) []chain {
	if val, ok := m.chains[prefix]; ok {
		return val
	}

	return []chain{}
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

func dumpChains(chains []chain) {
	for _, v := range chains {
		log.Printf("%s", v.toString())
	}
}
