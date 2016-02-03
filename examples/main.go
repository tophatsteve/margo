package main

import (
	"os"
	"bufio"    
	"log"
    flag "launchpad.net/gnuflag"
	"github.com/tophatsteve/margo"
)

var filename string
var prefixLength int

func init() {
    flag.StringVar(&filename, "filename", "", "The file containing lines of sample text")
    flag.StringVar(&filename, "f", "", "The file containing lines of sample text")
    flag.IntVar(&prefixLength, "prefix", 2, "The chain prefi length")
    flag.IntVar(&prefixLength, "p", 2, "The chain prefix length")
}

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