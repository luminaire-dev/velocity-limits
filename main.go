package main

import (
  "fmt"
  "bufio"
  "os"
	"github.com/luminaire-dev/velocity-limits/processor"
  "log"
)

func main() {
  // open input file
  file, err := os.Open("./input.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()

  // create output file
  dest, err := os.Create("./generated_output.txt")
  	if err != nil {
  		log.Fatal(err)
  	}
  	defer dest.Close()

    // read from input file
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
      res:= processor.Process(scanner.Bytes())

      // ignore duplicate IDs
  		if res == nil {
  			continue
  		}

      /// write to output file
      _, err = fmt.Fprintln(dest, string(res))
      if err != nil {
        log.Fatal(err)
      }
    }

    // get file path and print location of output file
    path, err := os.Getwd()
    if err != nil {
        log.Println(err)
    }

		fmt.Println("Done! 🎉")
    fmt.Println("Output file has been generated at " + path + "/generated_output")
}
