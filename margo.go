/*
Package margo is
*/
package margo

import (
	"fmt"
	"log"
	"strings"
)

// Chain is
type Chain struct {
	Prefix []string
	Suffix string
}

// ChainSet is
type ChainSet struct {
	Name         string
	PrefixLength int
	Chains       map[string][]Chain
}

// ToStringPrefix is
func (c Chain) ToStringPrefix() string {
	return strings.Join(c.Prefix, " ")
}

// ToString is
func (c Chain) ToString() string {
	return fmt.Sprint(c.ToStringPrefix(), " ", c.Suffix)
}

func (c Chain) buildNextLookupKey() string {
	words := make([]string, len(c.Prefix))
	copy(words, c.Prefix)
	if len(words) > 0 {
		words = words[1:len(c.Prefix)]
	}
	words = append(words, c.Suffix)

	return strings.Join(words, " ")
}

func (cs ChainSet) lookupChains(prefix string) []Chain {
	if val, ok := cs.Chains[prefix]; ok {
		return val
	}

	return []Chain{}
}

func buildChainsFromLine(msg chan []Chain, line string, prefixSize int) {
	chains := []Chain{}
	words := strings.Split(line, " ")

	for i := 0; i < len(words)-prefixSize; i++ {
		c := Chain{}
		for p := 0; p < prefixSize; p++ {
			c.Prefix = append(c.Prefix, words[i+p])
		}
		c.Suffix = words[i+prefixSize]
		chains = append(chains, c)
	}

	msg <- chains
}

func (cs ChainSet) pickFirstChain() Chain {
	keys := make([]string, 0, len(cs.Chains))
	for k := range cs.Chains {
		keys = append(keys, k)
	}

	firstChainValue := cs.Chains[keys[randomNumber(len(keys))]]
	return firstChainValue[randomNumber(len(firstChainValue))]
}

func (cs ChainSet) pickNextChain(c Chain) Chain {
	chains := cs.lookupChains(c.buildNextLookupKey())

	if len(chains) == 0 {
		return Chain{}
	}

	return chains[randomNumber(len(chains))]
}

func dumpChains(chains []Chain) {
	for _, v := range chains {
		log.Printf("%s", v.ToString())
	}
}

// BuildChainSet is
func BuildChainSet(lines []string, prefixSize int) ChainSet {
	chainSet := ChainSet{}
	chainSet.PrefixLength = prefixSize
	chainSet.Chains = make(map[string][]Chain)
	msg := make(chan []Chain)

	defer close(msg)

	for _, v := range lines {
		go buildChainsFromLine(msg, v, prefixSize)
	}

	for x := 0; x < len(lines); x++ {
		chains := <-msg
		for _, v := range chains {
			if _, ok := chainSet.Chains[v.ToStringPrefix()]; !ok {
				chainSet.Chains[v.ToStringPrefix()] = []Chain{}
			}
			chainSet.Chains[v.ToStringPrefix()] = append(chainSet.Chains[v.ToStringPrefix()], v)
		}
	}
	return chainSet
}

// BuildSentence is
func BuildSentence(chainset ChainSet, maxLength int) string {
	c1 := chainset.pickFirstChain()
	result := c1.ToString()
	for len(c1.Suffix) > 0 {
		c1 = chainset.pickNextChain(c1)
		result = fmt.Sprint(result, " ", c1.Suffix)
	}

	return result    
}

// GenerateSentence is
func GenerateSentence(lines []string, prefixSize int, maxLength int) string {
	chainset := BuildChainSet(lines, prefixSize)
    return BuildSentence(chainset, maxLength)
}