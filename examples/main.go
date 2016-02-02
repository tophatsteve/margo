package main

import (
	"os"
	"bufio"    
	"log"
	"github.com/tophatsteve/margo/margo"
)

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
		// convert to lowecase
		// remove punctuation
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
    lines := loadLines("./data/test.txt")
	log.Printf("%s", margo.GenerateSentence(lines, 2, 140))	
}