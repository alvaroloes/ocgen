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

var (
	interfaceRegexp      = regexp.MustCompile(`(?ms:^\s?@interface\s+([^:<\s]*).*?` + ocgenMarker + `.*?@end)`)
	implementationRegexp = regexp.MustCompile(`(?ms:^\s?@implementation\s+([^\s]*).*?@end)`)
)

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

	implFileName := implFileNameFromHeader(headerFileName)
	implFileBytes, err := ioutil.ReadFile(implFileName)
	if err != nil {
		log.Printf("Unable to open implementation file: %v\n", err)
		return nil, err
	}

	classesInfo := getClasses(headerFileBytes, implFileBytes, implFileName)

	return classesInfo, nil
}

func implFileNameFromHeader(headerFileName string) string {
	return headerFileName[:len(headerFileName)-len(headerFileExt)] + ".m"
}

func getClasses(headerFileBytes, implFileBytes []byte, implFileName string) []ObjCClass {
	matchedHInterfaces := interfaceRegexp.FindAllSubmatchIndex(headerFileBytes, -1)
	if matchedHInterfaces == nil {
		return []ObjCClass{} // No interfaces in header file
	}

	// TODO: No need for this. Use dynamic regexp here
	matchedImplementations := implementationRegexp.FindAllSubmatchIndex(implFileBytes, -1)
	if matchedImplementations == nil {
		return []ObjCClass{} // No implementations? This would be so weird
	}

	// TODO: No need for this. Use dynamic regexp here
	//matchedMInterfaces := interfaceRegexp.FindAllSubmatchIndex(implFileBytes, -1)

	classesInfo := make([]ObjCClass, len(matchedHInterfaces))

	for i, matchedInterface := range matchedHInterfaces {
		className := string(headerFileBytes[matchedInterface[2]:matchedInterface[3]])
		interfaceHBytes := headerFileBytes[matchedInterface[0]:matchedInterface[1]]
		interfaceMBytes := []byte{} //TODO
		implBytes := implBytesForClassName(className, matchedImplementations, implFileBytes)

		classesInfo[i] = NewObjCClass(interfaceHBytes, interfaceMBytes, implBytes, implFileName)
	}
	return classesInfo
}

func implBytesForClassName(className string, matchedImplementations [][]int, implFileBytes []byte) []byte {
	for _, matchedImpl := range matchedImplementations {
		implName := string(implFileBytes[matchedImpl[2]:matchedImpl[3]])
		if implName == className {
			return implFileBytes[matchedImpl[0]:matchedImpl[1]]
		}
	}
	return nil
}
