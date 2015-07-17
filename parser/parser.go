package parser

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var propertyRegexp = regexp.MustCompile(`\s?@property\s?(?:\((.*)\))?\s?([^\s\*]*)\s?(\*)?(.*);`)

type ObjCClassInfo struct {
	Properties []Property
}

type Property struct {
	Name, Class string
	Attributes  []string
	IsPointer   bool
}

// Right now we read the whole file in memory. Normally there should not be a problem
// as source files are not extremely big. Anyway this could be improved
func ParseFile(filename string) (*ObjCClassInfo, error) {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	info := ObjCClassInfo{}
	propertyMatches := propertyRegexp.FindAllSubmatch(fileBytes, -1)

	// Extract all properties
	for _, propertyMatch := range propertyMatches {
		//Split the attributes string and trim each of them
		attributes := strings.Split(string(propertyMatch[1]), ",")
		for i, attr := range attributes {
			attributes[i] = strings.TrimSpace(attr)
		}
		class := string(propertyMatch[2])
		pointer := string(propertyMatch[3])
		name := string(propertyMatch[4])

		// Add this property to the info
		info.Properties = append(info.Properties, Property{
			Name:       name,
			Class:      class,
			Attributes: attributes,
			IsPointer:  pointer != "",
		})
	}

	return &info, nil
}

func GetParseableFiles(rootPath string) []string {
	var headerFiles []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		isHeader, err := filepath.Match("*.h", info.Name())
		if err != nil {
			return err
		}

		if isHeader {
			headerFiles = append(headerFiles, path)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return headerFiles
}
