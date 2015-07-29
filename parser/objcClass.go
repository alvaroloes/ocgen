package parser

import (
	"regexp"
	"strings"
)

var (
	classNameRegexp = regexp.MustCompile(`@interface\s+([^:<\s]*)`)
	propertyRegexp  = regexp.MustCompile(`\s?@property\s?(?:\((.*)\))?\s?([^\s\*]*)\s?(\*)?(.*);`)
)

const classNameRegexpIndex = 1
const (
	propertyRegexpAttrIndex = iota + 1
	propertyRegexpClassIndex
	propertyRegexpPointerIndex
	propertyRegexpNameIndex
)

type ObjCClass struct {
	// These fields are extracted from the header file
	Name              string
	ImplFileName      string
	Properties        []Property // These are also extracted form the implementation file
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

//TODO: pass the class name here
func NewObjCClass(hInterfaceBytes, mInterfaceBytes, implBytes []byte, implFileName string) ObjCClass {
	matchedName := classNameRegexp.FindSubmatch(hInterfaceBytes)
	propertiesFromHeader := extractProperties(hInterfaceBytes)
	//TODO: Extract properties from @interface statements in the implementation file

	class := ObjCClass{
		Name:         string(matchedName[classNameRegexpIndex]),
		ImplFileName: implFileName,
		Properties:   propertiesFromHeader,
		//TODO: Detect if the class conforms the protocols taking into account the parent protocols too
		ConformsNSCoding:  true,
		ConformsNSCopying: true,
	}

	extractProtocolMethodsInfo(&class, implBytes)

	return class
}

func extractProperties(interfaceBytes []byte) []Property {
	var properties []Property

	matchedProperties := propertyRegexp.FindAllSubmatch(interfaceBytes, -1)
	for _, matchedProperty := range matchedProperties {
		//Split the attributes string and trim each of them
		attributes := strings.Split(string(matchedProperty[propertyRegexpAttrIndex]), ",")
		for i, attr := range attributes {
			attributes[i] = strings.TrimSpace(attr)
		}
		class := string(matchedProperty[propertyRegexpClassIndex])
		pointer := string(matchedProperty[propertyRegexpPointerIndex])
		name := string(matchedProperty[propertyRegexpNameIndex])

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

func extractProtocolMethodsInfo(class *ObjCClass, implFileBytes []byte) {

}
