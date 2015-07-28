package parser

import (
	"os"
	"regexp"
	"strings"
)

var (
	classNameRegexp = regexp.MustCompile(`@interface\s+([^:<\s]*)`)
	propertyRegexp  = regexp.MustCompile(`\s?@property\s?(?:\((.*)\))?\s?([^\s\*]*)\s?(\*)?(.*);`)
)

const (
	classNameRegexpIndex    = 1
	propertyRegexpAttrIndex = 1
	propertyRegexpClassIndex
	propertyRegexpPointerIndex
	propertyRegexpNameIndex
)

type ObjCClass struct {
	// These fields are extracted from the header file
	Name              string
	MFile             *os.File
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

func NewObjCClass(interfaceBytes []byte, implFile *os.File) ObjCClass {
	matchedName := classNameRegexp.FindSubmatch(interfaceBytes)

	class := ObjCClass{
		Name:       string(matchedName[classNameRegexpIndex]),
		MFile:      implFile,
		Properties: extractProperties(interfaceBytes),
	}
	return class
}

func extractProperties(interfaceBytes []byte) []Property {
	var properties []Property

	matchedProperties := propertyRegexp.FindAllSubmatch(interfaceBytes, -1)
	for _, matchedProperty := range matchedProperties {
		//Split the attributes string and trim each of them
		attributes := strings.Split(string(matchedProperty[1]), ",")
		for i, attr := range attributes {
			attributes[i] = strings.TrimSpace(attr)
		}
		class := string(matchedProperty[2])
		pointer := string(matchedProperty[3])
		name := string(matchedProperty[4])

		// Add this property to the info
		properties = append(properties, Property{
			Name:       name,
			Class:      class,
			Attributes: attributes,
			IsPointer:  pointer != "",
		})
	}

	return properties
}
