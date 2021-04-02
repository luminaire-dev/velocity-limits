package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/luminaire-dev/velocity-limits/processor"
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

	pro := processor.NewProcessor()

	for scanner.Scan() {
		// parse the JSON-encoded data
		var incomingLoad processor.Load
		if err := json.Unmarshal(scanner.Bytes(), &incomingLoad); err != nil {
			panic(err)
		}
		res := pro.Load(incomingLoad)

		// ignore duplicate IDs
		if res == nil {
			continue
		}

		jsonText, err := json.Marshal(*res)
		if err != nil {
			panic(err)
		}

		/// write to output file
		_, err = fmt.Fprintln(dest, string(jsonText))
		if err != nil {
			panic(err)
		}
	}

	// get file path and print location of output file
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Println("Done! 🎉")
	fmt.Println("Output file has been generated at " + path + "/generated_output")
}
