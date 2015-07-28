package parser

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

const (
	ocgenMarker   = "OCGEN_AUTO"
	headerFileExt = ".h"
)

var interfaceRegexp = regexp.MustCompile(`(?ms:^\s?@interface.*?` + ocgenMarker + `.*?@end)`)

func GetParseableFiles(rootPath string) []string {
	var headerFiles []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		isHeader, err := filepath.Match("*"+headerFileExt, info.Name())
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

func ParseAndGetClassesInfo(headerFileName string) ([]ObjCClass, error) {
	headerFileBytes, err := ioutil.ReadFile(headerFileName)
	if err != nil {
		log.Printf("Unable to open header file %v\n", err)
		return nil, err
	}

	implFile, err := os.Open(implFileNameFromHeader(headerFileName))
	if err != nil {
		log.Printf("Unable to open implementation file: %v\n", err)
		return nil, err
	}

	classesInfo := getClasses(headerFileBytes, implFile)

	return classesInfo, nil
}

func implFileNameFromHeader(headerFileName string) string {
	return headerFileName[:len(headerFileName)-len(headerFileExt)] + ".m"
}

func getClasses(headerFileBytes []byte, implFile *os.File) []ObjCClass {
	matchedInterfaces := interfaceRegexp.FindAllIndex(headerFileBytes, -1)

	if matchedInterfaces == nil {
		return []ObjCClass{} // No classes in this file
	}

	classesInfo := make([]ObjCClass, len(matchedInterfaces))

	for i, matchedInterface := range matchedInterfaces {
		start := matchedInterface[0]
		end := matchedInterface[1]

		classesInfo[i] = NewObjCClass(headerFileBytes[start:end], implFile)
	}
	return classesInfo
}
