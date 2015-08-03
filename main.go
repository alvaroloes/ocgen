package main

import (
	"fmt"

	"github.com/alvaroloes/ocgen/generator"
	"github.com/alvaroloes/ocgen/parser"
)

func main() {
	// Get all the header files under the directory
	fileNames := parser.GetParseableFiles("./fixtures")
	fmt.Println(fileNames)

	classFile, _ := parser.Parse(fileNames[0])
	generator.GenerateMethods(classFile)

	// Parse each file
	// for _, fileName := range fileNames {
	// 	parser.ParseAndGetClassesInfo(fileName, "asdf")
	// }

	//TODO: In the future, this will have some flags
	// flag.Parse()
	// filenames := flag.Args()

	// if len(filenames) == 0 {
	// 	log.Fatal(errors.New("No input files have been specified"))
	// }

	// for _, filename := range filenames {
	// 	classInfo, err := parser.ParseFile(filename)
	// 	if err != nil {
	// 		log.Printf("Error processing file %v: %v", filename, err)
	// 	}

	// 	nsCopyingCode, err := generator.NSCopying(classInfo)
	// 	if err != nil {
	// 		log.Printf("Error generating NSCopying code: %v", err)
	// 	}
	// 	fmt.Println(nsCopyingCode)

	// 	nsCodingCode, err := generator.NSCoding(classInfo)
	// 	if err != nil {
	// 		log.Printf("Error generating NSCoding code: %v", err)
	// 	}
	// 	fmt.Println(nsCodingCode)
	// }
}
