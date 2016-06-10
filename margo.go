/*
	Package margo is a Markov Chain generator.

	Example code:

        package main

        import (
            "bufio"
            "log"
            "os"

            "github.com/tophatsteve/margo"
            flag "launchpad.net/gnuflag"
        )

        var filename string
        var prefixLength int

        func init() {
            flag.StringVar(&filename, "filename", "", "The file containing lines of sample text")
            flag.StringVar(&filename, "f", "", "The file containing lines of sample text")
            flag.IntVar(&prefixLength, "prefix", 2, "The chain prefi length")
            flag.IntVar(&prefixLength, "p", 2, "The chain prefix length")
        }

        // load lines from a file into a []string
        func loadLines(filename string) []string {
            lines := []string{}

            // open a file
            if file, err := os.Open(filename); err == nil {

                // make sure it gets closed
                defer file.Close()

                // create a new scanner and read the file line by line
                scanner := bufio.NewScanner(file)
                for scanner.Scan() {
                    lines = append(lines, scanner.Text())
                }

                // check for errors
                if err = scanner.Err(); err != nil {
                    log.Fatal(err)
                }

            } else {
                log.Fatal(err)
            }

            return lines
        }

        func main() {
            flag.Parse(true)
            lines := loadLines(filename)
            log.Printf("%s", margo.GenerateSentence(lines, prefixLength, 140))
        }


*/
package margo

import (
	"fmt"
	"log"
	"strings"
)

// Chain is a set of prefixes followed by a suffix.
type chain struct {
	prefix []string
	suffix string
}

// ChainSet is a collection of Chains which also defines how long each prefix should be.
type chainSet struct {
	name         string
	prefixLength int
	chains       map[string][]chain
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

func (cs chainSet) lookupChains(prefix string) []chain {
	if val, ok := cs.chains[prefix]; ok {
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

func (cs chainSet) pickFirstChain() chain {
	keys := make([]string, 0, len(cs.chains))
	for k := range cs.chains {
		keys = append(keys, k)
	}

	firstChainValue := cs.chains[keys[randomNumber(len(keys))]]
	return firstChainValue[randomNumber(len(firstChainValue))]
}

func (cs chainSet) pickNextChain(c chain) chain {
	chains := cs.lookupChains(c.buildNextLookupKey())

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

func buildChainSet(lines []string, prefixSize int) chainSet {
	chainSet := chainSet{}
	chainSet.prefixLength = prefixSize
	chainSet.chains = make(map[string][]chain)
	msg := make(chan []chain)

	defer close(msg)

	for _, v := range lines {
		go buildChainsFromLine(msg, v, prefixSize)
	}

	for x := 0; x < len(lines); x++ {
		chains := <-msg
		for _, v := range chains {
			if _, ok := chainSet.chains[v.toStringPrefix()]; !ok {
				chainSet.chains[v.toStringPrefix()] = []chain{}
			}
			chainSet.chains[v.toStringPrefix()] = append(chainSet.chains[v.toStringPrefix()], v)
		}
	}
	return chainSet
}

func buildSentence(chainset chainSet, maxLength int) string {
	c1 := chainset.pickFirstChain()
	result := c1.toString()
	for len(c1.suffix) > 0 {
		c1 = chainset.pickNextChain(c1)

		if len(result)+len(c1.suffix) > maxLength {
			break
		}

		result = fmt.Sprint(result, " ", c1.suffix)
	}

	return result
}

// GenerateSentence is used to generate a sentence from a set of strings, upto a maximum length.
// prefixSize specifies how long the prefix should be in each chain.
func GenerateSentence(lines []string, prefixSize int, maxLength int) string {
	chainset := buildChainSet(lines, prefixSize)
	return buildSentence(chainset, maxLength)
}
