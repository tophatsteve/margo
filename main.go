package main

import (
	"log"
	"github.com/tophatsteve/margo/markov"
)

func main() {
	log.Printf("%s", markov.Generate("./data/test.txt", 2, 140))	
}