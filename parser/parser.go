package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

const ocgenMarker = "OCGEN_AUTO"

var headerRegexp = regexp.MustCompile(`(?ms:^\s?@interface\s+([^:<\s]*).*?` + ocgenMarker + `.*?@end)`)

const headerRegexpClassNameIndex = 1

var propertyRegexp = regexp.MustCompile(`\s?@property\s?(?:\((.*)\))?\s?([^\s\*]*)\s?(\*)?(.*);`)

const headerFileExt = ".h"

type ObjCClassInfo struct {
	MFile             os.File
	Properties        []Property
	ConformsNSCoding  bool
	ConformsNSCopying bool

	NSCodingInfo struct {
		InitWithCoder   MethodInfo
		EncodeWithCoder MethodInfo
	}
	NSCopyingInfo struct {
		CopyWithZone MethodInfo
	}
}

type MethodInfo struct {
	PosStart int64
	PosEnd   int64
}

type Property struct {
	Name, Class string
	Attributes  []string
	IsPointer   bool
}

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

func ParseAndGetClassesInfo(headerFileName string) ([]ObjCClassInfo, error) {
	headerFileBytes, err := ioutil.ReadFile(headerFileName)
	if err != nil {
		log.Printf("Unable to open header file %v\n", err)
		return nil, err
	}

	/*implFileBytes*/ _, err = ioutil.ReadFile(implFileNameFromHeader(headerFileName))
	if err != nil {
		log.Printf("Unable to open implementation file: %v\n", err)
		return nil, err
	}

	matchedInterfaces := headerRegexp.FindAllSubmatchIndex(headerFileBytes, -1)

	fmt.Println(len(matchedInterfaces))

	for _, matchedInterface := range matchedInterfaces {

		for i := 0; i < len(matchedInterface); i += 2 {
			fmt.Println(string(headerFileBytes[matchedInterface[i]:matchedInterface[i+1]]))
		}
	}

	return []ObjCClassInfo{}, nil

	// mFileName := mFileNameFromHeader(headerFileName)
	// mFile, err := os.Open(mFileName)
	// if err != nil {
	// 	log.Println("Unable to open implementation file %s", mFileName)
	// 	return nil, err
	// }

}

func classesFromHeader(hFileName string) ([]ObjCClassInfo, error) {
	return []ObjCClassInfo{}, nil
	//

	// info := ObjCClassInfo{}
	// propertyMatches := propertyRegexp.FindAllSubmatch(fileBytes, -1)

	// // Extract all properties
	// for _, propertyMatch := range propertyMatches {
	// 	//Split the attributes string and trim each of them
	// 	attributes := strings.Split(string(propertyMatch[1]), ",")
	// 	for i, attr := range attributes {
	// 		attributes[i] = strings.TrimSpace(attr)
	// 	}
	// 	class := string(propertyMatch[2])
	// 	pointer := string(propertyMatch[3])
	// 	name := string(propertyMatch[4])

	// 	// Add this property to the info
	// 	info.Properties = append(info.Properties, Property{
	// 		Name:       name,
	// 		Class:      class,
	// 		Attributes: attributes,
	// 		IsPointer:  pointer != "",
	// 	})
	// }

	// return &info, nil
}

func implFileNameFromHeader(headerFileName string) string {
	return headerFileName[:len(headerFileName)-len(headerFileExt)] + ".m"
}
