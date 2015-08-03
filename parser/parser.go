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

var interfaceRegexp = regexp.MustCompile(`(?ms:^\s?@interface\s+([^:<\s]*).*?` + ocgenMarker + `.*?@end)`)

const interfaceRegexpNameIndex = 1

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

func Parse(headerFileName string) (*ObjCClassFile, error) {
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

	classFile := ObjCClassFile{
		HName:   headerFileName,
		MName:   implFileName,
		Classes: getClasses(headerFileBytes, implFileBytes),
	}

	return &classFile, nil
}

func implFileNameFromHeader(headerFileName string) string {
	return headerFileName[:len(headerFileName)-len(headerFileExt)] + ".m"
}

func getClasses(headerFileBytes, implFileBytes []byte) []ObjCClass {
	// Search for all the interfaces in the header file
	matchedHInterfaces := interfaceRegexp.FindAllSubmatch(headerFileBytes, -1)
	if matchedHInterfaces == nil {
		return []ObjCClass{} // No interfaces in header file
	}

	classesInfo := make([]ObjCClass, len(matchedHInterfaces))
	for i, matchedInterface := range matchedHInterfaces {
		// Get the whole @interface bytes from header file
		interfaceHBytes := matchedInterface[0]
		// Get the class name to create the regexp for searching in the implementation file
		className := string(matchedInterface[interfaceRegexpNameIndex])

		// Get the whole @interface bytes from the implementation file
		classInterfaceRegexp := regexp.MustCompile(`(?ms:^\s?@interface\s+` + className + `\s+.*?@end)`)
		interfaceMBytes := classInterfaceRegexp.Find(implFileBytes)

		// Get the whole @implementation from the implementation file
		implRegexp := regexp.MustCompile(`(?ms:^\s?@implementation\s+` + className + `\s+.*?@end)`)
		matchedImpl := implRegexp.FindIndex(implFileBytes)
		implBytes := implFileBytes[matchedImpl[0]:matchedImpl[1]]

		classesInfo[i] = NewObjCClass(className, interfaceHBytes, interfaceMBytes, implBytes, matchedImpl[0])
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
