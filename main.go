package main

import (
  "fmt"
  "bufio"
  "os"
	"github.com/luminaire-dev/velocity-limits/processor"
)

func main() {
  file, err := os.Open("./input.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
      res:= processor.Process(scanner.Bytes(), )
      fmt.Println(res)
    }
    //TODO: generate output file
		fmt.Println("Done!")
}
