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

var interfaceRegexp = regexp.MustCompile(`(?ms:^\s?@interface\s+([^:<\s]*).*?` + ocgenMarker + `.*?@end)`)
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

func ParseAndGetClassesInfo(hFileName string) ([]ObjCClassInfo, error) {
	classes, err := classesFromHeader(hFileName)

	return classes, err

	// mFileName := mFileNameFromHeader(headerFileName)
	// mFile, err := os.Open(mFileName)
	// if err != nil {
	// 	log.Println("Unable to open implementation file %s", mFileName)
	// 	return nil, err
	// }

}

func classesFromHeader(hFileName string) ([]ObjCClassInfo, error) {
	fileBytes, err := ioutil.ReadFile(hFileName)
	if err != nil {
		log.Println("Unable to open header file %s", hFileName)
		return nil, err
	}

	matchedInterfaces := interfaceRegexp.FindAllSubmatchIndex(fileBytes, -1)
	fmt.Println(len(matchedInterfaces))
	for _, match := range matchedInterfaces {

		for i := 0; i < len(match); i += 2 {
			fmt.Println(string(fileBytes[match[i]:match[i+1]]))
		}
	}

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

func mFileNameFromHeader(headerFileName string) string {
	return headerFileName[:len(headerFileName)-len(headerFileExt)] + ".m"
}
