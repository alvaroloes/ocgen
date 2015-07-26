package parser

import (
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
	// These fields are extracted from the header file
	Name              string
	MFile             os.File
	Properties        []Property
	ConformsNSCoding  bool
	ConformsNSCopying bool

	// These fields are extracted from the implementation file
	NSCodingInfo struct {
		InitWithCoder   MethodInfo
		EncodeWithCoder MethodInfo
	}
	NSCopyingInfo struct {
		CopyWithZone MethodInfo
	}
}

type MethodInfo struct {
	PosStart, PosEnd int64
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

	implFile, err := os.Open(implFileNameFromHeader(headerFileName))
	if err != nil {
		log.Printf("Unable to open implementation file: %v\n", err)
		return nil, err
	}

	//TODO: Create these functions
	classesInfo := getClassesFromHeaderFile(headerFileBytes)
	fillClassesInfoFromImplFile(implFile, classesInfo)

	return classesInfo, nil

	// mFileName := mFileNameFromHeader(headerFileName)
	// mFile, err := os.Open(mFileName)
	// if err != nil {
	// 	log.Println("Unable to open implementation file %s", mFileName)
	// 	return nil, err
	// }

}

func implFileNameFromHeader(headerFileName string) string {
	return headerFileName[:len(headerFileName)-len(headerFileExt)] + ".m"
}

func getClassesFromHeaderFile(headerFileBytes []byte) []ObjCClassInfo {
	matchedInterfaces := headerRegexp.FindAllSubmatchIndex(headerFileBytes, -1)
	classesInfo := make([]ObjCClassInfo, len(matchedInterfaces))

	for i, matchedInterface := range matchedInterfaces {
		start := matchedInterface[headerRegexpClassNameIndex*2]
		end := matchedInterface[headerRegexpClassNameIndex*2+1]

		classesInfo[i] = ObjCClassInfo{
			Name: string(headerFileBytes[start:end]),
		}
	}
	return classesInfo
}

func fillClassesInfoFromImplFile(implFile *os.File, classesInfo []ObjCClassInfo) {

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
