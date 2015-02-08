package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/alvaroloes/ocgen/parser"
	"github.com/alvaroloes/ocgen/generator")

func main() {
	//TODO: In the future, this will have some flags
	flag.Parse()
	filenames := flag.Args()

	if len(filenames) == 0 {
		log.Fatal(errors.New("No input files have been specified"))
	}

	for _, filename := range filenames {
		classInfo, err := parser.ParseFile(filename)
		if err != nil {
			log.Printf("Error processing file %v: %v", filename, err)
		}

		nsCopyingCode, err := generator.NSCopying(classInfo)
		if (err != nil) {
			log.Printf("Error generating NSCopying code: %v", err)
		}
		fmt.Println(nsCopyingCode)

		nsCodingCode, err := generator.NSCoding(classInfo)
		if (err != nil) {
			log.Printf("Error generating NSCoding code: %v", err)
		}
		fmt.Println(nsCodingCode)
	}
}
